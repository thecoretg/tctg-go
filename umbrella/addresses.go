package umbrella

import (
	"context"
	"fmt"
)

// CustomerAddress is a customer's address information.
type CustomerAddress struct {
	AccountID        int    `json:"accountId"`
	AccountSiteID    int    `json:"accountSiteId"`
	MappedCrPartyID  int    `json:"mappedCrPartyId"`
	OrganizationName string `json:"organizationName,omitempty"`
	StreetAddress    string `json:"streetAddress,omitempty"`
	StreetAddress2   string `json:"streetAddress2,omitempty"`
	City             string `json:"city,omitempty"`
	State            string `json:"state,omitempty"`
	CountryCode      string `json:"countryCode,omitempty"`
	ZipCode          string `json:"zipCode,omitempty"`
}

// ListCustomerAddresses lists a single page of customer addresses for the
// provider. Supported params are "page" and "limit" (max 100).
func (c *Client) ListCustomerAddresses(ctx context.Context, params map[string]string) ([]CustomerAddress, error) {
	result, err := get[[]CustomerAddress](ctx, c, endpointURL("providers/customerAddresses"), params)
	if err != nil {
		return nil, fmt.Errorf("list customer addresses: %w", err)
	}
	return *result, nil
}
