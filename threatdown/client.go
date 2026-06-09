package threatdown

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"resty.dev/v3"
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
	restClient *resty.Client
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

	rc := resty.NewWithClient(oauth2.NewClient(ctx, ts))
	rc.SetHeader("Accept", "application/json")
	rc.SetRetryCount(3)
	rc.SetLoggerWarnLevel(false)

	return &Client{restClient: rc}, nil
}

func endpointURLV1(endpoint string) string {
	return fmt.Sprintf("%s/%s", baseURLV1, endpoint)
}

func endpointURLV2(endpoint string) string {
	return fmt.Sprintf("%s/%s", baseURLV2, endpoint)
}
