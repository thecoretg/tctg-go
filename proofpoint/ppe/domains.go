package ppe

import "context"

// ListDomains returns every Domain of the Organization whose primary domain
// is domain.
func (c *Client) ListDomains(ctx context.Context, domain string) ([]Domain, error) {
	return c.getList[Domain](ctx, c.endpointURL("orgs", domain, "domains"), "domains")
}

// GetDomain returns the Domain named name in domain's Organization.
func (c *Client) GetDomain(ctx context.Context, domain, name string) (*Domain, error) {
	return c.Get[Domain](ctx, c.endpointURL("orgs", domain, "domains", name), nil)
}

// CreateDomain adds d to domain's Organization.
func (c *Client) CreateDomain(ctx context.Context, domain string, d Domain) (*Domain, error) {
	return c.Post[Domain](ctx, c.endpointURL("orgs", domain, "domains"), d)
}

// UpdateDomain sends d's set fields to the Domain named name.
func (c *Client) UpdateDomain(ctx context.Context, domain, name string, d Domain) error {
	return c.Put(ctx, c.endpointURL("orgs", domain, "domains", name), d)
}

// DeleteDomain removes the Domain named name from domain's Organization. This
// cannot be undone.
func (c *Client) DeleteDomain(ctx context.Context, domain, name string) error {
	return c.Delete(ctx, c.endpointURL("orgs", domain, "domains", name))
}
