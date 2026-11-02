package ppe

import "context"

// GetOrg returns the Organization whose primary domain is domain.
func (c *Client) GetOrg(ctx context.Context, domain string) (*Organization, error) {
	return c.Get[Organization](ctx, c.endpointURL("orgs", domain), nil)
}

// ListChildOrgs returns the Organizations directly beneath domain's
// Organization, such as a partner's customers.
func (c *Client) ListChildOrgs(ctx context.Context, domain string) ([]Organization, error) {
	return c.getList[Organization](ctx, c.endpointURL("orgs", domain, "orgs"), "orgs")
}

// CreateOrg creates org beneath the parent Organization whose primary domain
// is parent.
func (c *Client) CreateOrg(ctx context.Context, parent string, org Organization) (*Organization, error) {
	return c.Post[Organization](ctx, c.endpointURL("orgs", parent, "orgs"), org)
}

// UpdateOrg applies update to the Organization whose primary domain is domain.
func (c *Client) UpdateOrg(ctx context.Context, domain string, update OrgUpdate) error {
	return c.Patch(ctx, c.endpointURL("orgs", domain), update)
}

// DeleteOrg deletes the Organization whose primary domain is domain. This
// cannot be undone.
func (c *Client) DeleteOrg(ctx context.Context, domain string) error {
	return c.Delete(ctx, c.endpointURL("orgs", domain))
}
