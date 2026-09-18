package entra

import (
	"context"
	"time"
)

// FlowState is the server-side half of an in-progress sign-in. It is written
// by Start and consumed exactly once by Callback.
type FlowState struct {
	State     string
	Nonce     string
	Verifier  string // PKCE code verifier
	Next      string // validated post-login path, may be empty
	ExpiresAt time.Time
}

// StateStore persists FlowState between the redirect to Microsoft and the
// callback. Keys are hex SHA-256 digests of the flow cookie value; the store
// never sees the raw token. MemoryStore satisfies it; multi-instance
// deployments should back it with a shared database.
type StateStore interface {
	PutState(ctx context.Context, key string, s FlowState) error
	// TakeState returns and deletes the state in one step so a replayed
	// authorization code finds nothing. ok is false when the key is unknown.
	TakeState(ctx context.Context, key string) (s FlowState, ok bool, err error)
}

// Session is a signed-in browser. UserID is whatever the Provisioner returned;
// the user record itself is reloaded on every request.
type Session struct {
	UserID     string
	CreatedAt  time.Time
	LastSeenAt time.Time
	ExpiresAt  time.Time
}

// SessionStore persists sessions. Keys are hex SHA-256 digests of the session
// cookie value; the store never sees the raw token. MemoryStore satisfies it.
type SessionStore interface {
	PutSession(ctx context.Context, key string, s Session) error
	// GetSession returns ok=false for an unknown key. Expiry is enforced by the
	// caller, not the store.
	GetSession(ctx context.Context, key string) (s Session, ok bool, err error)
	// TouchSession slides LastSeenAt and ExpiresAt forward.
	TouchSession(ctx context.Context, key string, lastSeen, expires time.Time) error
	DeleteSession(ctx context.Context, key string) error
	// DeleteUserSessions signs the user out everywhere. Applications call
	// this through Auth.Sessions when they disable an account.
	DeleteUserSessions(ctx context.Context, userID string) error
}
