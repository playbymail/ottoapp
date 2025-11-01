// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package sessmgr defines a session manager interface and
// provides an in-memory store as a reference implementation.
//
// This implementation uses math/rand/v2 for simplicity and reproducibility,
// not for security. In production, replace it with crypto/rand and consider
// using a standard session or token library instead.
package sessmgr

import (
	"fmt"
	"net/http"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
)

// Based on the excellent articles by Mohamed Said at themsaid.com.
// * https://themsaid.com/session-authentication-go

// Client                                    Server
//  |                                          |
//  |  1. User visits the web app              |
//  |----------------------------------------> |
//  |                                          |
//  |  2. Server creates a session and sends   |
//  |     its ID back                          |
//  | <--------------------------------------- |
//  |                                          |
//  |  3. User submits login form              |
//  |----------------------------------------> |
//  |                                          |
//  |  4. Server verifies credentials and      |
//  |     stores the user identifier in the    |
//  |     session                              |
//  |                                          |
//  |  5. Server sends a success response      |
//  | <----------------------------------------|
//  |                                          |
//  |  6. Client makes authenticated requests  |
//  |     (Cookie is automatically sent)       |
//  |----------------------------------------> |
//  |                                          |
//  |  7. Server retrieves session data and    |
//  |     extracts the user identifier         |
//  |                                          |
//  |  8. Server processes request and         |
//  |     responds                             |
//  | <----------------------------------------|

type AuthStore interface {
	AuthenticateUser(handle, plainTextSecret string) (domains.ID, error)
	CreateUser(handle, email, plainTextSecret string, timezone *time.Location) (*domains.User_t, error)
	GetUserByHandle(handle string) (*domains.User_t, error)
	GetUserByID(id domains.ID) (*domains.User_t, error)
}

func Register(
	auth AuthStore,
	handle string,
	email string,
	password string,
	timezone *time.Location,
) error {
	_, err := auth.CreateUser(handle, email, password, timezone)
	if err != nil {
		return err
	}
	return nil
}

func VerifyCredentials(
	auth AuthStore,
	handle string,
	password string,
) error {
	_, err := auth.AuthenticateUser(handle, password)
	if err != nil {
		return err
	}
	return nil
}

func Login(
	r *http.Request,
	sm *SessionManager,
	handle string,
) error {
	session := GetSession(r)

	err := sm.migrate(session)
	if err != nil {
		return fmt.Errorf("failed to migrate session: %w", err)
	}

	session.Put("handle", handle)

	return nil
}

func Logout(
	r *http.Request,
	sm *SessionManager,
) error {
	session := GetSession(r)

	err := sm.migrate(session)
	if err != nil {
		return fmt.Errorf("failed to migrate session: %w", err)
	}

	session.Put("handle", "")

	return nil
}

func Auth(auth AuthStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := GetSession(r)

		handle := session.Get("handle").(string)
		if handle == "" {
			http.Error(w, "Unauthenticated", http.StatusForbidden)
			return
		}

		_, err := auth.GetUserByHandle(handle)
		if err != nil {
			// return a 403, but could redirect to /login form
			http.Error(w, "Unauthenticated", http.StatusForbidden)
			return
		}
	})
}

// LoginWithDelay to avoid timing attack
func LoginWithDelay(w http.ResponseWriter, r *http.Request, auth AuthStore, sm *SessionManager, minDelay time.Duration) {
	//timingAttackGuard := time.NewTimer(500 * time.Microsecond)
	timingAttackGuard := time.NewTimer(minDelay)

	err := VerifyCredentials(auth, r.FormValue("username"), r.FormValue("password"))
	if err == nil {
		err = Login(r, sm, r.FormValue("username"))
	}

	// delay until the timer expires to prevent timing attacks
	<-timingAttackGuard.C

	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// example authorization
// session := GetSession(r)
// username := session.Get("username").(string)
// if username == "" {
//   // No user (un-authenticated session)
// } else {
//   // We have a user (authenticated session)
// }
