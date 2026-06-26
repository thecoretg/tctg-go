// Package automate is a hand-written client for the ConnectWise Automate REST
// API (/cwa/api/v1). It follows the same conventions as the other tctg-go
// integrations: a static-header client, generic request helpers over
// internal/httpx, and one file per resource group.
//
// Automate is deployed per-tenant, so the server base URL is supplied via
// config. Authentication uses a bearer token obtained by POSTing Automate
// credentials to /cwa/api/v1/apitoken; every request also carries a registered
// integrator ClientId header. See auth.go for the login/refresh flow.
package automate

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/thecoretg/tctg-go/internal/httpx"
)

type Config struct {
	// ServerURL is the root of the Automate server, e.g.
	// https://companyABC.hostedrmm.com. The /cwa/api/v1 path is part of each
	// endpoint.
	ServerURL string `json:"server_url,omitempty" mapstructure:"server_url"`
	// ClientID is the integrator GUID registered with the ConnectWise Developer
	// Network, sent in the ClientId header on every request.
	ClientID string `json:"client_id,omitempty" mapstructure:"client_id"`

	// Username and Password are the Automate credentials used to obtain a bearer
	// token. Required unless Token is supplied.
	Username string `json:"username,omitempty" mapstructure:"username"`
	Password string `json:"password,omitempty" mapstructure:"password"`
	// TwoFactorPasscode is required only when the Automate user has two-factor
	// authentication enabled.
	TwoFactorPasscode string `json:"two_factor_passcode,omitempty" mapstructure:"two_factor_passcode"`

	// Token is an optional pre-issued bearer token. When set, NewClient skips the
	// username/password login and uses it directly.
	Token string `json:"token,omitempty" mapstructure:"token"`
}

type Client struct {
	httpClient *http.Client
	baseURL    string
	clientID   string

	username  string
	password  string
	twoFactor string

	mu          sync.RWMutex
	token       string
	tokenResult *TokenResult
}

// NewClient builds a ConnectWise Automate client from cfg. Unless cfg.Token is
// supplied, it logs in with the configured username/password to obtain a bearer
// token, so the provided ctx governs that initial request.
func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	var missing []string
	if cfg.ServerURL == "" {
		missing = append(missing, "ServerURL")
	}
	if cfg.ClientID == "" {
		missing = append(missing, "ClientID")
	}
	if cfg.Token == "" && (cfg.Username == "" || cfg.Password == "") {
		missing = append(missing, "Token or Username+Password")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing automate config fields: %s", strings.Join(missing, ", "))
	}

	c := &Client{
		baseURL:   strings.TrimRight(cfg.ServerURL, "/"),
		clientID:  cfg.ClientID,
		username:  cfg.Username,
		password:  cfg.Password,
		twoFactor: cfg.TwoFactorPasscode,
		token:     cfg.Token,
	}

	// The bearer token is injected dynamically (it changes on login/refresh);
	// the ClientId header is constant and the retry/JSON headers come from httpx.
	base := &authTransport{base: http.DefaultTransport, clientID: cfg.ClientID, token: c.currentToken}
	c.httpClient = httpx.NewClient(base, map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}, 3)

	if cfg.Token == "" {
		if _, err := c.Login(ctx); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func NewClientFromEnv(ctx context.Context) (*Client, error) {
	return NewClient(ctx, Config{
		ServerURL:         os.Getenv("CWA_SERVER_URL"),
		ClientID:          os.Getenv("CWA_CLIENT_ID"),
		Username:          os.Getenv("CWA_USERNAME"),
		Password:          os.Getenv("CWA_PASSWORD"),
		TwoFactorPasscode: os.Getenv("CWA_2FA"),
		Token:             os.Getenv("CWA_TOKEN"),
	})
}

func (c *Client) fullURL(endpoint string) string {
	return fmt.Sprintf("%s/%s", c.baseURL, endpoint)
}
