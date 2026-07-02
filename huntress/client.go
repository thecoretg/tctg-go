// Package huntress is a client for the Huntress API, which manages
// organizations, agents, users, escalations, identities, and reseller-level
// accounts. It authenticates with HTTP Basic auth (a base64-encoded API key
// and secret) and is built on the shared internal/httpx transport for retries
// and JSON handling.
package huntress

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/thecoretg/tctg-go/internal/httpx"
)

const baseURL = "https://api.huntress.io/v1"

type Config struct {
	APIKey    string
	APISecret string
}

type Client struct {
	httpClient *http.Client
}

// NewClient builds a Client that authenticates every request with HTTP Basic
// auth using the API key and secret. ctx is accepted for signature parity with
// the other tctg-go clients; it is not used here.
func NewClient(_ context.Context, cfg Config) (*Client, error) {
	var missing []string
	if cfg.APIKey == "" {
		missing = append(missing, "APIKey")
	}
	if cfg.APISecret == "" {
		missing = append(missing, "APISecret")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing huntress config fields: %s", strings.Join(missing, ", "))
	}

	cred := base64.StdEncoding.EncodeToString([]byte(cfg.APIKey + ":" + cfg.APISecret))
	headers := map[string]string{
		"Authorization": "Basic " + cred,
		"Accept":        "application/json",
	}

	return &Client{httpClient: httpx.NewClient(nil, headers, 3)}, nil
}

func NewClientFromEnv(ctx context.Context) (*Client, error) {
	return NewClient(ctx, Config{
		APIKey:    os.Getenv("HUNTRESS_API_KEY"),
		APISecret: os.Getenv("HUNTRESS_API_SECRET"),
	})
}

func endpointURL(path string) string {
	return fmt.Sprintf("%s/%s", baseURL, path)
}
