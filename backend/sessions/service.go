// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sessions

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"slices"
	"time"

	"github.com/playbymail/ottoapp/backend/auth"
	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
	"github.com/playbymail/ottoapp/backend/users"
)

// Service provides session management operations.
//
// Based on the excellent articles by Mohamed Said at themsaid.com.
// * https://themsaid.com/building-secure-session-manager-in-go
// * https://themsaid.com/session-authentication-go
// * https://themsaid.com/csrf-protection-go-web-applications
type Service struct {
	db         *sqlite.DB
	authSvc    *auth.Service
	usersSvc   *users.Service
	ttl        time.Duration
	gcInterval time.Duration
}

func New(db *sqlite.DB, authSvc *auth.Service, usersSvc *users.Service, ttl time.Duration, gcInterval time.Duration) (*Service, error) {
	if ttl < 1*time.Minute {
		return nil, domains.ErrInvalidTtl
	} else if gcInterval < 1*time.Minute {
		return nil, domains.ErrInvalidGcInterval
	}

	s := &Service{
		db:         db,
		authSvc:    authSvc,
		usersSvc:   usersSvc,
		ttl:        ttl,
		gcInterval: gcInterval,
	}

	// start the go routine that will delete expired sessions
	log.Printf("[sessions] reaping %v\n", gcInterval)
	if err := s.ReapSessions(); err != nil {
		log.Printf("[sessions] reap %v: %v\n", gcInterval, err)
	}
	go func(d time.Duration) {
		ticker := time.NewTicker(d)
		for range ticker.C {
			//log.Printf("[sessions] reaping %v\n", time.Now().UTC())
			if err := s.ReapSessions(); err != nil {
				log.Printf("[sessions] reap %v: %v\n", d, err)
			}
		}
	}(gcInterval)

	return s, nil
}

func (s *Service) CreateSession(user *domains.User_t, ttl time.Duration) (*domains.Session, error) {
	q := s.db.Queries()
	ctx := s.db.Context()

	now := time.Now().UTC()
	sess := &domains.Session{
		Id:             newSessionId(),
		Csrf:           newCsrf(),
		User:           *user,
		ExpiresAt:      now.Add(ttl),
		LastActivityAt: now,
	}

	err := q.CreateSession(ctx, sqlc.CreateSessionParams{
		SessionID: string(sess.Id),
		Csrf:      sess.Csrf,
		UserID:    int64(sess.User.ID),
		ExpiresAt: sess.ExpiresAt.Unix(),
		CreatedAt: now.Unix(),
		UpdatedAt: now.Unix(),
	})
	if err != nil {
		return nil, err
	}

	return sess, nil
}

// DeleteCookie deletes the session cookie
func (s *Service) DeleteCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func (s *Service) DeleteSession(id domains.SessionId) error {
	q := s.db.Queries()
	ctx := s.db.Context()
	err := q.DeleteSession(ctx, string(id))
	if err != nil {
		log.Printf("[sessions] deleteSession %v\n", err)
	}
	return nil
}

func (s *Service) GetCurrentSession(r *http.Request) (*domains.Session, error) {
	//log.Printf("%s %s: gcs: entered\n", r.Method, r.URL.Path)
	//log.Printf("%s %s: gcs: storeId %q\n", r.Method, r.URL.Path, m.stores.sessions.SessionStoreID())
	sid, ok := readSID(r)
	//log.Printf("%s %s: gcs: sid %q: ok %v\n", r.Method, r.URL.Path, sid, ok)
	if !ok {
		return nil, domains.ErrNotExists
	}
	sess, err := s.ReadSession(sid)
	if err != nil {
		//log.Printf("%s %s: gcs: read %q: failed %v\n", r.Method, r.URL.Path, sid, err)
		return nil, err
	}
	return sess, nil
}

// GetCurrentUserID returns the user ID from the current session.
func (s *Service) GetCurrentUserID(r *http.Request) (domains.ID, error) {
	sess, err := s.GetCurrentSession(r)
	if err != nil {
		return domains.InvalidID, err
	}
	return sess.User.ID, nil
}

func (s *Service) HandleGetSession(w http.ResponseWriter, r *http.Request) {
	//log.Printf("%s %s: entered\n", r.Method, r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	sess, err := s.GetCurrentSession(r)
	if err != nil {
		//log.Printf("%s %s: gcs: error %v\n", r.Method, r.URL.Path, err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	//log.Printf("%s %s: gcs %+v\n", r.Method, r.URL.Path, sess)
	// Ember Simple Auth requires an HTTP 200 response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(toPayload(sess.Csrf, &sess.User))
}

// HandlePostLogin creates a session and sets the cookie.
func (s *Service) HandlePostLogin(w http.ResponseWriter, r *http.Request) {
	// always invalidate any old cookies
	s.DeleteCookie(w, r)

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var body struct{ Email, Password string }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("%s %s: bad json\n", r.Method, r.URL.Path)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	//log.Printf("%s %s: checkUser(%q, %q)\n", r.Method, r.URL.Path, body.Email, body.Password)
	userID, err := s.authSvc.AuthenticateWithEmailSecret(body.Email, body.Password)
	if err != nil {
		//log.Printf("%s %s: checkUser: %v\n", r.Method, r.URL.Path, err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	user, err := s.usersSvc.GetUserByID(userID)
	if err != nil {
		//log.Printf("%s %s: checkUser: %v\n", r.Method, r.URL.Path, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	user.Roles, err = s.authSvc.GetUserRoles(user.ID)
	if err != nil {
		//log.Printf("%s %s: checkUser: %v\n", r.Method, r.URL.Path, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	//log.Printf("%s %s: user %q: roles %+v: authenticated\n", r.Method, r.URL.Path, user.Username, user.Roles)

	sess, err := s.CreateSession(user, s.ttl)
	if err != nil {
		log.Printf("%s %s: user %q: createSession %v\n", r.Method, r.URL.Path, user.Username, err)
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
	_ = json.NewEncoder(w).Encode(toPayload(sess.Csrf, user))
}

func (s *Service) HandlePostLogout(w http.ResponseWriter, r *http.Request) {
	// always invalidate any old cookies
	s.DeleteCookie(w, r)

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	if sid, ok := readSID(r); ok {
		//log.Printf("%s %s: deleting session\n", r.Method, r.URL)
		err := s.DeleteSession(sid)
		if err != nil && !errors.Is(err, domains.ErrNotExists) {
			//log.Printf("%s %s: deleteSession\n", r.Method, r.URL)
		}
	}

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

func (s *Service) ReadSession(sessionId domains.SessionId) (*domains.Session, error) {
	q := s.db.Queries()
	ctx := s.db.Context()
	//log.Printf("readSession %q\n", sessionId)
	ss, err := q.GetSession(ctx, string(sessionId))
	if err != nil {
		//log.Printf("readSession %q: getSession %v\n", sessionId, err)
		return nil, errors.Join(domains.ErrNotExists, err)
	} else if !time.Now().UTC().Before(time.Unix(ss.ExpiresAt, 0)) {
		//log.Printf("readSession %q: expired %v\n", sessionId, ss.ExpiresAt)
		//log.Printf("readSession %q: not is  %v\n", sessionId, time.Now().UTC())
		return nil, domains.ErrNotExists
	}
	user, err := s.usersSvc.GetUserByID(domains.ID(ss.UserID))
	if err != nil {
		return nil, errors.Join(domains.ErrNotExists, err)
	}
	user.Roles, err = s.authSvc.GetUserRoles(user.ID)
	if err != nil {
		return nil, domains.ErrNoRolesAssigned
	}
	sess := &domains.Session{
		Id:             sessionId,
		Csrf:           ss.Csrf,
		User:           *user,
		ExpiresAt:      time.Unix(ss.ExpiresAt, 0),
		LastActivityAt: time.Now().UTC(),
	}
	return sess, nil
}

func (s *Service) ReapSessions() error {
	q := s.db.Queries()
	ctx := s.db.Context()
	return q.ReapSessions(ctx, time.Now().UTC().Unix())
}

func (s *Service) SessionStoreID() string {
	return "backend/stores/sqlite:DB"
}

func (s *Service) UpdateSession(session *domains.Session) error {
	q := s.db.Queries()
	ctx := s.db.Context()
	if session == nil {
		return domains.ErrSessionInvalid
	}
	return q.UpdateSession(ctx, sqlc.UpdateSessionParams{
		ExpiresAt: session.ExpiresAt.Unix(),
		SessionID: string(session.Id),
	})
}

func newCsrf() string {
	idBuffer := make([]byte, 16)
	binary.LittleEndian.PutUint64(idBuffer[0*8:], rand.Uint64())
	binary.LittleEndian.PutUint64(idBuffer[1*8:], rand.Uint64())
	return base64.RawURLEncoding.EncodeToString(idBuffer)
}

func newSessionId() domains.SessionId {
	idBuffer := make([]byte, 32)
	binary.LittleEndian.PutUint64(idBuffer[0*8:], rand.Uint64())
	binary.LittleEndian.PutUint64(idBuffer[1*8:], rand.Uint64())
	binary.LittleEndian.PutUint64(idBuffer[2*8:], rand.Uint64())
	binary.LittleEndian.PutUint64(idBuffer[3*8:], rand.Uint64())
	return domains.SessionId(base64.RawURLEncoding.EncodeToString(idBuffer))
}

func readSID(r *http.Request) (domains.SessionId, bool) {
	c, err := r.Cookie("sid")
	if err != nil {
		return "", false
	}
	if c.Value == "" {
		return "", false
	}
	return domains.SessionId(c.Value), true
}

type sessionPayload struct {
	CSRF string `json:"csrf,omitempty"`
	User struct {
		ID       string   `json:"id"`
		Username string   `json:"username"`
		Roles    []string `json:"roles"`
	} `json:"user"`
}
type userPayload struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

func toPayload(csrf string, user *domains.User_t) sessionPayload {
	// convert roles to a slice
	var roles []string
	for role, valid := range user.Roles {
		if valid {
			roles = append(roles, string(role))
		}
	}
	if roles == nil {
		roles = []string{}
	} else {
		slices.Sort(roles)
	}
	return sessionPayload{
		CSRF: csrf,
		User: userPayload{
			ID:       fmt.Sprintf("%d", user.ID),
			Username: user.Username,
			Roles:    roles,
		},
	}
}
