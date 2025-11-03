// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package sessions implements a very simple session manager.
package sessions

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
)

/*
Sessions:what are the 4 routes for?

  POST /api/login   → create session + Set-Cookie

  POST /api/logout  → delete session + clear cookie

  GET /api/session  → “is this cookie valid? give me csrf + user”

  GET /api/me       → “is this cookie valid? give me csrf + user”
*/

func NewSessionManager(authStore AuthStore, sessionStore SessionStore, ttl time.Duration, gcInterval time.Duration) (SessionManager, error) {
	if ttl < 1*time.Minute {
		return nil, domains.ErrInvalidTtl
	} else if gcInterval < 1*time.Minute {
		return nil, domains.ErrInvalidGcInterval
	}
	m := &Manager{
		ttl: ttl,
	}
	m.stores.auth = authStore
	m.stores.sessions = sessionStore

	// start the go routine that will delete expired sessions
	log.Printf("[sessions] reaping %v\n", gcInterval)
	if err := m.stores.sessions.ReapSessions(); err != nil {
		log.Printf("[sessions] reap %v: %v\n", gcInterval, err)
	}
	go func(d time.Duration) {
		ticker := time.NewTicker(d)
		for range ticker.C {
			//log.Printf("[sessions] reaping %v\n", time.Now().UTC())
			if err := m.stores.sessions.ReapSessions(); err != nil {
				log.Printf("[sessions] reap %v: %v\n", d, err)
			}
		}
	}(gcInterval)

	return m, nil
}

type AuthStore interface {
	AuthenticateUser(username, password string) (domains.ID, error)
	GetUserByID(domains.ID) (*domains.User_t, error)
}

type SessionStore interface {
	SessionStoreID() string // unique identifier for the store
	CreateSession(userId domains.ID, ttl time.Duration) (*domains.Session, error)
	ReadSession(id domains.SessionId) (*domains.Session, error)
	UpdateSession(*domains.Session) error
	DeleteSession(id domains.SessionId) error
	ReapSessions() error
}

type SessionManager interface {
	DeleteCookie(w http.ResponseWriter, r *http.Request)
	GetCurrentSession(r *http.Request) (*domains.Session, error)
	GetMeHandler(w http.ResponseWriter, r *http.Request)
	GetSessionHandler(w http.ResponseWriter, r *http.Request)
	PostLoginHandler(w http.ResponseWriter, r *http.Request)
	PostLogoutHandler(w http.ResponseWriter, r *http.Request)
}

type Manager struct {
	autoLogin bool
	debug     bool
	ttl       time.Duration
	stores    struct {
		auth     AuthStore
		sessions SessionStore
	}
}

// DeleteCookie deletes the session cookie and redirects the browser
func (m *Manager) DeleteCookie(w http.ResponseWriter, r *http.Request) {
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

func (m *Manager) GetCurrentSession(r *http.Request) (*domains.Session, error) {
	//log.Printf("%s %s: gcs: entered\n", r.Method, r.URL.Path)
	//log.Printf("%s %s: gcs: storeId %q\n", r.Method, r.URL.Path, m.stores.sessions.SessionStoreID())
	sid, ok := readSID(r)
	//log.Printf("%s %s: gcs: sid %q: ok %v\n", r.Method, r.URL.Path, sid, ok)
	if !ok {
		return nil, domains.ErrNotExists
	}
	sess, err := m.stores.sessions.ReadSession(sid)
	if err != nil {
		//log.Printf("%s %s: gcs: read %q: failed %v\n", r.Method, r.URL.Path, sid, err)
		return nil, err
	}
	return sess, nil
}

func (m *Manager) GetMeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	sess, err := m.GetCurrentSession(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	// Ember Simple Auth requires an HTTP 200 response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(toPayload(sess.Csrf, sess.User).User)
	return
}

func (m *Manager) GetSessionHandler(w http.ResponseWriter, r *http.Request) {
	//log.Printf("%s %s: entered\n", r.Method, r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	sess, err := m.GetCurrentSession(r)
	if err != nil {
		//log.Printf("%s %s: gcs: error %v\n", r.Method, r.URL.Path, err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	//log.Printf("%s %s: gcs %+v\n", r.Method, r.URL.Path, sess)
	// Ember Simple Auth requires an HTTP 200 response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(toPayload(sess.Csrf, sess.User))
}

// PostLoginHandler creates a session and sets the cookie.
func (m *Manager) PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var body struct{ Username, Password string }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		//log.Printf("%s %s: bad json\n", r.Method, r.URL.Path)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	var u domains.User_t
	var ok bool
	if m.autoLogin {
		u, ok = domains.User_t{ID: 42, Username: body.Username, Roles: map[string]bool{"authenticated": true}}, true
		if body.Password == "admin" || body.Password == "chief" {
			u.Roles[body.Password] = true
		} else {
			u.Roles["guest"] = true
		}
	} else {
		//log.Printf("%s %s: checkUser(%q, %q)\n", r.Method, r.URL.Path, body.Username, body.Password)
		u, ok = authenticateCredentials(m.stores.auth, body.Username, body.Password)
	}
	if !ok {
		//log.Printf("%s %s: authenticateCredentials failed\n", r.Method, r.URL.Path)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	//log.Printf("%s %s: user %q: roles %+v: authenticated\n", r.Method, r.URL.Path, u.Username, u.Roles)

	sess, err := m.stores.sessions.CreateSession(u.ID, m.ttl)
	if err != nil {
		//log.Printf("%s %s: user %q: createSession %v\n", r.Method, r.URL.Path, u.ID, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Ember Simple Auth doesn't use the cookie - it's for our session manager
	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    string(sess.Id),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,                 // HTTPS via Caddy (dev & prod)
		SameSite: http.SameSiteLaxMode, // same-site SPA+API
		MaxAge:   60 * 60 * 24 * 14,
	})

	// Ember Simple Auth requires an HTTP 200 response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(toPayload(sess.Csrf, u))
}

func (m *Manager) PostLogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	if sid, ok := readSID(r); ok {
		//log.Printf("%s %s: deleting session\n", r.Method, r.URL)
		err := m.stores.sessions.DeleteSession(sid)
		if err != nil && !errors.Is(err, domains.ErrNotExists) {
			//log.Printf("%s %s: deleteSession\n", r.Method, r.URL)
		}
	}
	//log.Printf("%s %s: deleting cookie\n", r.Method, r.URL)
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
	//log.Printf("%s %s: sending response\n", r.Method, r.URL)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	})
	if err != nil {
		//log.Printf("%s %s: write: json %v\n", r.Method, r.URL, err)
	}
}
