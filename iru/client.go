package iru

import (
	"fmt"
	"net/http"

	"github.com/thecoretg/tctg-go/internal/httpx"
)

type (
	Config struct {
		APIKey    string
		Subdomain string // subdomain.api.kandji.io
	}

	Client struct {
		httpClient *http.Client
		baseURL    string
	}
)

func NewClient(cfg Config) *Client {
	headers := map[string]string{
		"Accept":        "application/json",
		"Authorization": "Bearer " + cfg.APIKey,
	}

	return &Client{
		httpClient: httpx.NewClient(nil, headers, 3),
		baseURL:    baseURL(cfg.Subdomain),
	}
}

func baseURL(subdomain string) string {
	return fmt.Sprintf("https://%s.api.kandji.io", subdomain)
}
