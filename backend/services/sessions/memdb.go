// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sessions

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

func NewInMemorySessionManager(authStore AuthStore_i) (SessionManager_i, error) {
	m := &memdbManager{
		sessions: make(map[string]*Session_t),
	}
	m.stores.auth = authStore
	return m, nil
}

type memdbManager struct {
	stores struct {
		auth AuthStore_i
	}
	mu        sync.Mutex
	autoLogin bool
	debug     bool
	sessions  map[string]*Session_t
}

// DeleteCookie deletes the session cookie and redirects the browser
func (m *memdbManager) DeleteCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *memdbManager) GetCurrentSession(r *http.Request) (*Session_t, bool) {
	log.Printf("%s %s: currentSession: entered\n", r.Method, r.URL.Path)
	sid, ok := readSID(r)
	if !ok {
		log.Printf("%s %s: currentSession: readSID false\n", r.Method, r.URL.Path)
		return nil, false
	}
	log.Printf("%s %s: currentSession: readSID %q\n", r.Method, r.URL.Path, sid)
	sess, ok := m.getSession(sid)
	if !ok {
		log.Printf("%s %s: currentSession: sessions false\n", r.Method, r.URL.Path)
		return nil, false
	}
	log.Printf("%s %s: currentSession: sessions %+v\n", r.Method, r.URL.Path, sess)
	return sess, true
}

func (m *memdbManager) GetMeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	sid, ok := readSID(r)
	if !ok {
		log.Printf("%s %s: readSID false\n", r.Method, r.URL.Path)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	sess, ok := m.getSession(sid)
	if !ok {
		log.Printf("%s %s: getSession false\n", r.Method, r.URL.Path)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	// Ember Simple Auth requires an HTTP 200 response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(toPayload(sess.CSRF, sess.User).User)
	return
}

func (m *memdbManager) GetSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	sid, ok := readSID(r)
	if !ok {
		log.Printf("%s %s: readSID false\n", r.Method, r.URL.Path)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("%s %s: readSID %q\n", r.Method, r.URL.Path, sid)
	sess, ok := m.getSession(sid)
	if !ok {
		log.Printf("%s %s: getSession false\n", r.Method, r.URL.Path)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	// Ember Simple Auth requires an HTTP 200 response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(toPayload(sess.CSRF, sess.User))
}

// PostLoginHandler creates a session and sets the cookie.
func (m *memdbManager) PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	if m.debug {
		log.Printf("%s %s: entered\n", r.Method, r.URL.Path)
	}
	if r.Method != http.MethodPost {
		log.Printf("%s %s: method not allowed\n", r.Method, r.URL.Path)
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}

	var body struct{ Username, Password string }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("%s %s: bad json\n", r.Method, r.URL.Path)
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	var u User_t
	var ok bool
	if m.autoLogin {
		u, ok = User_t{ID: "catbird", Username: body.Username, Roles: map[string]bool{"authenticated": true, "guest": true}}, true
		if body.Password == "admin" || body.Password == "chief" {
			u.Roles[body.Password] = true
		}
	} else {
		log.Printf("%s %s: checkUser(%q, %q)\n", r.Method, r.URL.Path, body.Username, body.Password)
		u, ok = authenticateCredentials(m.stores.auth, body.Username, body.Password)
	}
	if !ok {
		log.Printf("%s %s: authenticateCredentials failed\n", r.Method, r.URL.Path)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("%s %s: user %q: roles %+v: authenticated\n", r.Method, r.URL.Path, u.Username, u.Roles)

	sid, csrf := newSID(), newCSRF()
	sess := &Session_t{User: u, CSRF: csrf, Expiry: time.Now().Add(14 * 24 * time.Hour)}
	m.sessions[sid] = sess

	// Ember Simple Auth doesn't use the cookie - it's for our session manager
	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,                 // HTTPS via Caddy (dev & prod)
		SameSite: http.SameSiteLaxMode, // same-site SPA+API
		MaxAge:   60 * 60 * 24 * 14,
	})

	// Ember Simple Auth requires an HTTP 200 response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(toPayload(sess.CSRF, u))
}

func (m *memdbManager) PostLogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("%s %s: method not allowed\n", r.Method, r.URL)
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	if sid, ok := readSID(r); ok {
		log.Printf("%s %s: deleting session\n", r.Method, r.URL)
		m.mu.Lock()
		delete(m.sessions, sid)
		m.mu.Unlock()
	}
	log.Printf("%s %s: deleting cookie\n", r.Method, r.URL)
	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	// Ember Simple Auth requires an HTTP 200 response
	log.Printf("%s %s: sending response\n", r.Method, r.URL)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	})
	if err != nil {
		log.Printf("%s %s: write: json %v\n", r.Method, r.URL, err)
	}
}

func (m *memdbManager) getSession(sid string) (*Session_t, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	s, ok := m.sessions[sid]
	if !ok {
		return nil, false
	} else if time.Now().After(s.Expiry) {
		delete(m.sessions, sid)
		return nil, false
	}
	return s, true
}
