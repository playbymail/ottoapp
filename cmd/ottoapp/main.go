package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

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

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/ping", pingHandler)
	mux.HandleFunc("/api/login", loginHandler)
	mux.Handle("/api/logout", csrfOnly(http.HandlerFunc(logoutHandler)))
	mux.HandleFunc("/api/session", sessionHandler) // returns CSRF
	mux.Handle("/api/me", authOnly(http.HandlerFunc(meHandler)))

	// Protect all state-changing routes with CSRF:
	protected := csrfOnly(mux)

	srv := &http.Server{
		Addr:         ":8181",
		Handler:      protected,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Println("Go API on :8181")
	log.Fatal(srv.ListenAndServe())
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","msg":"pong"}`))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}

	var body struct{ Username, Password string }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	u, ok := checkUser(body.Username, body.Password)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	sid := newID(32)
	csrf := newID(16)
	sessions[sid] = &Session{User: u, CSRF: csrf, Expiry: time.Now().Add(14 * 24 * time.Hour)}

	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,                 // HTTPS via Caddy (dev & prod)
		SameSite: http.SameSiteLaxMode, // same-site SPA+API
		MaxAge:   60 * 60 * 24 * 14,
	})

	w.WriteHeader(http.StatusNoContent)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	if sid, ok := readSID(r); ok {
		delete(sessions, sid)
		http.SetCookie(w, &http.Cookie{
			Name: "sid", Value: "", Path: "/", MaxAge: -1,
			HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode,
		})
	}
	w.WriteHeader(http.StatusNoContent)
}

func sessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	if s, ok := currentSession(r); ok {
		_ = json.NewEncoder(w).Encode(struct {
			CSRF string `json:"csrf"`
		}{CSRF: s.CSRF})
		return
	}
	http.Error(w, "unauthorized", http.StatusUnauthorized)
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	if s, ok := currentSession(r); ok {
		_ = json.NewEncoder(w).Encode(s.User)
		return
	}
	http.Error(w, "unauthorized", http.StatusUnauthorized)
}

/*** middleware & helpers ***/

func authOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := currentSession(r); !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func csrfOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only enforce on state-changing methods
		switch r.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			s, ok := currentSession(r)
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			if got := r.Header.Get("X-CSRF-Token"); got == "" || got != s.CSRF {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func currentSession(r *http.Request) (*Session, bool) {
	sid, ok := readSID(r)
	if !ok {
		return nil, false
	}
	s, ok := sessions[sid]
	if !ok {
		return nil, false
	}
	if time.Now().After(s.Expiry) {
		delete(sessions, sid)
		return nil, false
	}
	return s, true
}

func readSID(r *http.Request) (string, bool) {
	c, err := r.Cookie("sid")
	if err != nil || c.Value == "" {
		return "", false
	}
	return c.Value, true
}

func newID(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func checkUser(username, password string) (User, bool) {
	// TODO: replace with real lookup + password hash check
	if username == "admin" && password == "secret" {
		return User{ID: "1", Username: "admin", Role: "admin"}, true
	}
	return User{}, false
}
