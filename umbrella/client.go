// Package umbrella is a client for the Cisco Umbrella Providers API, which
// manages the customers belonging to a provider. It authenticates with the
// OAuth2 client-credentials flow and is built on the shared internal/httpx
// transport for retries and JSON handling.
package umbrella

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/thecoretg/tctg-go/internal/httpx"
	"github.com/thecoretg/tctg-go/ratelimit"
)

const (
	// apiHost is the root every Umbrella API path hangs off.
	apiHost  = "https://api.umbrella.com"
	tokenURL = apiHost + "/auth/v2/token"
	baseURL  = apiHost + "/admin/v2"
	// deploymentsBaseURL is the base for the deployments APIs (policies,
	// roaming computers), which live under a different path than the
	// providers/admin endpoints.
	deploymentsBaseURL = apiHost + "/deployments/v2"
)

type Config struct {
	ClientID     string
	ClientSecret string
	// Limiter paces requests. nil uses DefaultLimiter(); ratelimit.None turns
	// pacing off. Share one limiter between clients on the same API key.
	Limiter ratelimit.Limiter
}

type Client struct {
	httpClient *http.Client
	// creds, ctx and policy are what httpClient was built from, kept so ForOrg can build
	// a child client with a token of its own and the same limiter.
	creds  *clientcredentials.Config
	ctx    context.Context //nolint:containedctx // the token sources' lifetime, as NewClient's ctx
	policy httpx.Policy

	mu   sync.Mutex
	orgs map[int]*Client
}

// orgHeader names the child org a token is requested for.
const orgHeader = "X-Umbrella-OrgId"

var acceptJSON = map[string]string{"Accept": "application/json"}

// NewClient builds a Client using the OAuth2 client-credentials flow. The
// credential is granted the provider scopes it is entitled to, so no scopes
// are requested per token.
func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	var missing []string
	if cfg.ClientID == "" {
		missing = append(missing, "ClientID")
	}
	if cfg.ClientSecret == "" {
		missing = append(missing, "ClientSecret")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing umbrella config fields: %s", strings.Join(missing, ", "))
	}

	lim := cfg.Limiter
	if lim == nil {
		lim = DefaultLimiter()
	}
	creds := &clientcredentials.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		TokenURL:     tokenURL,
	}
	return newClient(ctx, creds, httpx.Policy{Retries: 3, Limiter: lim}), nil
}

func newClient(ctx context.Context, creds *clientcredentials.Config, policy httpx.Policy) *Client {
	oauthClient := oauth2.NewClient(ctx, creds.TokenSource(ctx))
	return &Client{
		httpClient: httpx.NewClient(oauthClient.Transport, acceptJSON, policy),
		creds:      creds,
		ctx:        ctx,
		policy:     policy,
		orgs:       map[int]*Client{},
	}
}

// ForOrg returns a client that acts on one of the provider's child orgs, the customer
// with customerID. It holds a token issued for that org (the token request carries the
// X-Umbrella-OrgId header), cached and refreshed on its own, and shares this client's
// limiter, so it spends the same per-key budget. The same customerID returns the same
// client, so its token is reused.
//
// The header on the provider's own token is not enough: the live API ignores it and
// answers child-scoped calls for the provider's org.
//
// Use it for the child-scoped APIs, such as GetOrganizationInfo and the roaming computer
// calls. The provider endpoints under /admin/v2/providers/customers take the customer in
// the path and need the provider's client.
func (c *Client) ForOrg(customerID int) *Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	if child, ok := c.orgs[customerID]; ok {
		return child
	}
	tokenHTTP := &http.Client{Transport: orgTransport{base: http.DefaultTransport, org: strconv.Itoa(customerID)}}
	if hc, ok := c.ctx.Value(oauth2.HTTPClient).(*http.Client); ok && hc.Transport != nil {
		tokenHTTP.Transport = orgTransport{base: hc.Transport, org: strconv.Itoa(customerID)}
	}
	child := newClient(context.WithValue(c.ctx, oauth2.HTTPClient, tokenHTTP), c.creds, c.policy)
	c.orgs[customerID] = child
	return child
}

// orgTransport adds the child org's header to the token request.
type orgTransport struct {
	base http.RoundTripper
	org  string
}

func (t orgTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r = r.Clone(r.Context())
	r.Header.Set(orgHeader, t.org)
	return t.base.RoundTrip(r)
}

// DefaultLimiter paces a client to Umbrella's three per-key limits, which are
// enforced together: 5 per second, 14 per minute and 350 per 30 minutes. One
// token every 30m/348 with bursts of 2 satisfies all three at once: at most 2
// in a second, 13.6 in a minute and 350 in 30 minutes. That is a steady
// request every ~5.2s, the 30-minute window's sustained rate. Each call returns
// a new limiter.
func DefaultLimiter() *rate.Limiter {
	return rate.NewLimiter(rate.Every(30*time.Minute/348), 2)
}

func NewClientFromEnv(ctx context.Context) (*Client, error) {
	return NewClient(ctx, Config{
		ClientID:     os.Getenv("UMBRELLA_CLIENT_ID"),
		ClientSecret: os.Getenv("UMBRELLA_CLIENT_SECRET"),
	})
}

// hostURL builds a URL from a path relative to the API host, such as
// "admin/v2/providers/customers".
func hostURL(path string) string {
	return fmt.Sprintf("%s/%s", apiHost, path)
}

func endpointURL(path string) string {
	return fmt.Sprintf("%s/%s", baseURL, path)
}

func deploymentsURL(path string) string {
	return fmt.Sprintf("%s/%s", deploymentsBaseURL, path)
}
