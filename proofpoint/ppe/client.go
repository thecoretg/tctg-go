// Package ppe is a client for the Proofpoint Essentials Interface API (v1),
// which manages partner and customer organizations along with their users,
// domains, licensing and products. It covers Proofpoint Essentials only, not
// the Enterprise or TAP APIs.
//
// Every request authenticates with the X-User and X-Password headers of a
// Proofpoint Essentials admin account. What that account may do depends on its
// role: a partner admin can manage the organizations beneath it, an
// organization admin only its own users.
//
// Organizations live on a regional stack (us1–us5, eu1). Config.Stack picks
// the stack the client talks to; Stacks reports which stacks a domain is
// active on.
//
// Proofpoint publishes no rate limit for this API, so the client does not pace
// requests unless Config.Limiter is set. A 429 is retried, honoring
// Retry-After, and surfaces as a *ratelimit.Error if retries run out.
package ppe

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/thecoretg/tctg-go/internal/httpx"
	"github.com/thecoretg/tctg-go/ratelimit"
)

// DefaultStack is the stack used when Config.Stack is empty.
const DefaultStack = "us3"

type Config struct {
	// Username and Password are a Proofpoint Essentials admin login, sent as
	// the X-User and X-Password headers on every request.
	Username string
	Password string
	// Stack is the regional instance to call, such as "us3" or "eu1". Empty
	// uses DefaultStack.
	Stack string
	// Limiter paces requests. nil sends requests unpaced, since Proofpoint
	// publishes no limit for this API.
	Limiter ratelimit.Limiter
}

type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient builds a Client for cfg.Stack. ctx is accepted for signature
// parity with the other tctg-go clients; it is not used here.
func NewClient(_ context.Context, cfg Config) (*Client, error) {
	var missing []string
	if cfg.Username == "" {
		missing = append(missing, "Username")
	}
	if cfg.Password == "" {
		missing = append(missing, "Password")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing ppe config fields: %s", strings.Join(missing, ", "))
	}

	stack := cfg.Stack
	if stack == "" {
		stack = DefaultStack
	}

	headers := map[string]string{
		"X-User":     cfg.Username,
		"X-Password": cfg.Password,
		"Accept":     "application/json",
	}

	return &Client{
		httpClient: httpx.NewClient(nil, headers, httpx.Policy{Retries: 3, Limiter: cfg.Limiter}),
		baseURL:    StackURL(stack) + "/api/v1",
	}, nil
}

// NewClientFromEnv builds a Client from PROOFPOINT_USERNAME,
// PROOFPOINT_PASSWORD and PROOFPOINT_STACK (optional, defaults to
// DefaultStack).
func NewClientFromEnv(ctx context.Context) (*Client, error) {
	return NewClient(ctx, Config{
		Username: os.Getenv("PROOFPOINT_USERNAME"),
		Password: os.Getenv("PROOFPOINT_PASSWORD"),
		Stack:    os.Getenv("PROOFPOINT_STACK"),
	})
}

// StackURL returns the root URL of a stack, such as
// "https://us3.proofpointessentials.com" for "us3".
func StackURL(stack string) string {
	return fmt.Sprintf("https://%s.proofpointessentials.com", stack)
}

// endpointURL joins path segments onto the API base URL, escaping each one so
// domains and email addresses are safe to use as segments.
func (c *Client) endpointURL(segments ...string) string {
	escaped := make([]string, len(segments))
	for i, s := range segments {
		escaped[i] = url.PathEscape(s)
	}
	return c.baseURL + "/" + strings.Join(escaped, "/")
}
