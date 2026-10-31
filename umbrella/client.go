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
	"time"

	"golang.org/x/time/rate"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/thecoretg/tctg-go/internal/httpx"
	"github.com/thecoretg/tctg-go/ratelimit"
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
	// Limiter paces requests. nil uses DefaultLimiter(); ratelimit.None turns
	// pacing off. Share one limiter between clients on the same API key.
	Limiter ratelimit.Limiter
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

	lim := cfg.Limiter
	if lim == nil {
		lim = DefaultLimiter()
	}

	oauthClient := oauth2.NewClient(ctx, ts)
	hc := httpx.NewClient(oauthClient.Transport, map[string]string{"Accept": "application/json"}, httpx.Policy{Retries: 3, Limiter: lim})

	return &Client{httpClient: hc}, nil
}

// DefaultLimiter paces a client to Umbrella's three per-key limits, which are
// enforced together: 5 per second, 14 per minute and 350 per 30 minutes. One
// token every 30m/348 with bursts of 2 satisfies all three at once: at most 2
// in a second, 13.6 in a minute and 350 in 30 minutes. That is a steady
// request every ~5.2s, the 30-minute window's sustained rate. Each call returns
// a new limiter.
func DefaultLimiter() *rate.Limiter {
	return rate.NewLimiter(rate.Every(30*time.Minute/348), 2)
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
