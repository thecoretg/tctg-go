package huntress

import (
	"context"
	"fmt"
)

// Invoice represents a Huntress invoice.
type Invoice struct {
	ID           int64  `json:"id,omitempty"`
	Amount       int    `json:"amount,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	CurrencyType string `json:"currency_type,omitempty"`
	Plan         string `json:"plan,omitempty"`
	Quantity     int64  `json:"quantity,omitempty"`
	Receipt      string `json:"receipt,omitempty"`
	Status       string `json:"status,omitempty"`
	HasUsage     bool   `json:"has_usage,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

// AccountUsageLineItem is a per-account usage line item on a reseller invoice.
type AccountUsageLineItem struct {
	ID           int64          `json:"id,omitempty"`
	PeriodStart  string         `json:"period_start,omitempty"`
	PeriodEnd    string         `json:"period_end,omitempty"`
	Account      map[string]any `json:"account,omitempty"`
	Product      string         `json:"product,omitempty"`
	Subscription map[string]any `json:"subscription,omitempty"`
	Usage        map[string]any `json:"usage,omitempty"`
}

// OrganizationUsageLineItem is a per-organization usage line item on a reseller invoice.
type OrganizationUsageLineItem struct {
	ID           int64          `json:"id,omitempty"`
	PeriodStart  string         `json:"period_start,omitempty"`
	PeriodEnd    string         `json:"period_end,omitempty"`
	Account      map[string]any `json:"account,omitempty"`
	Organization map[string]any `json:"organization,omitempty"`
	ActualUsage  map[string]any `json:"actual_usage,omitempty"`
}

// InvoicesResponse is a page of invoices.
type InvoicesResponse struct {
	Invoices   []Invoice  `json:"invoices"`
	Pagination Pagination `json:"pagination"`
}

// AccountUsageLineItemsResponse is a page of account usage line items.
type AccountUsageLineItemsResponse struct {
	AccountUsageLineItems []AccountUsageLineItem `json:"account_usage_line_items"`
	Pagination            Pagination             `json:"pagination"`
}

// OrganizationUsageLineItemsResponse is a page of organization usage line items.
type OrganizationUsageLineItemsResponse struct {
	OrganizationUsageLineItems []OrganizationUsageLineItem `json:"organization_usage_line_items"`
	Pagination                 Pagination                  `json:"pagination"`
}

type invoiceEnvelope struct {
	Invoice Invoice `json:"invoice"`
}

// ListInvoices returns a single page of invoices.
func (c *Client) ListInvoices(ctx context.Context, params map[string]string) (*InvoicesResponse, error) {
	result, err := get[InvoicesResponse](ctx, c, endpointURL("invoices"), params)
	if err != nil {
		return nil, fmt.Errorf("list invoices: %w", err)
	}
	return result, nil
}

// GetInvoice returns a single invoice by ID.
func (c *Client) GetInvoice(ctx context.Context, id int64) (*Invoice, error) {
	result, err := get[invoiceEnvelope](ctx, c, endpointURL(fmt.Sprintf("invoices/%d", id)), nil)
	if err != nil {
		return nil, fmt.Errorf("get invoice: %w", err)
	}
	return &result.Invoice, nil
}

// ListAccountInvoices returns a single page of invoices for an account (Reseller
// credentials only).
func (c *Client) ListAccountInvoices(ctx context.Context, accountID int64, params map[string]string) (*InvoicesResponse, error) {
	result, err := get[InvoicesResponse](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/invoices", accountID)), params)
	if err != nil {
		return nil, fmt.Errorf("list account invoices: %w", err)
	}
	return result, nil
}

// GetAccountInvoice returns a single invoice for an account (Reseller credentials only).
func (c *Client) GetAccountInvoice(ctx context.Context, accountID, id int64) (*Invoice, error) {
	result, err := get[invoiceEnvelope](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/invoices/%d", accountID, id)), nil)
	if err != nil {
		return nil, fmt.Errorf("get account invoice: %w", err)
	}
	return &result.Invoice, nil
}

// ListResellerInvoices returns a single page of reseller invoices (Reseller
// credentials only).
func (c *Client) ListResellerInvoices(ctx context.Context, params map[string]string) (*InvoicesResponse, error) {
	result, err := get[InvoicesResponse](ctx, c, endpointURL("reseller/invoices"), params)
	if err != nil {
		return nil, fmt.Errorf("list reseller invoices: %w", err)
	}
	return result, nil
}

// GetResellerInvoice returns a single reseller invoice by ID (Reseller
// credentials only).
func (c *Client) GetResellerInvoice(ctx context.Context, id int64) (*Invoice, error) {
	result, err := get[Invoice](ctx, c, endpointURL(fmt.Sprintf("reseller/invoices/%d", id)), nil)
	if err != nil {
		return nil, fmt.Errorf("get reseller invoice: %w", err)
	}
	return result, nil
}

// ListInvoiceAccountUsageLineItems returns a single page of account usage line
// items for a reseller invoice (Reseller credentials only).
func (c *Client) ListInvoiceAccountUsageLineItems(ctx context.Context, invoiceID int64, params map[string]string) (*AccountUsageLineItemsResponse, error) {
	result, err := get[AccountUsageLineItemsResponse](ctx, c, endpointURL(fmt.Sprintf("reseller/invoices/%d/account_usage_line_items", invoiceID)), params)
	if err != nil {
		return nil, fmt.Errorf("list invoice account usage line items: %w", err)
	}
	return result, nil
}

// ListInvoiceOrganizationUsageLineItems returns a single page of organization
// usage line items for a reseller invoice (Reseller credentials only).
func (c *Client) ListInvoiceOrganizationUsageLineItems(ctx context.Context, invoiceID int64, params map[string]string) (*OrganizationUsageLineItemsResponse, error) {
	result, err := get[OrganizationUsageLineItemsResponse](ctx, c, endpointURL(fmt.Sprintf("reseller/invoices/%d/organization_usage_line_items", invoiceID)), params)
	if err != nil {
		return nil, fmt.Errorf("list invoice organization usage line items: %w", err)
	}
	return result, nil
}
