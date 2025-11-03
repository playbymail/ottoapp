// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sessions

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math/rand/v2"
	"net/http"
	"slices"

	"github.com/playbymail/ottoapp/backend/domains"
)

func authenticateCredentials(s AuthStore, username, password string) (domains.User_t, bool) {
	id, err := s.AuthenticateUser(username, password)
	if err != nil {
		//log.Printf("authenticateCredentials(%q, %q) %v\n", username, password, err)
		return domains.User_t{}, false
	}
	user, err := s.GetUserByID(id)
	if err != nil {
		//log.Printf("authenticateCredentials(%q, %q) %v\n", username, password, err)
		return domains.User_t{}, false
	}
	//log.Printf("authenticateCredentials(%q, %q) %+v\n", username, password, *user)
	return domains.User_t{
		ID:       user.ID,
		Username: user.Username,
		Roles:    map[string]bool{"authenticated": true},
	}, true
}

func newCSRF() string {
	id := make([]byte, 32)
	binary.LittleEndian.PutUint64(id[0*8:], rand.Uint64())
	binary.LittleEndian.PutUint64(id[1*8:], rand.Uint64())
	binary.LittleEndian.PutUint64(id[2*8:], rand.Uint64())
	binary.LittleEndian.PutUint64(id[3*8:], rand.Uint64())
	return base64.RawURLEncoding.EncodeToString(id)
}

func newSID() domains.SessionId {
	id := make([]byte, 32)
	binary.LittleEndian.PutUint64(id[0*8:], rand.Uint64())
	binary.LittleEndian.PutUint64(id[1*8:], rand.Uint64())
	binary.LittleEndian.PutUint64(id[2*8:], rand.Uint64())
	binary.LittleEndian.PutUint64(id[3*8:], rand.Uint64())
	return domains.SessionId(base64.RawURLEncoding.EncodeToString(id))
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

func toPayload(csrf string, user domains.User_t) sessionPayload {
	// convert roles to a slice
	var roles []string
	for k, v := range user.Roles {
		if v {
			roles = append(roles, k)
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
