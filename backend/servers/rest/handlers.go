// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/playbymail/ottoapp/backend/services/sessmgr"
)

/*
Sessions:what are the 4 routes for?

  POST /api/login → create session + Set-Cookie

  POST /api/logout → delete session + clear cookie

  GET /api/session → “is this cookie valid? give me csrf + user”

  GET /api/me → (optional) “give me just user again”
*/

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type Session struct {
	User   User
	CSRF   string
	Expiry time.Time
}

var sessions = map[string]*Session{}

func deleteCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name: "sid", Value: "", Path: "/", MaxAge: -1,
		HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// loginHandler creates a session and sets the cookie.
func loginHandler(auth sessmgr.AuthStore, debug, autoLogin bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if debug {
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
		var u User
		var ok bool
		if autoLogin {
			u, ok = User{ID: "catbird", Username: body.Username, Role: "guest"}, true
			if body.Password == "admin" || body.Password == "chief" {
				u.Role = body.Password
			}
		} else {
			log.Printf("%s %s: checkUser(%q, %q)\n", r.Method, r.URL.Path, body.Username, body.Password)
			u, ok = checkUser(auth, body.Username, body.Password)
		}
		if !ok {
			log.Printf("%s %s: checkUser failed\n", r.Method, r.URL.Path)
			http.Error(w, "unauthorized w", http.StatusUnauthorized)
			return
		}
		log.Printf("%s %s: user %q: role %q: authenticated\n", r.Method, r.URL.Path, u.Username, u.Role)

		sid, csrf := newID(), newID()
		sess := &Session{User: u, CSRF: csrf, Expiry: time.Now().Add(14 * 24 * time.Hour)}
		sessions[sid] = sess

		payload := struct {
			CSRF string `json:"csrf,omitempty"`
			User struct {
				ID       string `json:"id"`
				Username string `json:"username"`
				Role     string `json:"role"`
			} `json:"user"`
		}{
			CSRF: sess.CSRF,
			User: struct {
				ID       string `json:"id"`
				Username string `json:"username"`
				Role     string `json:"role"`
			}{
				ID:       u.ID,
				Username: u.Username,
				Role:     u.Role,
			},
		}

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
		_ = json.NewEncoder(w).Encode(payload)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("%s %s: method not allowed\n", r.Method, r.URL)
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	if sid, ok := readSID(r); ok {
		log.Printf("%s %s: deleting session\n", r.Method, r.URL)
		delete(sessions, sid)
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

func sessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	sess, ok := currentSession(r)
	if !ok {
		http.Error(w, "unauthorized e", http.StatusUnauthorized)
		return
	}
	payload := struct {
		CSRF string `json:"csrf,omitempty"`
		User struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			Role     string `json:"role"`
		} `json:"user"`
	}{
		CSRF: sess.CSRF,
		User: struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			Role     string `json:"role"`
		}{
			ID:       sess.User.ID,
			Username: sess.User.Username,
			Role:     sess.User.Role,
		},
	}
	// Ember Simple Auth requires an HTTP 200 response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(payload)

}

func meHandler(w http.ResponseWriter, r *http.Request) {
	if s, ok := currentSession(r); ok {
		_ = json.NewEncoder(w).Encode(s.User)
		return
	}
	http.Error(w, "unauthorized x", http.StatusUnauthorized)
}
