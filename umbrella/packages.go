package umbrella

import (
	"context"
	"fmt"
)

// Package is the information about an Umbrella package.
type Package struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	PkgSeatMin  int    `json:"pkgSeatMin"`
	PpovSeatMin int    `json:"ppovSeatMin"`
}

// ListCustomerPackages lists the packages available to the trial customer.
func (c *Client) ListCustomerPackages(ctx context.Context) ([]Package, error) {
	result, err := c.Get[[]Package](ctx, endpointURL("providers/customers/packages"), nil)
	if err != nil {
		return nil, fmt.Errorf("list customer packages: %w", err)
	}
	return *result, nil
}
