package entra

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"
)

// guidRE matches a directory or object ID. Entra substitutes the tenant GUID
// into the discovery document's issuer however the tenant is addressed, so a
// domain name (or common/organizations) fails the OIDC issuer check. Requiring
// a GUID here is also what makes the tid check meaningful.
var guidRE = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// Config wires an Auth. TenantID, ClientID, ClientSecret and BaseURL are
// required; everything else has a default noted on the field.
type Config struct {
	// TenantID is the directory (tenant) GUID. Domain names are rejected.
	TenantID string
	// ClientID is the application (client) ID of the app registration.
	ClientID string
	// ClientSecret is a client secret from Certificates & secrets. They expire;
	// when one does, every sign-in fails with ErrInvalidClient.
	ClientSecret string
	// BaseURL is the public origin browsers use, e.g. https://app.example.com.
	// https is required except for localhost. No path, query or fragment. The
	// redirect URI is BaseURL + CallbackPath and is never derived from the
	// incoming request, because a reverse proxy makes the request's scheme and
	// host unreliable and a mismatch surfaces as AADSTS50011 at the token
	// endpoint, after the user has already authenticated.
	BaseURL string

	// CallbackPath is where Microsoft redirects after sign-in. Default
	// "/auth/sso/callback". Mount Callback here.
	CallbackPath string
	// LoginPath is where unauthenticated browsers and failed sign-ins are sent.
	// Default "/login".
	LoginPath string
	// LogoutRedirect is where Logout sends the browser. Default LoginPath.
	LogoutRedirect string
	// SuccessPath is where a completed sign-in lands when no valid ?next was
	// given. Default "/".
	SuccessPath string

	// SessionTTL bounds a session. It slides forward on activity. Default 8h.
	// Group removal in Entra only blocks the next sign-in, so this is also the
	// revocation window for a running session.
	SessionTTL time.Duration
	// FlowTTL bounds the round trip to Microsoft. Default 10m.
	FlowTTL time.Duration
	// SessionCookie and FlowCookie are cookie names. Defaults "entra_session"
	// and "entra_flow".
	SessionCookie string
	FlowCookie    string

	// Scopes requested from Microsoft. Default openid, profile, email. profile
	// supplies oid, name and preferred_username; email is needed because
	// managed users have no email claim by default. offline_access is always
	// removed: the package mints its own session and never calls Microsoft
	// again, so a refresh token would be a liability with no use.
	Scopes []string

	// Authorize runs after tenant verification and before Provision. nil
	// allows every user in the tenant. See RequireGroup and RequireRole.
	Authorize Authorizer

	// Sessions and States persist sessions and in-flight sign-ins. Both
	// default to one shared MemoryStore, which is single-instance only.
	Sessions SessionStore
	States   StateStore

	// Logger receives security-relevant detail (refused tenant IDs, group
	// lists) that is deliberately kept out of browser-facing messages.
	// Default slog.Default().
	Logger *slog.Logger
	// HTTPClient is used for discovery, JWKS and the token endpoint. Default
	// http.DefaultClient.
	HTTPClient *http.Client

	// OnError handles a failed sign-in in Start or Callback. err wraps one of
	// the package sentinels where one applies. The default logs the cause and
	// redirects to LoginPath?err=<Message(err)>.
	OnError func(w http.ResponseWriter, r *http.Request, err error)
	// Unauthorized handles a request that RequireAuth rejects. The default
	// answers API-shaped requests with a JSON 401 and redirects browsers to
	// LoginPath?next=<request path>.
	Unauthorized func(w http.ResponseWriter, r *http.Request)

	// issuer overrides the discovery URL. Tests point it at a fake provider.
	issuer string
}

// Provisioner is the application's side of just-in-time provisioning. It owns
// the user model and its storage; the package never sees either beyond the
// opaque userID it stores in the session.
//
// Link local records on Claims.OID only. Do not resurrect a record the
// application has disabled: return ErrUserDisabled from Provision instead, or
// a naive upsert-on-sign-in silently re-enables whoever was just off-boarded.
type Provisioner[U any] interface {
	// Provision is called once per successful sign-in, after the tenant check
	// and Authorizer have passed. It creates or updates the local record and
	// returns a stable identifier for the session. The record itself is not
	// needed here: RequireAuth loads it through Lookup on every request.
	Provision(ctx context.Context, c Claims) (userID string, err error)
	// Lookup reloads the user on every authenticated request, so disabling an
	// account takes effect on the next request rather than at session expiry.
	// Return ErrUserDisabled or ErrUserNotFound to end the session.
	Lookup(ctx context.Context, userID string) (U, error)
}

// New validates cfg and returns an Auth. It performs no network I/O; OIDC
// discovery runs lazily on the first sign-in so a Microsoft outage cannot stop
// the application from starting. ctx is accepted for signature parity with the
// other tctg-go constructors.
func New[U any](_ context.Context, cfg Config, p Provisioner[U]) (*Auth[U], error) {
	if p == nil {
		return nil, fmt.Errorf("entra: Provisioner is required")
	}

	var missing []string
	if cfg.TenantID == "" {
		missing = append(missing, "TenantID")
	}
	if cfg.ClientID == "" {
		missing = append(missing, "ClientID")
	}
	if cfg.ClientSecret == "" {
		missing = append(missing, "ClientSecret")
	}
	if cfg.BaseURL == "" {
		missing = append(missing, "BaseURL")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing entra config fields: %s", strings.Join(missing, ", "))
	}

	if !guidRE.MatchString(cfg.TenantID) {
		return nil, fmt.Errorf("entra: TenantID must be the directory GUID, not a domain name or %q", cfg.TenantID)
	}
	base, err := validateBaseURL(cfg.BaseURL)
	if err != nil {
		return nil, err
	}
	cfg.BaseURL = base

	cfg.CallbackPath = cmp.Or(cfg.CallbackPath, "/auth/sso/callback")
	cfg.LoginPath = cmp.Or(cfg.LoginPath, "/login")
	cfg.LogoutRedirect = cmp.Or(cfg.LogoutRedirect, cfg.LoginPath)
	cfg.SuccessPath = cmp.Or(cfg.SuccessPath, "/")
	for _, pth := range []string{cfg.CallbackPath, cfg.LoginPath, cfg.LogoutRedirect, cfg.SuccessPath} {
		if !strings.HasPrefix(pth, "/") {
			return nil, fmt.Errorf("entra: path %q must start with /", pth)
		}
	}

	cfg.SessionTTL = cmp.Or(cfg.SessionTTL, 8*time.Hour)
	cfg.FlowTTL = cmp.Or(cfg.FlowTTL, 10*time.Minute)
	cfg.SessionCookie = cmp.Or(cfg.SessionCookie, "entra_session")
	cfg.FlowCookie = cmp.Or(cfg.FlowCookie, "entra_flow")

	if len(cfg.Scopes) == 0 {
		cfg.Scopes = []string{"openid", "profile", "email"}
	}
	cfg.Scopes = slices.DeleteFunc(slices.Clone(cfg.Scopes), func(s string) bool { return s == "offline_access" })
	if !slices.Contains(cfg.Scopes, "openid") {
		cfg.Scopes = append([]string{"openid"}, cfg.Scopes...)
	}

	if cfg.Sessions == nil || cfg.States == nil {
		mem := NewMemoryStore()
		if cfg.Sessions == nil {
			cfg.Sessions = mem
		}
		if cfg.States == nil {
			cfg.States = mem
		}
	}
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
	}
	cfg.issuer = cmp.Or(cfg.issuer, "https://login.microsoftonline.com/"+cfg.TenantID+"/v2.0")

	a := &Auth[U]{cfg: cfg, prov: p}
	if a.cfg.OnError == nil {
		a.cfg.OnError = a.defaultOnError
	}
	if a.cfg.Unauthorized == nil {
		a.cfg.Unauthorized = a.defaultUnauthorized
	}
	return a, nil
}

// NewFromEnv builds an Auth from ENTRA_TENANT_ID, ENTRA_CLIENT_ID,
// ENTRA_CLIENT_SECRET and ENTRA_BASE_URL. When ENTRA_GROUP_ID is set, sign-in
// is restricted to members of that group via RequireGroup.
func NewFromEnv[U any](ctx context.Context, p Provisioner[U]) (*Auth[U], error) {
	cfg := Config{
		TenantID:     os.Getenv("ENTRA_TENANT_ID"),
		ClientID:     os.Getenv("ENTRA_CLIENT_ID"),
		ClientSecret: os.Getenv("ENTRA_CLIENT_SECRET"),
		BaseURL:      os.Getenv("ENTRA_BASE_URL"),
	}
	if g := strings.TrimSpace(os.Getenv("ENTRA_GROUP_ID")); g != "" {
		if !guidRE.MatchString(g) {
			return nil, fmt.Errorf("entra: ENTRA_GROUP_ID must be the group's object ID (a GUID), not its name")
		}
		cfg.Authorize = RequireGroup(g)
	}
	return New(ctx, cfg, p)
}

// validateBaseURL accepts an https origin, or http for localhost, with no
// path, query or fragment. The returned value has no trailing slash.
func validateBaseURL(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("entra: BaseURL: %w", err)
	}
	if u.Host == "" {
		return "", fmt.Errorf("entra: BaseURL %q must be an absolute URL", raw)
	}
	if u.RawQuery != "" || u.Fragment != "" || (u.Path != "" && u.Path != "/") {
		return "", fmt.Errorf("entra: BaseURL %q must be an origin only, with no path, query or fragment", raw)
	}
	switch u.Scheme {
	case "https":
	case "http":
		if h := u.Hostname(); h != "localhost" && h != "127.0.0.1" && h != "::1" {
			return "", fmt.Errorf("entra: BaseURL %q must use https (http is allowed for localhost only)", raw)
		}
	default:
		return "", fmt.Errorf("entra: BaseURL %q must use https", raw)
	}
	return u.Scheme + "://" + u.Host, nil
}
