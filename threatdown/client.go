package threatdown

import (
	"context"
	"fmt"
	"os"
	"strings"

	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/thecoretg/tctg-go/internal/httpx"
)

const (
	tokenURL  = "https://api.threatdown.com/oneview/oauth2/token"
	baseURLV1 = "https://api.threatdown.com/oneview/v1"
	baseURLV2 = "https://api.threatdown.com/oneview/v2"
)

type Config struct {
	ClientID     string
	ClientSecret string
}

type Client struct {
	httpClient *http.Client
}

func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	var missing []string
	if cfg.ClientID == "" {
		missing = append(missing, "ClientID")
	}
	if cfg.ClientSecret == "" {
		missing = append(missing, "ClientSecret")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing threatdown config fields: %s", strings.Join(missing, ", "))
	}

	ts := (&clientcredentials.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		TokenURL:     tokenURL,
		Scopes:       []string{"read", "write"},
	}).TokenSource(ctx)

	oauthClient := oauth2.NewClient(ctx, ts)
	hc := httpx.NewClient(oauthClient.Transport, map[string]string{"Accept": "application/json"}, 3)

	return &Client{httpClient: hc}, nil
}

func NewClientFromEnv(ctx context.Context) (*Client, error) {
	return NewClient(ctx, Config{
		ClientID:     os.Getenv("THREATDOWN_CLIENT_ID"),
		ClientSecret: os.Getenv("THREATDOWN_CLIENT_SECRET"),
	})
}

func endpointURLV1(endpoint string) string {
	return fmt.Sprintf("%s/%s", baseURLV1, endpoint)
}

func endpointURLV2(endpoint string) string {
	return fmt.Sprintf("%s/%s", baseURLV2, endpoint)
}
