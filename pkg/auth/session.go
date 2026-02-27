package auth

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type Session struct {
	User         *UserInfo
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

func NewSessionStore() *SessionStore {
	ss := &SessionStore{
		sessions: make(map[string]*Session),
	}
	go ss.cleanup()
	return ss
}

func (ss *SessionStore) Create(user *UserInfo, accessToken, refreshToken string, expiresIn time.Duration) string {
	id := generateSessionID()
	ss.mu.Lock()
	ss.sessions[id] = &Session{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(expiresIn),
		CreatedAt:    time.Now(),
	}
	ss.mu.Unlock()
	return id
}

func (ss *SessionStore) Get(id string) *Session {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	s, ok := ss.sessions[id]
	if !ok {
		return nil
	}
	if time.Since(s.CreatedAt) > 24*time.Hour {
		return nil
	}
	return s
}

func (ss *SessionStore) Delete(id string) {
	ss.mu.Lock()
	delete(ss.sessions, id)
	ss.mu.Unlock()
}

func (ss *SessionStore) UpdateTokens(id, accessToken, refreshToken string, expiresIn time.Duration) {
	ss.mu.Lock()
	if s, ok := ss.sessions[id]; ok {
		s.AccessToken = accessToken
		s.RefreshToken = refreshToken
		s.ExpiresAt = time.Now().Add(expiresIn)
	}
	ss.mu.Unlock()
}

func (ss *SessionStore) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		ss.mu.Lock()
		now := time.Now()
		for id, s := range ss.sessions {
			if now.Sub(s.CreatedAt) > 24*time.Hour {
				delete(ss.sessions, id)
			}
		}
		ss.mu.Unlock()
	}
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
