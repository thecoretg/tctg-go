package huntress

import (
	"context"
	"fmt"
)

// Identity represents a Huntress ITDR identity.
type Identity struct {
	ID                    int64          `json:"id,omitempty"`
	Account               *Account       `json:"account,omitempty"`
	Organization          *Organization  `json:"organization,omitempty"`
	Tenant                map[string]any `json:"tenant,omitempty"`
	Username              string         `json:"username,omitempty"`
	Email                 string         `json:"email,omitempty"`
	Enabled               bool           `json:"enabled,omitempty"`
	Billable              bool           `json:"billable,omitempty"`
	EnabledProducts       []string       `json:"enabled_products,omitempty"`
	MFAEnabled            bool           `json:"mfa_enabled,omitempty"`
	RiskLevel             string         `json:"risk_level,omitempty"`
	RiskState             string         `json:"risk_state,omitempty"`
	RiskDetail            string         `json:"risk_detail,omitempty"`
	RiskLastUpdatedAt     string         `json:"risk_last_updated_at,omitempty"`
	PasswordLastChangedAt string         `json:"password_last_changed_at,omitempty"`
	OnPremisesSyncEnabled bool           `json:"on_premises_sync_enabled,omitempty"`
	External              bool           `json:"external,omitempty"`
	CreatedAt             string         `json:"created_at,omitempty"`
	UpdatedAt             string         `json:"updated_at,omitempty"`
}

// IdentitiesResponse is a page of identities.
type IdentitiesResponse struct {
	Identities []Identity `json:"identities"`
	Pagination Pagination `json:"pagination"`
}

// ListIdentities returns a single page of identities.
func (c *Client) ListIdentities(ctx context.Context, params map[string]string) (*IdentitiesResponse, error) {
	result, err := c.Get[IdentitiesResponse](ctx, endpointURL("identities"), params)
	if err != nil {
		return nil, fmt.Errorf("list identities: %w", err)
	}
	return result, nil
}

// GetIdentity returns a single identity by ID.
func (c *Client) GetIdentity(ctx context.Context, id int64) (*Identity, error) {
	result, err := c.Get[Identity](ctx, endpointURL(fmt.Sprintf("identities/%d", id)), nil)
	if err != nil {
		return nil, fmt.Errorf("get identity: %w", err)
	}
	return result, nil
}
