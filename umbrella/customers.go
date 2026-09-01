package umbrella

import (
	"context"
	"fmt"
)

// Customer is the properties of a provider's customer.
type Customer struct {
	CustomerID         int      `json:"customerId"`
	CustomerName       string   `json:"customerName"`
	LicenseType        string   `json:"licenseType,omitempty"`
	PackageName        string   `json:"packageName"`
	PackageID          int      `json:"packageId"`
	Seats              int      `json:"seats"`
	StreetAddress      string   `json:"streetAddress"`
	StreetAddress2     string   `json:"streetAddress2,omitempty"`
	City               string   `json:"city"`
	State              string   `json:"state,omitempty"`
	CountryCode        string   `json:"countryCode"`
	ZipCode            string   `json:"zipCode,omitempty"`
	DealID             string   `json:"dealId,omitempty"`
	AdminEmails        []string `json:"adminEmails"`
	CcwDealOwnerEmails []string `json:"ccwDealOwnerEmails,omitempty"`
	CreatedAt          string   `json:"createdAt,omitempty"`
	ModifiedAt         string   `json:"modifiedAt,omitempty"`
	IsTrial            bool     `json:"isTrial,omitempty"`
	AddonRbi           string   `json:"addonRbi,omitempty"`
	AddonDlp           bool     `json:"addonDlp,omitempty"`
	AddonCdfwL7        bool     `json:"addonCdfwL7,omitempty"`
}

// CustomerCreateRequest is the body for creating a customer. LicenseType,
// IsTrial, and the addon fields apply only to specific package types; see the
// spec for the applicable combinations.
type CustomerCreateRequest struct {
	CustomerName       string   `json:"customerName"`
	Seats              int      `json:"seats"`
	StreetAddress      string   `json:"streetAddress"`
	City               string   `json:"city"`
	CountryCode        string   `json:"countryCode"`
	PackageID          int      `json:"packageId"`
	AdminEmails        []string `json:"adminEmails"`
	StreetAddress2     string   `json:"streetAddress2,omitempty"`
	State              string   `json:"state,omitempty"`
	ZipCode            string   `json:"zipCode,omitempty"`
	DealID             string   `json:"dealId,omitempty"`
	CcwDealOwnerEmails []string `json:"ccwDealOwnerEmails,omitempty"`
	LicenseType        string   `json:"licenseType,omitempty"`
	IsTrial            *bool    `json:"isTrial,omitempty"`
	AddonRbi           string   `json:"addonRbi,omitempty"`
	AddonDlp           *bool    `json:"addonDlp,omitempty"`
	AddonCdfwL7        *bool    `json:"addonCdfwL7,omitempty"`
}

// CustomerUpdateRequest is the body for updating a customer.
type CustomerUpdateRequest struct {
	CustomerName       string   `json:"customerName"`
	Seats              int      `json:"seats"`
	StreetAddress      string   `json:"streetAddress"`
	City               string   `json:"city"`
	CountryCode        string   `json:"countryCode"`
	PackageID          int      `json:"packageId"`
	AdminEmails        []string `json:"adminEmails"`
	StreetAddress2     string   `json:"streetAddress2,omitempty"`
	State              string   `json:"state,omitempty"`
	ZipCode            string   `json:"zipCode,omitempty"`
	DealID             string   `json:"dealId,omitempty"`
	CcwDealOwnerEmails []string `json:"ccwDealOwnerEmails,omitempty"`
}

// CreateCustomer creates a customer for the provider.
func (c *Client) CreateCustomer(ctx context.Context, body CustomerCreateRequest) (*Customer, error) {
	result, err := c.Post[Customer](ctx, endpointURL("providers/customers"), body)
	if err != nil {
		return nil, fmt.Errorf("create customer: %w", err)
	}
	return result, nil
}

// ListCustomers lists a single page of customers for the provider. Supported
// params are "page" and "limit" (max 100).
func (c *Client) ListCustomers(ctx context.Context, params map[string]string) ([]Customer, error) {
	result, err := c.Get[[]Customer](ctx, endpointURL("providers/customers"), params)
	if err != nil {
		return nil, fmt.Errorf("list customers: %w", err)
	}
	return *result, nil
}

// GetCustomer gets a customer for the provider by ID.
func (c *Client) GetCustomer(ctx context.Context, customerID int) (*Customer, error) {
	result, err := c.Get[Customer](ctx, endpointURL(fmt.Sprintf("providers/customers/%d", customerID)), nil)
	if err != nil {
		return nil, fmt.Errorf("get customer: %w", err)
	}
	return result, nil
}

// UpdateCustomer updates a customer for the provider by ID.
func (c *Client) UpdateCustomer(ctx context.Context, customerID int, body CustomerUpdateRequest) (*Customer, error) {
	result, err := c.Put[Customer](ctx, endpointURL(fmt.Sprintf("providers/customers/%d", customerID)), body)
	if err != nil {
		return nil, fmt.Errorf("update customer: %w", err)
	}
	return result, nil
}

// DeleteCustomer deletes a customer for the provider by ID.
func (c *Client) DeleteCustomer(ctx context.Context, customerID int) error {
	if err := c.Delete(ctx, endpointURL(fmt.Sprintf("providers/customers/%d", customerID))); err != nil {
		return fmt.Errorf("delete customer: %w", err)
	}
	return nil
}
