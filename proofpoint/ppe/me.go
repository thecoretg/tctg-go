package ppe

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// Me returns the account the client is authenticated as. Its
// EntityPrimaryDomain is the organization the account belongs to.
func (c *Client) Me(ctx context.Context) (*Me, error) {
	return c.Get[Me](ctx, c.endpointURL("me"), nil)
}

// Stacks reports the stacks a domain is active on, such as ["us3"]. Any stack
// can answer this, so it works before you know which stack to configure.
func (c *Client) Stacks(ctx context.Context, domain string) ([]string, error) {
	res, err := c.Get[struct {
		Endpoints []string `json:"endpoints"`
	}](ctx, c.endpointURL("endpoints", domain), nil)
	if err != nil {
		return nil, err
	}

	stacks := make([]string, 0, len(res.Endpoints))
	for _, e := range res.Endpoints {
		u, err := url.Parse(e)
		if err != nil || u.Hostname() == "" {
			return nil, fmt.Errorf("unexpected endpoint %q", e)
		}
		stack, _, _ := strings.Cut(u.Hostname(), ".")
		stacks = append(stacks, stack)
	}
	return stacks, nil
}
