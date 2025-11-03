// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package memdb

import (
	"sync"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
)

func NewSessionStore() (*SessionStore, error) {
	m := &SessionStore{
		sessions: make(map[domains.SessionId]*domains.Session),
	}
	return m, nil
}

type SessionStore struct {
	mu       sync.Mutex
	sessions map[domains.SessionId]*domains.Session
}

func (m *SessionStore) CreateSession(sess *domains.Session) error {
	if sess == nil {
		return domains.ErrSessionInvalid
	} else if sess.Id == "" {
		return domains.ErrSessionIdInvalid
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now().UTC()
	if !now.Before(sess.ExpiresAt) {
		delete(m.sessions, sess.Id)
		return domains.ErrSessionExpired
	}
	m.sessions[sess.Id] = sess
	return nil
}

func (m *SessionStore) ReadSession(id domains.SessionId) (*domains.Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now().UTC()
	sess, ok := m.sessions[id]
	if !ok {
		return nil, domains.ErrNotExists
	} else if !now.Before(sess.ExpiresAt) {
		delete(m.sessions, id)
		return nil, domains.ErrSessionExpired
	}
	return sess, nil
}

func (m *SessionStore) UpdateSession(sess *domains.Session) error {
	if sess == nil {
		return domains.ErrSessionInvalid
	} else if sess.Id == "" {
		return domains.ErrSessionIdInvalid
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	sess, ok := m.sessions[sess.Id]
	if !ok {
		return domains.ErrNotExists
	}
	m.sessions[sess.Id] = sess
	return nil
}

func (m *SessionStore) DeleteSession(id domains.SessionId) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, id)
	return nil
}

func (m *SessionStore) ReapSessions() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now().UTC()
	for id, sess := range m.sessions {
		if !now.Before(sess.ExpiresAt) {
			delete(m.sessions, id)
		}
	}
	return nil
}
