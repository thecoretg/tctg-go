package threatdown

import (
	"context"
	"fmt"
	"os"
	"strings"

	"golang.org/x/time/rate"

	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/thecoretg/tctg-go/internal/httpx"
	"github.com/thecoretg/tctg-go/ratelimit"
)

const (
	tokenURL = "https://api.threatdown.com/oneview/oauth2/token"
	// oneViewRoot is the root of the versioned OneView APIs.
	oneViewRoot = "https://api.threatdown.com/oneview"
	baseURLV1   = oneViewRoot + "/v1"
	baseURLV2   = oneViewRoot + "/v2"
)

type Config struct {
	ClientID     string
	ClientSecret string
	// Limiter paces requests. nil uses DefaultLimiter(); ratelimit.None turns
	// pacing off. Share one limiter between clients on the same OAuth application.
	Limiter ratelimit.Limiter
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

	lim := cfg.Limiter
	if lim == nil {
		lim = DefaultLimiter()
	}

	oauthClient := oauth2.NewClient(ctx, ts)
	hc := httpx.NewClient(oauthClient.Transport, map[string]string{"Accept": "application/json"}, httpx.Policy{Retries: 3, Limiter: lim})

	return &Client{httpClient: hc}, nil
}

// DefaultLimiter paces a client to ThreatDown's 360 requests per minute, a
// leaky bucket counted per OAuth application: 6 per second sustained, with
// bursts of 6 since the bucket's depth is not published. ThreatDown raises the
// limit on request; pass a limiter with the new rate when it does. Each call
// returns a new limiter.
func DefaultLimiter() *rate.Limiter {
	return rate.NewLimiter(6, 6)
}

func NewClientFromEnv(ctx context.Context) (*Client, error) {
	return NewClient(ctx, Config{
		ClientID:     os.Getenv("THREATDOWN_CLIENT_ID"),
		ClientSecret: os.Getenv("THREATDOWN_CLIENT_SECRET"),
	})
}

// oneViewURL builds a URL from a path relative to the OneView root, such as
// "v1/users".
func oneViewURL(path string) string {
	return fmt.Sprintf("%s/%s", oneViewRoot, path)
}

func endpointURLV1(endpoint string) string {
	return fmt.Sprintf("%s/%s", baseURLV1, endpoint)
}

func endpointURLV2(endpoint string) string {
	return fmt.Sprintf("%s/%s", baseURLV2, endpoint)
}
