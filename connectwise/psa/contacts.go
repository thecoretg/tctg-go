package psa

import (
	"context"
	"fmt"
)

func contactIDEndpoint(contactID int) string {
	return fmt.Sprintf("company/contacts/%d", contactID)
}

func (c *Client) PostContact(ctx context.Context, contact *Contact) (*Contact, error) {
	return c.Post[Contact](ctx, "company/contacts", contact)
}

func (c *Client) ListContacts(ctx context.Context, params map[string]string, opts ...ListOption) ([]Contact, error) {
	return c.GetMany[Contact](ctx, "company/contacts", params, opts...)
}

func (c *Client) GetContact(ctx context.Context, contactID int, params map[string]string) (*Contact, error) {
	return c.Get[Contact](ctx, contactIDEndpoint(contactID), params)
}

func (c *Client) PutContact(ctx context.Context, contactID int, contact *Contact) (*Contact, error) {
	return c.Put[Contact](ctx, contactIDEndpoint(contactID), contact)
}

func (c *Client) PatchContact(ctx context.Context, contactID int, patchOps []PatchOp) (*Contact, error) {
	return c.Patch[Contact](ctx, contactIDEndpoint(contactID), patchOps)
}

func (c *Client) DeleteContact(ctx context.Context, contactID int) error {
	return c.Delete(ctx, contactIDEndpoint(contactID))
}
