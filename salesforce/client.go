package salesforce

import (
	"context"
	"fmt"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"resty.dev/v3"
)

const (
	envVarClientID       = "SALESFORCE_CLIENT_ID"
	envVarClientSecret   = "SALESFORCE_CLIENT_SECRET"
	envVarCompanyURLName = "SALESFORCE_COMPANY_URL_NAME"
	versionTag           = "v64.0"
)

type (
	Client struct {
		restClient *resty.Client
		baseURL    string
	}

	envVars struct {
		clientID       string
		clientSecret   string
		companyURLName string
	}
)

func NewClient(ctx context.Context) (*Client, error) {
	ev := getEnvVars()
	if err := ev.validate(); err != nil {
		return nil, fmt.Errorf("validating env vars: %w", err)
	}

	ts := (&clientcredentials.Config{
		ClientID:     ev.clientID,
		ClientSecret: ev.clientSecret,
		TokenURL:     tokenURL(ev.companyURLName),
	}).TokenSource(ctx)

	rc := resty.NewWithClient(oauth2.NewClient(ctx, ts))
	rc.SetHeader("Accept", "application/json")
	rc.SetRetryCount(3)
	rc.SetDisableWarn(true)

	return &Client{restClient: rc, baseURL: baseURL(ev.companyURLName)}, nil
}

func getEnvVars() envVars {
	return envVars{
		clientID:       os.Getenv(envVarClientID),
		clientSecret:   os.Getenv(envVarClientSecret),
		companyURLName: os.Getenv(envVarCompanyURLName),
	}
}

func (e envVars) validate() error {
	var missing []string
	if e.clientID == "" {
		missing = append(missing, envVarClientID)
	}

	if e.clientSecret == "" {
		missing = append(missing, envVarClientSecret)
	}

	if e.companyURLName == "" {
		missing = append(missing, envVarCompanyURLName)
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing salesforce env vars: %s", strings.Join(missing, ", "))
	}

	return nil
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
