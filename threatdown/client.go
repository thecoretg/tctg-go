package threatdown

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
	envVarClientID     = "THREATDOWN_CLIENT_ID"
	envVarClientSecret = "THREATDOWN_CLIENT_SECRET"
	tokenURL           = "https://api.threatdown.com/oneview/oauth2/token"
	baseURLV1          = "https://api.threatdown.com/oneview/v1"
	baseURLV2          = "https://api.threatdown.com/oneview/v2"
)

type Client struct {
	restClient *resty.Client
}

type envVars struct {
	clientID     string
	clientSecret string
}

func NewClient(ctx context.Context) (*Client, error) {
	ev := getEnvVars()
	if err := ev.validate(); err != nil {
		return nil, fmt.Errorf("validating env vars: %w", err)
	}

	ts := (&clientcredentials.Config{
		ClientID:     ev.clientID,
		ClientSecret: ev.clientSecret,
		TokenURL:     tokenURL,
		Scopes:       []string{"read", "write"},
	}).TokenSource(ctx)

	rc := resty.NewWithClient(oauth2.NewClient(ctx, ts))
	rc.SetHeader("Accept", "application/json")
	rc.SetRetryCount(3)
	rc.SetDisableWarn(true)

	return &Client{restClient: rc}, nil
}

func endpointURLV1(endpoint string) string {
	return fmt.Sprintf("%s/%s", baseURLV1, endpoint)
}

func endpointURLV2(endpoint string) string {
	return fmt.Sprintf("%s/%s", baseURLV2, endpoint)
}

func getEnvVars() envVars {
	return envVars{
		clientID:     os.Getenv(envVarClientID),
		clientSecret: os.Getenv(envVarClientSecret),
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
	if len(missing) > 0 {
		return fmt.Errorf("missing threatdown env vars: %s", strings.Join(missing, ", "))
	}
	return nil
}
