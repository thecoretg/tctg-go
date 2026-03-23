package rewst

import (
	"fmt"
	"os"
	"strings"

	"resty.dev/v3"
)

const (
	envVarSecret       = "REWST_WEBHOOK_SECRET"
	envVarListOrgsURL  = "REWST_LIST_ORGS_URL"
	envVarUpsertVarURL = "REWST_UPSERT_ORG_VAR_URL"
)

type Client struct {
	restClient   *resty.Client
	listOrgsURL  string
	upsertVarURL string
}

type envVars struct {
	secret       string
	listOrgsURL  string
	upsertVarURL string
}

func NewClient() (*Client, error) {
	ev := getEnvVars()
	if err := ev.validate(); err != nil {
		return nil, fmt.Errorf("validating env vars: %w", err)
	}

	rc := resty.New()
	rc.SetHeader("Accept", "application/json")
	rc.SetHeader("x-rewst-secret", ev.secret)
	rc.SetRetryCount(3)
	rc.SetDisableWarn(true)

	return &Client{
		restClient:   rc,
		listOrgsURL:  ev.listOrgsURL,
		upsertVarURL: ev.upsertVarURL,
	}, nil
}

func getEnvVars() envVars {
	return envVars{
		secret:       os.Getenv(envVarSecret),
		listOrgsURL:  os.Getenv(envVarListOrgsURL),
		upsertVarURL: os.Getenv(envVarUpsertVarURL),
	}
}

func (e envVars) validate() error {
	var missing []string
	if e.secret == "" {
		missing = append(missing, envVarSecret)
	}
	if e.listOrgsURL == "" {
		missing = append(missing, envVarListOrgsURL)
	}
	if e.upsertVarURL == "" {
		missing = append(missing, envVarUpsertVarURL)
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing rewst env vars: %s", strings.Join(missing, ", "))
	}
	return nil
}
