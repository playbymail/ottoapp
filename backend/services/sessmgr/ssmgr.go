// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sessmgr

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"log"
	"math/rand/v2"
	"net/http"
	"sync"
	"time"
)

// Based on the excellent articles by Mohamed Said at themsaid.com.
// * https://themsaid.com/building-secure-session-manager-in-go

type Session struct {
	mu             sync.Mutex
	createdAt      time.Time
	lastActivityAt time.Time
	id             string
	data           map[string]any
}

type SessionStore interface {
	read(id string) (*Session, error)
	write(session *Session) error
	destroy(id string) error
	gc(idleExpiration, absoluteExpiration time.Duration) error
}

type SessionManager struct {
	store              SessionStore
	idleExpiration     time.Duration
	absoluteExpiration time.Duration
	cookieName         string
}

// generateSessionId creates a Base64 URL–encoded string from 32 bytes of
// non-cryptographic random data. Intended for demos and tests only.
//
// ⚠️ Not secure! Use crypto/rand for production systems.
func generateSessionId() string {
	id := make([]byte, 32)
	binary.LittleEndian.PutUint64(id[0*8:], rand.Uint64())
	binary.LittleEndian.PutUint64(id[1*8:], rand.Uint64())
	binary.LittleEndian.PutUint64(id[2*8:], rand.Uint64())
	binary.LittleEndian.PutUint64(id[3*8:], rand.Uint64())
	return base64.RawURLEncoding.EncodeToString(id)
}

// Example secure version (production):
// func generateSecureSessionId() string {
// 	id := make([]byte, 32)
// 	if _, err := crypto/rand.Read(id); err != nil {
// 		panic(err)
// 	}
// 	return base64.RawURLEncoding.EncodeToString(id)
// }

func newSession() *Session {
	return &Session{
		id:             generateSessionId(),
		data:           map[string]any{"csrf_token": generateCSRFToken()},
		createdAt:      time.Now(),
		lastActivityAt: time.Now(),
	}
}

func (s *Session) Get(key string) any {
	return s.data[key]
}

func (s *Session) Put(key string, value any) {
	s.data[key] = value
}

func (s *Session) Delete(key string) {
	delete(s.data, key)
}

func NewSessionManager(
	store SessionStore,
	gcInterval,
	idleExpiration,
	absoluteExpiration time.Duration,
	cookieName string) *SessionManager {

	m := &SessionManager{
		store:              store,
		idleExpiration:     idleExpiration,
		absoluteExpiration: absoluteExpiration,
		cookieName:         cookieName,
	}

	go m.gc(gcInterval)

	return m
}

func (m *SessionManager) gc(d time.Duration) {
	ticker := time.NewTicker(d)

	for range ticker.C {
		m.store.gc(m.idleExpiration, m.absoluteExpiration)
	}
}

func (m *SessionManager) validate(session *Session) bool {
	if time.Since(session.createdAt) > m.absoluteExpiration ||
		time.Since(session.lastActivityAt) > m.idleExpiration {

		// Delete the session from the store
		err := m.store.destroy(session.id)
		if err != nil {
			panic(err)
		}

		return false
	}

	return true
}

type sessionContextKeyType string

const sessionContextKey = sessionContextKeyType("")

func (m *SessionManager) start(r *http.Request) (*Session, *http.Request) {
	var session *Session

	// Read From Cookie
	cookie, err := r.Cookie(m.cookieName)
	if err == nil {
		session, err = m.store.read(cookie.Value)
		if err != nil {
			log.Printf("Failed to read session from store: %v", err)
		}
	}

	// Generate a new session
	if session == nil || !m.validate(session) {
		session = newSession()
	}

	// Attach session to context
	ctx := context.WithValue(r.Context(), sessionContextKey, session)
	r = r.WithContext(ctx)

	return session, r
}

func (m *SessionManager) save(session *Session) error {
	session.lastActivityAt = time.Now()

	err := m.store.write(session)
	if err != nil {
		return err
	}

	return nil
}

func (m *SessionManager) migrate(session *Session) error {
	session.mu.Lock()
	defer session.mu.Unlock()

	err := m.store.destroy(session.id)
	if err != nil {
		return err
	}

	session.id = generateSessionId()

	return nil
}

func (m *SessionManager) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Start the session
		session, rws := m.start(r)

		// Create a new response writer
		sw := &sessionResponseWriter{
			ResponseWriter: w,
			sessionManager: m,
			request:        rws,
		}

		// Add essential headers
		w.Header().Add("Vary", "Cookie")
		w.Header().Add("Cache-Control", `no-cache="Set-Cookie"`)

		// Verify CSRF token for state-changing requests
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch || r.Method == http.MethodDelete {
			if !m.verifyCSRFToken(rws, session) {
				http.Error(sw, "CSRF token mismatch", http.StatusForbidden)
				return
			}
		}

		// Call the next handler and pass the new response writer and new request
		next.ServeHTTP(sw, rws)

		// Save the session
		m.save(session)

		// Write the session cookie to the response if not already written
		writeCookieIfNecessary(sw)
	})
}

type sessionResponseWriter struct {
	http.ResponseWriter
	sessionManager *SessionManager
	request        *http.Request
	done           bool
}

func (w *sessionResponseWriter) Write(b []byte) (int, error) {
	writeCookieIfNecessary(w)

	return w.ResponseWriter.Write(b)
}

func (w *sessionResponseWriter) WriteHeader(code int) {
	writeCookieIfNecessary(w)

	w.ResponseWriter.WriteHeader(code)
}

func (w *sessionResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func writeCookieIfNecessary(w *sessionResponseWriter) {
	if w.done {
		return
	}

	session, ok := w.request.Context().Value(sessionContextKey).(*Session)
	if !ok {
		panic("session not found in request context")
	}

	cookie := &http.Cookie{
		Name:     w.sessionManager.cookieName,
		Value:    session.id,
		Domain:   "mywebsite.com",
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(w.sessionManager.idleExpiration),
		MaxAge:   int(w.sessionManager.idleExpiration / time.Second),
	}

	http.SetCookie(w.ResponseWriter, cookie)

	w.done = true
}

func GetSession(r *http.Request) *Session {
	session, ok := r.Context().Value(sessionContextKey).(*Session)
	if !ok {
		panic("session not found in request context")
	}

	return session
}

type InMemorySessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{
		sessions: make(map[string]*Session),
	}
}

func (s *InMemorySessionStore) read(id string) (*Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, _ := s.sessions[id]

	return session, nil
}

func (s *InMemorySessionStore) write(session *Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[session.id] = session

	return nil
}

func (s *InMemorySessionStore) destroy(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, id)

	return nil
}

func (s *InMemorySessionStore) gc(idleExpiration, absoluteExpiration time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, session := range s.sessions {
		if time.Since(session.lastActivityAt) > idleExpiration ||
			time.Since(session.createdAt) > absoluteExpiration {
			delete(s.sessions, id)
		}
	}

	return nil
}

// example using SessionManager
//func main() {
//	sessionManager := ssmgr.NewSessionManager(
//		ssmgr.NewInMemorySessionStore(),
//		30 * time.Minute,
//		1 * time.Hour,
//		12 * time.Hour,
//		"session",
//	)
//  mux.HandleFunc("/project", func(w http.ResponseWriter, r *http.Request) {
//    session := sess.GetSession(r)
//    currentProject := session.Get("current_project")
//  })
//  mux.HandleFunc("/projects/switch/some-project-id", func(w http.ResponseWriter, r *http.Request) {
//    session := sess.GetSession(r)
//    session.Put("current_project", "some-project-id")
//  })
//  mux.HandleFunc("POST /projects/switch/some-project-id", func(w http.ResponseWriter, r *http.Request) {
//    session := sess.GetSession(r)
//    session.Put("current_project", "some-project-id")
//  })
//  server := &http.Server{
//    Addr:    ":8080",
//    Handler: sessionManager.Handle(mux), // Here
//  }
//  server.ListenAndServe()
//}
