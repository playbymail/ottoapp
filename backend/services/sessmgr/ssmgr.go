// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sessmgr

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

// Based on the excellent articles by Mohamed Said at themsaid.com.
// * https://themsaid.com/building-secure-session-manager-in-go

type Session_t struct {
	mu             sync.Mutex
	CreatedAt      time.Time
	LastActivityAt time.Time
	Id             string
	Data           map[string]any
}

type SessionStore_i interface {
	readSession(id string) (*Session_t, error)
	writeSession(session *Session_t) error
	destroySession(id string) error
	gcSessions(idleExpiration, absoluteExpiration time.Duration) error
}

type SessionManager_t struct {
	store              SessionStore_i
	idleExpiration     time.Duration
	absoluteExpiration time.Duration
	cookieName         string
}

func newSession() *Session_t {
	return &Session_t{
		Id:             generateSessionId(),
		Data:           map[string]any{"csrf_token": generateCSRFToken()},
		CreatedAt:      time.Now(),
		LastActivityAt: time.Now(),
	}
}

func (s *Session_t) Get(key string) any {
	return s.Data[key]
}

func (s *Session_t) Put(key string, value any) {
	s.Data[key] = value
}

func (s *Session_t) Delete(key string) {
	delete(s.Data, key)
}

func NewSessionManager(
	store SessionStore_i,
	gcInterval,
	idleExpiration,
	absoluteExpiration time.Duration,
	cookieName string) *SessionManager_t {

	m := &SessionManager_t{
		store:              store,
		idleExpiration:     idleExpiration,
		absoluteExpiration: absoluteExpiration,
		cookieName:         cookieName,
	}

	go m.gc(gcInterval)

	return m
}

func (m *SessionManager_t) gc(d time.Duration) {
	ticker := time.NewTicker(d)

	for range ticker.C {
		m.store.gcSessions(m.idleExpiration, m.absoluteExpiration)
	}
}

func (m *SessionManager_t) validate(session *Session_t) bool {
	if time.Since(session.CreatedAt) > m.absoluteExpiration ||
		time.Since(session.LastActivityAt) > m.idleExpiration {

		// Delete the session from the store
		err := m.store.destroySession(session.Id)
		if err != nil {
			panic(err)
		}

		return false
	}

	return true
}

type sessionContextKeyType string

const sessionContextKey = sessionContextKeyType("")

func (m *SessionManager_t) start(r *http.Request) (*Session_t, *http.Request) {
	var session *Session_t

	// Read From Cookie
	cookie, err := r.Cookie(m.cookieName)
	if err == nil {
		session, err = m.store.readSession(cookie.Value)
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

func (m *SessionManager_t) save(session *Session_t) error {
	session.LastActivityAt = time.Now()

	err := m.store.writeSession(session)
	if err != nil {
		return err
	}

	return nil
}

func (m *SessionManager_t) migrate(session *Session_t) error {
	session.mu.Lock()
	defer session.mu.Unlock()

	err := m.store.destroySession(session.Id)
	if err != nil {
		return err
	}

	session.Id = generateSessionId()

	return nil
}

func (m *SessionManager_t) Handle(next http.Handler) http.Handler {
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
	sessionManager *SessionManager_t
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

	session, ok := w.request.Context().Value(sessionContextKey).(*Session_t)
	if !ok {
		panic("session not found in request context")
	}

	cookie := &http.Cookie{
		Name:     w.sessionManager.cookieName,
		Value:    session.Id,
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

func GetSession(r *http.Request) *Session_t {
	session, ok := r.Context().Value(sessionContextKey).(*Session_t)
	if !ok {
		panic("session not found in request context")
	}

	return session
}

type InMemorySessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*Session_t
}

func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{
		sessions: make(map[string]*Session_t),
	}
}

func (s *InMemorySessionStore) readSession(id string) (*Session_t, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, _ := s.sessions[id]

	return session, nil
}

func (s *InMemorySessionStore) writeSession(session *Session_t) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[session.Id] = session

	return nil
}

func (s *InMemorySessionStore) destroySession(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, id)

	return nil
}

func (s *InMemorySessionStore) gcSessions(idleExpiration, absoluteExpiration time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, session := range s.sessions {
		if time.Since(session.LastActivityAt) > idleExpiration ||
			time.Since(session.CreatedAt) > absoluteExpiration {
			delete(s.sessions, id)
		}
	}

	return nil
}

// example using SessionManager_t
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
