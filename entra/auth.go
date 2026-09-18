package entra

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// Auth holds the configuration and discovered provider metadata for one app
// registration. Obtain one with New or NewFromEnv, mount Start, Callback and
// Logout, and wrap protected handlers with RequireAuth.
type Auth[U any] struct {
	cfg  Config
	prov Provisioner[U]

	mu       sync.Mutex
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
}

// RedirectURI is the exact value to register under Authentication > Web in
// the app registration.
func (a *Auth[U]) RedirectURI() string {
	return a.cfg.BaseURL + a.cfg.CallbackPath
}

// Sessions exposes the session store so the application can revoke sessions,
// for example with DeleteUserSessions when it disables an account.
func (a *Auth[U]) Sessions() SessionStore {
	return a.cfg.Sessions
}

// discover fetches the OIDC discovery document once and caches the provider
// and verifier. A failed attempt is not cached, so a transient Microsoft
// outage during the first sign-in does not poison later ones. The verifier's
// JWKS cache refetches on an unknown kid, which is all key rollover needs.
func (a *Auth[U]) discover(ctx context.Context) (*oidc.Provider, *oidc.IDTokenVerifier, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.provider != nil {
		return a.provider, a.verifier, nil
	}

	p, err := oidc.NewProvider(oidc.ClientContext(ctx, a.cfg.HTTPClient), a.cfg.issuer)
	if err != nil {
		return nil, nil, fmt.Errorf("entra: discovery: %w", err)
	}

	// JWKS fetches happen after the request that triggered discovery has
	// ended, so they get a background context bound to our HTTP client.
	a.provider = p
	a.verifier = p.VerifierContext(oidc.ClientContext(context.Background(), a.cfg.HTTPClient), &oidc.Config{
		ClientID: a.cfg.ClientID,
	})
	return a.provider, a.verifier, nil
}

// oauthConfig returns an oauth2.Config with the discovered endpoints filled in.
func (a *Auth[U]) oauthConfig(ctx context.Context) (*oauth2.Config, *oidc.IDTokenVerifier, error) {
	p, v, err := a.discover(ctx)
	if err != nil {
		return nil, nil, err
	}
	return &oauth2.Config{
		ClientID:     a.cfg.ClientID,
		ClientSecret: a.cfg.ClientSecret,
		Endpoint:     p.Endpoint(),
		RedirectURL:  a.RedirectURI(),
		Scopes:       a.cfg.Scopes,
	}, v, nil
}

// httpContext makes oauth2 use the configured HTTP client.
func (a *Auth[U]) httpContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, oauth2.HTTPClient, a.cfg.HTTPClient)
}

// TestConnection proves the tenant is reachable and the client secret is
// accepted by performing discovery and a client-credentials grant for
// https://graph.microsoft.com/.default. It cannot prove the redirect URI or
// group configuration; only a real sign-in does that. A rejected secret
// returns an error wrapping ErrInvalidClient.
func (a *Auth[U]) TestConnection(ctx context.Context) error {
	p, _, err := a.discover(ctx)
	if err != nil {
		return err
	}

	cc := clientcredentials.Config{
		ClientID:     a.cfg.ClientID,
		ClientSecret: a.cfg.ClientSecret,
		TokenURL:     p.Endpoint().TokenURL,
		Scopes:       []string{"https://graph.microsoft.com/.default"},
	}
	if _, err := cc.Token(a.httpContext(ctx)); err != nil {
		return fmt.Errorf("entra: test connection: %w", mapTokenError(err))
	}
	return nil
}

// mapTokenError turns an oauth2.RetrieveError into an *Error carrying the
// Entra code, wrapping ErrInvalidClient where that is the code.
func mapTokenError(err error) error {
	re, ok := errors.AsType[*oauth2.RetrieveError](err)
	if !ok || re.ErrorCode == "" {
		return err
	}
	e := &Error{Code: re.ErrorCode, Description: re.ErrorDescription}
	if re.ErrorCode == "invalid_client" {
		e.Err = ErrInvalidClient
	}
	return e
}
