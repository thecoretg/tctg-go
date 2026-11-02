package ppe

import "context"

// GetLicensing returns the licensing data of domain's Organization.
func (c *Client) GetLicensing(ctx context.Context, domain string) (*Licensing, error) {
	return c.Get[Licensing](ctx, c.endpointURL("orgs", domain, "licensing"), nil)
}

// UpdateLicensing sends l's set fields as domain's Organization's licensing
// data.
func (c *Client) UpdateLicensing(ctx context.Context, domain string, l Licensing) error {
	return c.Put(ctx, c.endpointURL("orgs", domain, "licensing"), l)
}

// UpdatePackage changes the package of domain's Organization.
func (c *Client) UpdatePackage(ctx context.Context, domain string, p Package) error {
	return c.Put(ctx, c.endpointURL("orgs", domain, "package"), p)
}

// GetFeatures returns the features of domain's Organization. The set depends
// on its package.
func (c *Client) GetFeatures(ctx context.Context, domain string) (Features, error) {
	f, err := c.Get[Features](ctx, c.endpointURL("orgs", domain, "features"), nil)
	if err != nil {
		return nil, err
	}
	return *f, nil
}

// UpdateFeatures sets the features in f on domain's Organization. Setting a
// feature its package does not include returns ErrForbidden.
func (c *Client) UpdateFeatures(ctx context.Context, domain string, f Features) error {
	return c.Put(ctx, c.endpointURL("orgs", domain, "features"), f)
}

// ListProducts returns the products of domain's Organization, both available
// and purchased.
func (c *Client) ListProducts(ctx context.Context, domain string) ([]Product, error) {
	return c.getList[Product](ctx, c.endpointURL("orgs", domain, "products"), "products")
}

// GetProduct returns the product with label.
func (c *Client) GetProduct(ctx context.Context, domain, label string) (*Product, error) {
	return c.Get[Product](ctx, c.endpointURL("orgs", domain, "products", label), nil)
}

// PurchaseProduct purchases a product for domain's Organization.
func (c *Client) PurchaseProduct(ctx context.Context, domain string, p ProductPurchase) (*Product, error) {
	return c.Post[Product](ctx, c.endpointURL("orgs", domain, "products"), p)
}

// UpdateProduct applies update to the product with label.
func (c *Client) UpdateProduct(ctx context.Context, domain, label string, update ProductUpdate) error {
	return c.Patch(ctx, c.endpointURL("orgs", domain, "products", label), update)
}

// DeleteProduct cancels the product subscription with label. This cannot be
// undone. The API may refuse it depending on the account's permissions.
func (c *Client) DeleteProduct(ctx context.Context, domain, label string) error {
	return c.Delete(ctx, c.endpointURL("orgs", domain, "products", label))
}
