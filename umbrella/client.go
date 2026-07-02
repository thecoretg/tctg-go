// Package umbrella is a client for the Cisco Umbrella Providers API, which
// manages the customers belonging to a provider. It authenticates with the
// OAuth2 client-credentials flow and is built on the shared internal/httpx
// transport for retries and JSON handling.
package umbrella

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/thecoretg/tctg-go/internal/httpx"
)

const (
	tokenURL = "https://api.umbrella.com/auth/v2/token"
	baseURL  = "https://api.umbrella.com/admin/v2"
	// deploymentsBaseURL is the base for the deployments APIs (policies,
	// roaming computers), which live under a different path than the
	// providers/admin endpoints.
	deploymentsBaseURL = "https://api.umbrella.com/deployments/v2"
)

type Config struct {
	ClientID     string
	ClientSecret string
}

type Client struct {
	httpClient *http.Client
}

// NewClient builds a Client using the OAuth2 client-credentials flow. The
// credential is granted the provider scopes it is entitled to, so no scopes
// are requested per token.
func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	var missing []string
	if cfg.ClientID == "" {
		missing = append(missing, "ClientID")
	}
	if cfg.ClientSecret == "" {
		missing = append(missing, "ClientSecret")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing umbrella config fields: %s", strings.Join(missing, ", "))
	}

	ts := (&clientcredentials.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		TokenURL:     tokenURL,
	}).TokenSource(ctx)

	oauthClient := oauth2.NewClient(ctx, ts)
	hc := httpx.NewClient(oauthClient.Transport, map[string]string{"Accept": "application/json"}, 3)

	return &Client{httpClient: hc}, nil
}

func NewClientFromEnv(ctx context.Context) (*Client, error) {
	return NewClient(ctx, Config{
		ClientID:     os.Getenv("UMBRELLA_CLIENT_ID"),
		ClientSecret: os.Getenv("UMBRELLA_CLIENT_SECRET"),
	})
}

func endpointURL(path string) string {
	return fmt.Sprintf("%s/%s", baseURL, path)
}

func deploymentsURL(path string) string {
	return fmt.Sprintf("%s/%s", deploymentsBaseURL, path)
}
