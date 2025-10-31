// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/playbymail/ottoapp/backend/stores/sqlite"
)

type Server struct {
	http.Server
	graceTimer    time.Duration
	shutdownTimer time.Duration
}

func New(db *sqlite.DB, options ...Option) (*Server, error) {
	s := &Server{
		Server: http.Server{
			Addr:         ":8181",
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}
	s.Handler = Routes(s)

	for _, opt := range options {
		err := opt(s)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Server) Run() error {
	log.Printf("[rest] serving API on %q\n", fmt.Sprintf("http://localhost%s", s.Addr))
	if s.shutdownTimer != 0 {
		log.Printf("[rest] server timeout %v\n", s.shutdownTimer)
	}
	if s.graceTimer != 0 {
		log.Printf("[rest] shutdown delay %v\n", s.graceTimer)
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("server listening on port %s\n", s.Addr)
		serverErrors <- s.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	if s.shutdownTimer > 0 {
		go func() {
			time.Sleep(s.shutdownTimer)
			log.Printf("[rest] timeout reached (%v), initiating shutdown\n", s.shutdownTimer)
			shutdown <- syscall.SIGTERM
		}()
	}

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Printf("[rest] received signal %v, starting graceful shutdown\n", sig)

		ctx, cancel := context.WithTimeout(context.Background(), s.graceTimer)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			log.Printf("[rest] error during shutdown: %v\n", err)
			return s.Close()
		}

		log.Println("[rest] server stopped gracefully")
	}

	return nil
}

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

func Routes(s *Server) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","msg":"pong"}`))
	})

	mux.HandleFunc("POST /api/login", loginHandler)
	mux.Handle("POST /api/logout", http.HandlerFunc(logoutHandler))
	mux.Handle("GET /api/me", authOnly(http.HandlerFunc(meHandler)))
	mux.HandleFunc("GET /api/session", sessionHandler) // returns CSRF

	// Protect all state-changing routes with CSRF:
	protected := csrfOnly(mux)

	// Add logging middleware
	logged := loggingMiddleware(protected)

	return logged
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s: entered\n", r.Method, r.URL.Path)
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
	log.Printf("%s %s: checkUser(%q, %q)\n", r.Method, r.URL.Path, body.Username, body.Password)
	u, ok := checkUser(body.Username, body.Password)
	if !ok {
		log.Printf("%s %s: checkUser failed\n", r.Method, r.URL.Path)
		http.Error(w, "unauthorized w", http.StatusUnauthorized)
		return
	}

	sid := newID()
	csrf := newID()
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
	http.Error(w, "unauthorized e", http.StatusUnauthorized)
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	if s, ok := currentSession(r); ok {
		_ = json.NewEncoder(w).Encode(s.User)
		return
	}
	http.Error(w, "unauthorized x", http.StatusUnauthorized)
}

/*** middleware & helpers ***/

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		log.Printf("%s %s: logging\n", r.Method, r.URL.Path)
		next.ServeHTTP(rec, r)
		log.Printf("%s %s - %d - %v", r.Method, r.URL.Path, rec.status, time.Since(start))
	})
}

func authOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := currentSession(r); !ok {
			http.Error(w, "unauthorized z", http.StatusUnauthorized)
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
				log.Printf("%s %s: csrf: !ok\n", r.Method, r.URL.Path)
				http.Error(w, "unauthorized b", http.StatusUnauthorized)
				return
			}
			if got := r.Header.Get("X-CSRF-Token"); got == "" {
				log.Printf("%s %s: csrf: forbidden: no token\n", r.Method, r.URL.Path)
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			} else if got != s.CSRF {
				log.Printf("%s %s: csrf: forbidden: %q != %q\n", r.Method, r.URL.Path, got, s.CSRF)
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

func newID() string {
	id := make([]byte, 32)
	binary.LittleEndian.PutUint64(id[0*8:], rand.Uint64())
	binary.LittleEndian.PutUint64(id[1*8:], rand.Uint64())
	binary.LittleEndian.PutUint64(id[2*8:], rand.Uint64())
	binary.LittleEndian.PutUint64(id[3*8:], rand.Uint64())
	return base64.RawURLEncoding.EncodeToString(id)
}

var (
	fakeUsers struct {
		sync.Mutex
		users map[string]string
	}
)

func init() {
	fakeUsers.users = make(map[string]string)
	fakeUsers.users["admin"] = "1"
}

func checkUser(username, password string) (User, bool) {
	log.Printf("checkUser(%q, %q)\n", username, password)
	// TODO: replace with real lookup + password hash check
	if username == "admin" && password == "secret" {
		return User{ID: "1", Username: "admin", Role: "admin"}, true
	} else if 3 < len(username) && len(username) < 8 {
		fakeUsers.Lock()
		defer fakeUsers.Unlock()
		if id, ok := fakeUsers.users[username]; ok {
			return User{ID: id, Username: username, Role: "guest"}, true
		}
		id := fmt.Sprintf("%d", len(fakeUsers.users)+1)
		fakeUsers.users[username] = id
		return User{ID: id, Username: username, Role: "guest"}, true
	}
	return User{}, false
}
