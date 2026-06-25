package webex

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/thecoretg/tctg-go/internal/httpx"
)

type Config struct {
	Token string
}

type Client struct {
	httpClient *http.Client
}

// NewClient builds a Webex client from cfg. The ctx parameter is accepted for
// signature consistency with the other tctg-go clients; Webex uses a static
// bearer token and performs no setup that requires it.
func NewClient(_ context.Context, cfg Config) (*Client, error) {
	var missing []string
	if cfg.Token == "" {
		missing = append(missing, "Token")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing webex config fields: %s", strings.Join(missing, ", "))
	}

	headers := map[string]string{
		"Authorization": "Bearer " + cfg.Token,
		"Content-Type":  "application/json",
		"Accept":        "application/json",
	}

	return &Client{httpClient: httpx.NewClient(nil, headers, 3)}, nil
}

func NewClientFromEnv(ctx context.Context) (*Client, error) {
	return NewClient(ctx, Config{
		Token: os.Getenv("WEBEX_TOKEN"),
	})
}
