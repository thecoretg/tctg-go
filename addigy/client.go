// Package addigy is a hand-written client for the Addigy API v2
// (https://app.addigy.com/integrations). It follows the same conventions as the
// other tctg-go integrations: a static-header client, generic request helpers
// over internal/httpx, and one file per endpoint group.
//
// Authentication uses an Addigy API key sent in the x-api-key header. Most v2
// endpoints are scoped to an organization under /o/{organization_id}/...; the
// organization ID is stored on the Client and injected by orgURL.
package addigy

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/thecoretg/tctg-go/internal/httpx"
)

// baseURL is the Addigy API v2 root. The Swagger spec declares only the
// basePath (/api/v2); the v2 API is served from the api.addigy.com host (the
// app.addigy.com host returns 401 "Missing auth_token").
const baseURL = "https://api.addigy.com/api/v2"

type Config struct {
	APIKey string
	// OrgID is the organization ID injected into org-scoped endpoints
	// (/o/{organization_id}/...). It is optional for clients that only call
	// non-org-scoped endpoints.
	OrgID string
}

type Client struct {
	httpClient *http.Client
	orgID      string
}

// NewClient builds an Addigy client from cfg. The ctx parameter is accepted for
// signature consistency with the other tctg-go clients; Addigy uses a static API
// key and performs no setup that requires it.
func NewClient(_ context.Context, cfg Config) (*Client, error) {
	var missing []string
	if cfg.APIKey == "" {
		missing = append(missing, "APIKey")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing addigy config fields: %s", strings.Join(missing, ", "))
	}

	headers := map[string]string{
		"x-api-key":    cfg.APIKey,
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	return &Client{
		httpClient: httpx.NewClient(nil, headers, 3),
		orgID:      cfg.OrgID,
	}, nil
}

func NewClientFromEnv(ctx context.Context) (*Client, error) {
	return NewClient(ctx, Config{
		APIKey: os.Getenv("ADDIGY_API_KEY"),
		OrgID:  os.Getenv("ADDIGY_ORG_ID"),
	})
}

// url builds a URL for a non-org-scoped endpoint.
func (c *Client) url(path string) string {
	return fmt.Sprintf("%s/%s", baseURL, path)
}

// orgURL builds a URL for an org-scoped endpoint (/o/{organization_id}/...). It
// errors if the client was created without an OrgID.
func (c *Client) orgURL(path string) (string, error) {
	if c.orgID == "" {
		return "", fmt.Errorf("addigy: OrgID is required for org-scoped endpoints")
	}
	return fmt.Sprintf("%s/o/%s/%s", baseURL, c.orgID, path), nil
}
