package huntress

import (
	"context"
	"fmt"
)

// Organization represents a Huntress organization.
type Organization struct {
	ID                       int64    `json:"id,omitempty"`
	AgentsCount              int64    `json:"agents_count,omitempty"`
	AccountID                int64    `json:"account_id,omitempty"`
	CreatedAt                string   `json:"created_at,omitempty"`
	IncidentReportsCount     int64    `json:"incident_reports_count,omitempty"`
	Key                      string   `json:"key,omitempty"`
	LogsSourcesCount         int64    `json:"logs_sources_count,omitempty"`
	IdentityProviderTenantID string   `json:"identity_provider_tenant_id,omitempty"`
	BillableIdentityCount    int64    `json:"billable_identity_count,omitempty"`
	Name                     string   `json:"name,omitempty"`
	ReportRecipients         []string `json:"report_recipients,omitempty"`
	SatLearnerCount          int64    `json:"sat_learner_count,omitempty"`
	UpdatedAt                string   `json:"updated_at,omitempty"`
}

// OrganizationWithActualProductUsages is an organization with product-usage
// detail, returned when fetching a single organization.
type OrganizationWithActualProductUsages struct {
	Organization
	ActualUsages string `json:"actual_usages,omitempty"`
}

// OrganizationCreationParameters is the request body for creating an organization.
type OrganizationCreationParameters struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// OrganizationUpdateParameters is the request body for updating an organization.
type OrganizationUpdateParameters struct {
	Name             string   `json:"name,omitempty"`
	Key              string   `json:"key,omitempty"`
	ReportRecipients []string `json:"report_recipients,omitempty"`
}

// OrganizationsResponse is a page of organizations.
type OrganizationsResponse struct {
	Organizations []Organization `json:"organizations"`
	Pagination    Pagination     `json:"pagination"`
}

type organizationEnvelope struct {
	Organization OrganizationWithActualProductUsages `json:"organization"`
}

// ListOrganizations returns a single page of organizations.
func (c *Client) ListOrganizations(ctx context.Context, params map[string]string) (*OrganizationsResponse, error) {
	result, err := get[OrganizationsResponse](ctx, c, endpointURL("organizations"), params)
	if err != nil {
		return nil, fmt.Errorf("list organizations: %w", err)
	}
	return result, nil
}

// GetOrganization returns a single organization by ID.
func (c *Client) GetOrganization(ctx context.Context, id int64) (*OrganizationWithActualProductUsages, error) {
	result, err := get[organizationEnvelope](ctx, c, endpointURL(fmt.Sprintf("organizations/%d", id)), nil)
	if err != nil {
		return nil, fmt.Errorf("get organization: %w", err)
	}
	return &result.Organization, nil
}

// CreateOrganization creates an organization.
func (c *Client) CreateOrganization(ctx context.Context, body OrganizationCreationParameters) (*Organization, error) {
	result, err := post[Organization](ctx, c, endpointURL("organizations"), body)
	if err != nil {
		return nil, fmt.Errorf("create organization: %w", err)
	}
	return result, nil
}

// UpdateOrganization updates an organization.
func (c *Client) UpdateOrganization(ctx context.Context, id int64, body OrganizationUpdateParameters) (*Organization, error) {
	result, err := patch[Organization](ctx, c, endpointURL(fmt.Sprintf("organizations/%d", id)), body)
	if err != nil {
		return nil, fmt.Errorf("update organization: %w", err)
	}
	return result, nil
}

// DeleteOrganization deletes an organization.
func (c *Client) DeleteOrganization(ctx context.Context, id int64) error {
	if err := del(ctx, c, endpointURL(fmt.Sprintf("organizations/%d", id))); err != nil {
		return fmt.Errorf("delete organization: %w", err)
	}
	return nil
}

// ListAccountOrganizations returns a single page of organizations for an account
// (Reseller credentials only).
func (c *Client) ListAccountOrganizations(ctx context.Context, accountID int64, params map[string]string) (*OrganizationsResponse, error) {
	result, err := get[OrganizationsResponse](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/organizations", accountID)), params)
	if err != nil {
		return nil, fmt.Errorf("list account organizations: %w", err)
	}
	return result, nil
}

// GetAccountOrganization returns a single organization for an account (Reseller
// credentials only).
func (c *Client) GetAccountOrganization(ctx context.Context, accountID, id int64) (*OrganizationWithActualProductUsages, error) {
	result, err := get[organizationEnvelope](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/organizations/%d", accountID, id)), nil)
	if err != nil {
		return nil, fmt.Errorf("get account organization: %w", err)
	}
	return &result.Organization, nil
}

// CreateAccountOrganization creates an organization within an account (Reseller
// credentials only).
func (c *Client) CreateAccountOrganization(ctx context.Context, accountID int64, body OrganizationCreationParameters) (*Organization, error) {
	result, err := post[Organization](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/organizations", accountID)), body)
	if err != nil {
		return nil, fmt.Errorf("create account organization: %w", err)
	}
	return result, nil
}

// UpdateAccountOrganization updates an organization within an account (Reseller
// credentials only).
func (c *Client) UpdateAccountOrganization(ctx context.Context, accountID, id int64, body OrganizationUpdateParameters) (*Organization, error) {
	result, err := patch[Organization](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/organizations/%d", accountID, id)), body)
	if err != nil {
		return nil, fmt.Errorf("update account organization: %w", err)
	}
	return result, nil
}

// DeleteAccountOrganization deletes an organization within an account (Reseller
// credentials only).
func (c *Client) DeleteAccountOrganization(ctx context.Context, accountID, id int64) error {
	if err := del(ctx, c, endpointURL(fmt.Sprintf("accounts/%d/organizations/%d", accountID, id))); err != nil {
		return fmt.Errorf("delete account organization: %w", err)
	}
	return nil
}
