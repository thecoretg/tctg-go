package automate

import (
	"context"
	"fmt"
	"net/http"
)

// TokenCredentials is the body POSTed to /cwa/api/v1/apitoken
// (Automate.Api.Domain.Contracts.Security.TokenCredentials).
type TokenCredentials struct {
	Username          string `json:"Username"`
	Password          string `json:"Password"`
	TwoFactorPasscode string `json:"TwoFactorPasscode,omitempty"`
}

// TokenResult is the response from the apitoken endpoints
// (Automate.Api.Domain.Contracts.Security.TokenResult).
type TokenResult struct {
	AccessToken                 string `json:"AccessToken"`
	TokenType                   string `json:"TokenType"`
	ExpirationDate              Time   `json:"ExpirationDate,omitzero"`
	AbsoluteExpirationDate      Time   `json:"AbsoluteExpirationDate,omitzero"`
	UserId                      string `json:"UserId"`
	InternalUserName            string `json:"InternalUserName"`
	IsTwoFactorRequired         bool   `json:"IsTwoFactorRequired"`
	IsInternalTwoFactorRequired bool   `json:"IsInternalTwoFactorRequired"`
	SSOAccessToken              string `json:"SSOAccessToken"`
}

// authTransport injects the Automate auth headers on every request: the static
// ClientId integrator GUID and, once available, the current bearer token. The
// token is read through a func so refreshes are picked up without rebuilding the
// transport.
type authTransport struct {
	base     http.RoundTripper
	clientID string
	token    func() string
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	if t.clientID != "" {
		req.Header.Set("ClientId", t.clientID)
	}
	if tok := t.token(); tok != "" && req.Header.Get("Authorization") == "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}
	return t.base.RoundTrip(req)
}

func (c *Client) currentToken() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.token
}

func (c *Client) setToken(tr *TokenResult) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token = tr.AccessToken
	c.tokenResult = tr
}

// Login exchanges the configured username/password (and optional two-factor
// passcode) for a bearer token via POST /cwa/api/v1/apitoken, storing it for
// subsequent requests. It is called automatically by NewClient unless a token
// was supplied.
func (c *Client) Login(ctx context.Context) (*TokenResult, error) {
	creds := TokenCredentials{
		Username:          c.username,
		Password:          c.password,
		TwoFactorPasscode: c.twoFactor,
	}
	result, err := post[TokenResult](ctx, c, "cwa/api/v1/apitoken", creds)
	if err != nil {
		return nil, fmt.Errorf("automate login: %w", err)
	}
	if result.AccessToken == "" {
		if result.IsTwoFactorRequired {
			return nil, fmt.Errorf("automate login: two-factor passcode required")
		}
		return nil, fmt.Errorf("automate login: no access token returned")
	}
	c.setToken(result)
	return result, nil
}

// Refresh extends the current session via POST /cwa/api/v1/apitoken/refresh,
// returning a new bearer token valid until the original AbsoluteExpirationDate.
func (c *Client) Refresh(ctx context.Context) (*TokenResult, error) {
	token := c.currentToken()
	if token == "" {
		return nil, fmt.Errorf("automate refresh: no token to refresh")
	}
	// The refresh body is the current token as a bare JSON string.
	result, err := post[TokenResult](ctx, c, "cwa/api/v1/apitoken/refresh", token)
	if err != nil {
		return nil, fmt.Errorf("automate refresh: %w", err)
	}
	if result.AccessToken == "" {
		return nil, fmt.Errorf("automate refresh: no access token returned")
	}
	c.setToken(result)
	return result, nil
}

// Token returns the bearer token currently in use.
func (c *Client) Token() string {
	return c.currentToken()
}
