// Copyright (c) 2025 Michael D Henderson. All rights reserved.

// Package sessions implements a very simple session manager.
package sessions

import (
	"net/http"

	"github.com/playbymail/ottoapp/backend/domains"
)

/*
Sessions:what are the 4 routes for?

  POST /api/login → create session + Set-Cookie

  POST /api/logout → delete session + clear cookie

  GET /api/session → “is this cookie valid? give me csrf + user”

  GET /api/me → (optional) “give me just user again”
*/

type AuthStore_i interface {
	AuthenticateUser(username, password string) (domains.ID, error)
	GetUserByID(domains.ID) (*domains.User_t, error)
}

type SessionManager_i interface {
	DeleteCookie(w http.ResponseWriter, r *http.Request)
	GetCurrentSession(r *http.Request) (*Session_t, bool)
	GetMeHandler(w http.ResponseWriter, r *http.Request)
	GetSessionHandler(w http.ResponseWriter, r *http.Request)
	PostLoginHandler(w http.ResponseWriter, r *http.Request)
	PostLogoutHandler(w http.ResponseWriter, r *http.Request)
}
