// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package rest

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/playbymail/ottoapp/backend/services/sessmgr"
)

func currentSession(r *http.Request) (*Session, bool) {
	log.Printf("%s %s: currentSession: entered\n", r.Method, r.URL.Path)
	sid, ok := readSID(r)
	if !ok {
		log.Printf("%s %s: currentSession: readSID false\n", r.Method, r.URL.Path)
		return nil, false
	}
	log.Printf("%s %s: currentSession: readSID %q\n", r.Method, r.URL.Path, sid)
	s, ok := sessions[sid]
	if !ok {
		log.Printf("%s %s: currentSession: sessions false\n", r.Method, r.URL.Path)
		return nil, false
	}
	if time.Now().After(s.Expiry) {
		log.Printf("%s %s: currentSession: sessions expired\n", r.Method, r.URL.Path)
		delete(sessions, sid)
		return nil, false
	}
	log.Printf("%s %s: currentSession: sessions %+v\n", r.Method, r.URL.Path, s)
	return s, true
}

func readSID(r *http.Request) (string, bool) {
	log.Printf("%s %s: readSID: entered\n", r.Method, r.URL.Path)
	c, err := r.Cookie("sid")
	if err != nil {
		log.Printf("%s %s: readSID: %v\n", r.Method, r.URL.Path, err)
		return "", false
	}
	log.Printf("%s %s: readSID: %q\n", r.Method, r.URL.Path, c.Value)
	if c.Value == "" {
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

func checkUser(auth sessmgr.AuthStore, handle, password string) (User, bool) {
	log.Printf("checkUser(%q, %q)\n", handle, password)
	id, err := auth.AuthenticateUser(handle, password)
	if err != nil {
		log.Printf("checkUser(%q, %q) %v\n", handle, password, err)
		return User{}, false
	}
	user, err := auth.GetUserByID(id)
	if err != nil {
		log.Printf("checkUser(%q, %q) %v\n", handle, password, err)
		return User{}, false
	}
	log.Printf("checkUser(%q, %q) %+v\n", handle, password, *user)
	return User{
		ID:       fmt.Sprintf("%d", user.ID),
		Username: user.Username,
		Role:     "guest",
	}, true
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
