package entra

import (
	"context"
	"maps"
	"sync"
	"time"
)

const defaultMaxStates = 1000

// MemoryStore is an in-process StateStore and SessionStore. It is safe for
// concurrent use and suitable for a single-instance deployment; everything is
// lost on restart and nothing is shared between replicas.
type MemoryStore struct {
	mu        sync.Mutex
	states    map[string]FlowState
	sessions  map[string]Session
	maxStates int
	now       func() time.Time
}

// NewMemoryStore returns an empty store that keeps at most 1000 pending
// sign-in flows, evicting the soonest-to-expire when full.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		states:    make(map[string]FlowState),
		sessions:  make(map[string]Session),
		maxStates: defaultMaxStates,
		now:       time.Now,
	}
}

func (m *MemoryStore) PutState(_ context.Context, key string, s FlowState) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := m.now()
	maps.DeleteFunc(m.states, func(_ string, fs FlowState) bool { return !fs.ExpiresAt.After(now) })

	for len(m.states) >= m.maxStates {
		var oldestKey string
		var oldest time.Time
		for k, fs := range m.states {
			if oldestKey == "" || fs.ExpiresAt.Before(oldest) {
				oldestKey, oldest = k, fs.ExpiresAt
			}
		}
		delete(m.states, oldestKey)
	}

	m.states[key] = s
	return nil
}

func (m *MemoryStore) TakeState(_ context.Context, key string) (FlowState, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.states[key]
	if ok {
		delete(m.states, key)
	}
	return s, ok, nil
}

func (m *MemoryStore) PutSession(_ context.Context, key string, s Session) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := m.now()
	maps.DeleteFunc(m.sessions, func(_ string, ss Session) bool { return !ss.ExpiresAt.After(now) })
	m.sessions[key] = s
	return nil
}

func (m *MemoryStore) GetSession(_ context.Context, key string) (Session, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.sessions[key]
	return s, ok, nil
}

func (m *MemoryStore) TouchSession(_ context.Context, key string, lastSeen, expires time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.sessions[key]
	if !ok {
		return nil
	}
	s.LastSeenAt = lastSeen
	s.ExpiresAt = expires
	m.sessions[key] = s
	return nil
}

func (m *MemoryStore) DeleteSession(_ context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.sessions, key)
	return nil
}

func (m *MemoryStore) DeleteUserSessions(_ context.Context, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	maps.DeleteFunc(m.sessions, func(_ string, s Session) bool { return s.UserID == userID })
	return nil
}
