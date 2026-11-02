package ppe

import "context"

// ListUsers returns every User in the Organization whose primary domain is
// domain.
func (c *Client) ListUsers(ctx context.Context, domain string) ([]User, error) {
	return c.getList[User](ctx, c.endpointURL("orgs", domain, "users"), "users")
}

// GetUser returns the User whose primary email is email.
func (c *Client) GetUser(ctx context.Context, domain, email string) (*User, error) {
	return c.Get[User](ctx, c.endpointURL("orgs", domain, "users", email), nil)
}

// CreateUser creates user in the Organization whose primary domain is domain.
func (c *Client) CreateUser(ctx context.Context, domain string, user User) (*User, error) {
	return c.Post[User](ctx, c.endpointURL("orgs", domain, "users"), user)
}

// UpdateUser sends user's set fields to the User whose primary email is email.
// The API does not document whether PUT replaces or merges, so only fields
// that are set are sent.
func (c *Client) UpdateUser(ctx context.Context, domain, email string, user User) error {
	return c.Put(ctx, c.endpointURL("orgs", domain, "users", email), user)
}

// DeleteUser deletes the User whose primary email is email. This cannot be
// undone.
func (c *Client) DeleteUser(ctx context.Context, domain, email string) error {
	return c.Delete(ctx, c.endpointURL("orgs", domain, "users", email))
}
