package salesforce

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"resty.dev/v3"
)

const versionTag = "v64.0"

type (
	Config struct {
		ClientID       string
		ClientSecret   string
		CompanyURLName string
	}

	Client struct {
		restClient *resty.Client
		baseURL    string
	}
)

func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	var missing []string
	if cfg.ClientID == "" {
		missing = append(missing, "ClientID")
	}
	if cfg.ClientSecret == "" {
		missing = append(missing, "ClientSecret")
	}
	if cfg.CompanyURLName == "" {
		missing = append(missing, "CompanyURLName")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing salesforce config fields: %s", strings.Join(missing, ", "))
	}

	ts := (&clientcredentials.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		TokenURL:     tokenURL(cfg.CompanyURLName),
	}).TokenSource(ctx)

	rc := resty.NewWithClient(oauth2.NewClient(ctx, ts))
	rc.SetHeader("Accept", "application/json")
	rc.SetRetryCount(3)
	rc.SetLoggerWarnLevel(false)

	return &Client{restClient: rc, baseURL: baseURL(cfg.CompanyURLName)}, nil
}

func (c *Client) endpointURL(endpoint string) string {
	return fmt.Sprintf("%s/services/data/%s/%s", c.baseURL, versionTag, endpoint)
}

func baseURL(companyName string) string {
	return fmt.Sprintf("https://%s.my.salesforce.com", companyName)
}

func tokenURL(companyName string) string {
	return fmt.Sprintf("%s/%s", baseURL(companyName), "services/oauth2/token")
}
