package huntress

import (
	"context"
	"fmt"
)

// Membership represents a user's membership in an account or organization.
type Membership struct {
	ID           int64          `json:"id,omitempty"`
	Permissions  string         `json:"permissions,omitempty"`
	Account      map[string]any `json:"account,omitempty"`
	Organization map[string]any `json:"organization,omitempty"`
	User         map[string]any `json:"user,omitempty"`
	CreatedAt    string         `json:"created_at,omitempty"`
	UpdatedAt    string         `json:"updated_at,omitempty"`
}

// MemberInvitation is returned when a new membership is created for a user who
// must still accept an invitation.
type MemberInvitation struct {
	Permissions  string         `json:"permissions,omitempty"`
	Account      map[string]any `json:"account,omitempty"`
	Organization map[string]any `json:"organization,omitempty"`
	User         map[string]any `json:"user,omitempty"`
}

// MembershipCreationParameters is the request body for creating a membership.
type MembershipCreationParameters struct {
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Permissions    string `json:"permissions"`
	OrganizationID int    `json:"organization_id,omitempty"`
}

// MembershipUpdateParameters is the request body for updating a membership.
type MembershipUpdateParameters struct {
	Permissions string `json:"permissions,omitempty"`
}

// MembershipsResponse is a page of memberships.
type MembershipsResponse struct {
	Memberships []Membership `json:"memberships"`
	Pagination  Pagination   `json:"pagination"`
}

type membershipEnvelope struct {
	Membership Membership `json:"membership"`
}

// ListMemberships returns a single page of memberships.
func (c *Client) ListMemberships(ctx context.Context, params map[string]string) (*MembershipsResponse, error) {
	result, err := get[MembershipsResponse](ctx, c, endpointURL("memberships"), params)
	if err != nil {
		return nil, fmt.Errorf("list memberships: %w", err)
	}
	return result, nil
}

// GetMembership returns a single membership by ID.
func (c *Client) GetMembership(ctx context.Context, id int64) (*Membership, error) {
	result, err := get[membershipEnvelope](ctx, c, endpointURL(fmt.Sprintf("memberships/%d", id)), nil)
	if err != nil {
		return nil, fmt.Errorf("get membership: %w", err)
	}
	return &result.Membership, nil
}

// CreateMembership invites a user and creates a membership.
func (c *Client) CreateMembership(ctx context.Context, body MembershipCreationParameters) (*MemberInvitation, error) {
	result, err := post[MemberInvitation](ctx, c, endpointURL("memberships"), body)
	if err != nil {
		return nil, fmt.Errorf("create membership: %w", err)
	}
	return result, nil
}

// UpdateMembership updates a user's membership permissions.
func (c *Client) UpdateMembership(ctx context.Context, id int64, body MembershipUpdateParameters) (*Membership, error) {
	result, err := patch[Membership](ctx, c, endpointURL(fmt.Sprintf("memberships/%d", id)), body)
	if err != nil {
		return nil, fmt.Errorf("update membership: %w", err)
	}
	return result, nil
}

// DeleteMembership deletes a membership.
func (c *Client) DeleteMembership(ctx context.Context, id int64) error {
	if err := del(ctx, c, endpointURL(fmt.Sprintf("memberships/%d", id))); err != nil {
		return fmt.Errorf("delete membership: %w", err)
	}
	return nil
}

// ListAccountMemberships returns a single page of memberships for an account
// (Reseller credentials only).
func (c *Client) ListAccountMemberships(ctx context.Context, accountID int64, params map[string]string) (*MembershipsResponse, error) {
	result, err := get[MembershipsResponse](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/memberships", accountID)), params)
	if err != nil {
		return nil, fmt.Errorf("list account memberships: %w", err)
	}
	return result, nil
}

// GetAccountMembership returns a single membership for an account (Reseller
// credentials only).
func (c *Client) GetAccountMembership(ctx context.Context, accountID, id int64) (*Membership, error) {
	result, err := get[membershipEnvelope](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/memberships/%d", accountID, id)), nil)
	if err != nil {
		return nil, fmt.Errorf("get account membership: %w", err)
	}
	return &result.Membership, nil
}

// CreateAccountMembership invites a user to an account (Reseller credentials only).
func (c *Client) CreateAccountMembership(ctx context.Context, accountID int64, body MembershipCreationParameters) (*MemberInvitation, error) {
	result, err := post[MemberInvitation](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/memberships", accountID)), body)
	if err != nil {
		return nil, fmt.Errorf("create account membership: %w", err)
	}
	return result, nil
}

// UpdateAccountMembership updates a membership within an account (Reseller
// credentials only).
func (c *Client) UpdateAccountMembership(ctx context.Context, accountID, id int64, body MembershipUpdateParameters) (*Membership, error) {
	result, err := patch[Membership](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/memberships/%d", accountID, id)), body)
	if err != nil {
		return nil, fmt.Errorf("update account membership: %w", err)
	}
	return result, nil
}

// DeleteAccountMembership deletes a membership within an account (Reseller
// credentials only).
func (c *Client) DeleteAccountMembership(ctx context.Context, accountID, id int64) error {
	if err := del(ctx, c, endpointURL(fmt.Sprintf("accounts/%d/memberships/%d", accountID, id))); err != nil {
		return fmt.Errorf("delete account membership: %w", err)
	}
	return nil
}
