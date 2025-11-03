// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package sqlite

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"log"
	"math/rand/v2"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
)

func (db *DB) CreateSession(userId domains.ID, ttl time.Duration) (*domains.Session, error) {
	//log.Printf("createSession: user %d: ttl %v: expires %v\n", userId, ttl, time.Now().Add(ttl).UTC())
	user, err := db.q.GetUser(db.ctx, int64(userId))
	if err != nil {
		return nil, errors.Join(domains.ErrNotExists, err)
	}
	loc, ok := NormalizeTimeZone(user.Timezone)
	if !ok {
		return nil, domains.ErrInvalidTimezone
	}
	roles := map[string]bool{}
	if userRoles, err := db.q.GetUserRoles(db.ctx, int64(userId)); err == nil {
		for _, role := range userRoles {
			roles[role] = true
		}
	}
	sess := &domains.Session{
		Id:   newSessionId(),
		Csrf: newCsrf(),
		User: domains.User_t{
			ID:       userId,
			Username: user.Handle,
			Email:    user.Email,
			Locale: domains.UserLocale_t{
				DateFormat: "2006-01-02",
				Timezone: domains.UserTimezone_t{
					Location: loc,
				},
			},
			Roles:   roles,
			Created: user.CreatedAt,
			Updated: user.UpdatedAt,
		},
		ExpiresAt:      time.Now().Add(ttl).UTC(),
		LastActivityAt: time.Now().UTC(),
	}

	err = db.q.CreateSession(db.ctx, sqlc.CreateSessionParams{
		SessionID: string(sess.Id),
		Csrf:      sess.Csrf,
		UserID:    int64(sess.User.ID),
		ExpiresAt: sess.ExpiresAt.Unix(),
	})
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func (db *DB) DeleteSession(id domains.SessionId) error {
	err := db.q.DeleteSession(db.ctx, string(id))
	if err != nil {
		log.Printf("[sqldb] deleteSession %v\n", err)
	}
	return nil
}

func (db *DB) ReadSession(sessionId domains.SessionId) (*domains.Session, error) {
	//log.Printf("readSession %q\n", sessionId)
	ss, err := db.q.GetSession(db.ctx, string(sessionId))
	if err != nil {
		//log.Printf("readSession %q: getSession %v\n", sessionId, err)
		return nil, errors.Join(domains.ErrNotExists, err)
	} else if !time.Now().UTC().Before(time.Unix(ss.ExpiresAt, 0)) {
		//log.Printf("readSession %q: expired %v\n", sessionId, ss.ExpiresAt)
		//log.Printf("readSession %q: not is  %v\n", sessionId, time.Now().UTC())
		return nil, domains.ErrNotExists
	}
	user, err := db.q.GetUser(db.ctx, ss.UserID)
	if err != nil {
		return nil, errors.Join(domains.ErrNotExists, err)
	}
	loc, ok := NormalizeTimeZone(user.Timezone)
	if !ok {
		return nil, domains.ErrInvalidTimezone
	}
	roles := map[string]bool{}
	if userRoles, err := db.q.GetUserRoles(db.ctx, ss.UserID); err == nil {
		for _, role := range userRoles {
			roles[role] = true
		}
	}
	sess := &domains.Session{
		Id:   sessionId,
		Csrf: ss.Csrf,
		User: domains.User_t{
			ID:       domains.ID(ss.UserID),
			Username: user.Handle,
			Email:    user.Email,
			Locale: domains.UserLocale_t{
				DateFormat: "2006-01-02",
				Timezone: domains.UserTimezone_t{
					Location: loc,
				},
			},
			Roles:   roles,
			Created: user.CreatedAt,
			Updated: user.UpdatedAt,
		},
		ExpiresAt:      time.Unix(ss.ExpiresAt, 0),
		LastActivityAt: time.Now().UTC(),
	}
	return sess, nil
}

func (db *DB) ReapSessions() error {
	return db.q.ReapSessions(db.ctx, time.Now().UTC().Unix())
}

func (db *DB) SessionStoreID() string {
	return "backend/stores/sqlite:DB"
}

func (db *DB) UpdateSession(session *domains.Session) error {
	if session == nil {
		return domains.ErrSessionInvalid
	}
	return db.q.UpdateSession(db.ctx, sqlc.UpdateSessionParams{
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
