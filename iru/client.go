package iru

import (
	"fmt"

	"resty.dev/v3"
)

type (
	Config struct {
		APIKey    string
		Subdomain string // subdomain.api.kandji.io
	}

	Client struct {
		restyClient *resty.Client
	}
)

func NewClient(cfg Config) *Client {
	rc := resty.New()
	rc.SetHeader("Accept", "application/json")
	rc.SetAuthToken(cfg.APIKey)
	rc.SetRetryCount(3)
	rc.SetLoggerWarnLevel(false)
	rc.SetBaseURL(baseURL(cfg.Subdomain))

	return &Client{
		restyClient: rc,
	}
}

func baseURL(subdomain string) string {
	return fmt.Sprintf("https://%s.api.kandji.io", subdomain)
}
