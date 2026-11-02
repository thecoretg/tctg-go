package ppe

import "context"

// GetBilling returns the billing data of domain's Organization.
func (c *Client) GetBilling(ctx context.Context, domain string) ([]BillingReport, error) {
	return c.getList[BillingReport](ctx, c.endpointURL("billing", domain))
}

// GetBillingWithChildren returns the billing data of domain's Organization
// and every Organization beneath it.
func (c *Client) GetBillingWithChildren(ctx context.Context, domain string) ([]BillingReport, error) {
	return c.getList[BillingReport](ctx, c.endpointURL("billing", domain, "orgs"))
}
