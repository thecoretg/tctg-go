package entra

import "errors"

// Sentinel errors. Handlers pass these (possibly wrapped) to Config.OnError so
// applications can distinguish outcomes with errors.Is.
var (
	// ErrAccessDenied means the user cancelled or declined at Microsoft.
	ErrAccessDenied = errors.New("entra: sign-in cancelled")
	// ErrInvalidClient means the token endpoint rejected the client
	// credentials. Almost always the client secret has expired (AADSTS7000222).
	ErrInvalidClient = errors.New("entra: client credentials rejected")
	// ErrFlowExpired means the callback arrived without a flow cookie, with an
	// unknown or already-consumed flow (a replay), or after the flow's TTL.
	ErrFlowExpired = errors.New("entra: sign-in flow expired or missing")
	// ErrStateMismatch means the state or nonce did not match the flow.
	ErrStateMismatch = errors.New("entra: state or nonce mismatch")
	// ErrWrongTenant means the ID token's tid differs from Config.TenantID.
	ErrWrongTenant = errors.New("entra: token issued by a different tenant")
	// ErrNotAuthorized means the Authorizer refused the sign-in. Errors from
	// the shipped Authorizers are *DeniedError values with a user-safe Message.
	ErrNotAuthorized = errors.New("entra: not authorized")
	// ErrUserDisabled is returned by a Provisioner to refuse a disabled account
	// at sign-in, or to end an existing session on the next request.
	ErrUserDisabled = errors.New("entra: user disabled")
	// ErrUserNotFound is returned by Provisioner.Lookup when the session's user
	// no longer exists; the session is deleted.
	ErrUserNotFound = errors.New("entra: user not found")
	// ErrNoOID means the ID token carried no oid claim, which usually means the
	// profile scope is missing from the app registration.
	ErrNoOID = errors.New("entra: id token has no oid claim")
)

// Error carries the error code and description Entra returned, either on the
// authorization response (?error=...) or from the token endpoint. Err is the
// sentinel it maps to, if any.
type Error struct {
	Code        string
	Description string
	Err         error
}

func (e *Error) Error() string {
	s := "entra: " + e.Code
	if e.Description != "" {
		s += ": " + e.Description
	}
	return s
}

func (e *Error) Unwrap() error { return e.Err }

// DeniedError is returned by the shipped Authorizers. Message is written for
// the person signing in and carries no tenant or group identifiers, so it is
// safe to show in a browser.
type DeniedError struct {
	Message string
}

func (e *DeniedError) Error() string { return "entra: not authorized: " + e.Message }

// Is reports true for ErrNotAuthorized so errors.Is works without knowing the
// concrete type.
func (e *DeniedError) Is(target error) bool { return target == ErrNotAuthorized }
