// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sessmgr

import (
	"encoding/base64"
	"encoding/binary"
	"math/rand/v2"
	"net/http"
)

// Based on the excellent articles by Mohamed Said at themsaid.com.
// * https://themsaid.com/csrf-protection-go-web-applications

func generateCSRFToken() string {
	id := make([]byte, 32)
	binary.LittleEndian.PutUint64(id[0*8:], rand.Uint64())
	binary.LittleEndian.PutUint64(id[1*8:], rand.Uint64())
	binary.LittleEndian.PutUint64(id[2*8:], rand.Uint64())
	binary.LittleEndian.PutUint64(id[3*8:], rand.Uint64())
	return base64.RawURLEncoding.EncodeToString(id)
}

func (m *SessionManager_t) verifyCSRFToken(r *http.Request, session *Session_t) bool {
	sToken, ok := session.Get("csrf_token").(string)
	if !ok {
		return false
	}

	token := r.FormValue("csrf_token")

	if token == "" {
		token = r.Header.Get("X-XSRF-Token")
	}

	return token == sToken
}

//func SpaHandlerSnippet(w http.ResponseWriter, r *http.Request) {
//	session := GetSession(r)
//
//	csrfToken := session.Get("csrf_token").(string)
//
//	cookie := &http.Cookie{
//		Name:     "XSRF-TOKEN",
//		Value:    csrfToken,
//		Domain:   "domain.com",
//		HttpOnly: true,
//		Path:     "/",
//		Secure:   true,
//		SameSite: http.SameSiteLaxMode,
//	}
//
//	http.SetCookie(w, cookie)
//}
//
