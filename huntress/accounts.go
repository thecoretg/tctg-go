package huntress

import (
	"context"
	"fmt"
)

// Account represents a Huntress account.
type Account struct {
	ID                int64          `json:"id,omitempty"`
	Name              string         `json:"name,omitempty"`
	Subdomain         string         `json:"subdomain,omitempty"`
	Status            string         `json:"status,omitempty"`
	SupportType       string         `json:"support_type,omitempty"`
	NeighborhoodWatch map[string]any `json:"neighborhood_watch,omitempty"`
	BillingAddress    *Address       `json:"billing_address,omitempty"`
	ShippingAddress   *Address       `json:"shipping_address,omitempty"`
}

// AccountCreationParameters is the request body for creating an account.
type AccountCreationParameters struct {
	Name                   string         `json:"name"`
	Subdomain              string         `json:"subdomain"`
	PhoneNumber            string         `json:"phone_number"`
	Admin                  map[string]any `json:"admin"`
	AdditionalAdminEmails  []string       `json:"additional_admin_emails,omitempty"`
	SupportType            string         `json:"support_type,omitempty"`
	Products               []string       `json:"products,omitempty"`
	ProductTrials          []string       `json:"product_trials,omitempty"`
	ProductTrialsStartDate string         `json:"product_trials_start_date,omitempty"`
	BillingAddress         *Address       `json:"billing_address,omitempty"`
	ShippingAddress        *Address       `json:"shipping_address,omitempty"`
}

// AccountUpdateParameters is the request body for updating an account.
type AccountUpdateParameters struct {
	Name            string   `json:"name,omitempty"`
	Subdomain       string   `json:"subdomain,omitempty"`
	PhoneNumber     string   `json:"phone_number,omitempty"`
	SupportType     string   `json:"support_type,omitempty"`
	BillingAddress  *Address `json:"billing_address,omitempty"`
	ShippingAddress *Address `json:"shipping_address,omitempty"`
}

// AccountsResponse is a page of accounts.
type AccountsResponse struct {
	Accounts   []Account  `json:"accounts"`
	Pagination Pagination `json:"pagination"`
}

// GetAccount returns the account associated with the API credentials.
func (c *Client) GetAccount(ctx context.Context) (*Account, error) {
	result, err := get[Account](ctx, c, endpointURL("account"), nil)
	if err != nil {
		return nil, fmt.Errorf("get account: %w", err)
	}
	return result, nil
}

// ListAccounts returns a single page of accounts (Reseller credentials only).
func (c *Client) ListAccounts(ctx context.Context, params map[string]string) (*AccountsResponse, error) {
	result, err := get[AccountsResponse](ctx, c, endpointURL("accounts"), params)
	if err != nil {
		return nil, fmt.Errorf("list accounts: %w", err)
	}
	return result, nil
}

// CreateAccount creates an account (Reseller credentials only).
func (c *Client) CreateAccount(ctx context.Context, body AccountCreationParameters) (*Account, error) {
	result, err := post[Account](ctx, c, endpointURL("accounts"), body)
	if err != nil {
		return nil, fmt.Errorf("create account: %w", err)
	}
	return result, nil
}

// GetAccountByID returns a single account by ID (Reseller credentials only).
func (c *Client) GetAccountByID(ctx context.Context, accountID int64) (*Account, error) {
	result, err := get[Account](ctx, c, endpointURL(fmt.Sprintf("accounts/%d", accountID)), nil)
	if err != nil {
		return nil, fmt.Errorf("get account: %w", err)
	}
	return result, nil
}

// UpdateAccount updates an account (Reseller credentials only).
func (c *Client) UpdateAccount(ctx context.Context, accountID int64, body AccountUpdateParameters) (*Account, error) {
	result, err := patch[Account](ctx, c, endpointURL(fmt.Sprintf("accounts/%d", accountID)), body)
	if err != nil {
		return nil, fmt.Errorf("update account: %w", err)
	}
	return result, nil
}

// DeleteAccount deletes an account (Reseller credentials only).
func (c *Client) DeleteAccount(ctx context.Context, accountID int64) error {
	if err := del(ctx, c, endpointURL(fmt.Sprintf("accounts/%d", accountID))); err != nil {
		return fmt.Errorf("delete account: %w", err)
	}
	return nil
}
