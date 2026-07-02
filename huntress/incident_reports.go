package huntress

import (
	"context"
	"fmt"
)

// IncidentReport represents a Huntress incident report.
type IncidentReport struct {
	ID              int64          `json:"id,omitempty"`
	AccountID       int64          `json:"account_id,omitempty"`
	AgentID         int64          `json:"agent_id,omitempty"`
	Body            string         `json:"body,omitempty"`
	ClosedAt        string         `json:"closed_at,omitempty"`
	IndicatorCounts map[string]any `json:"indicator_counts,omitempty"`
	IndicatorTypes  []string       `json:"indicator_types,omitempty"`
	OrganizationID  int64          `json:"organization_id,omitempty"`
	Platform        string         `json:"platform,omitempty"`
	Remediations    map[string]any `json:"remediations,omitempty"`
	SentAt          string         `json:"sent_at,omitempty"`
	Severity        string         `json:"severity,omitempty"`
	Status          string         `json:"status,omitempty"`
	StatusUpdatedAt string         `json:"status_updated_at,omitempty"`
	Subject         string         `json:"subject,omitempty"`
	Summary         string         `json:"summary,omitempty"`
	UpdatedAt       string         `json:"updated_at,omitempty"`
}

// Remediation represents a remediation action attached to an incident report.
type Remediation struct {
	ID          int64    `json:"id,omitempty"`
	Type        string   `json:"type,omitempty"`
	Action      string   `json:"action,omitempty"`
	Parameters  []string `json:"parameters,omitempty"`
	Status      string   `json:"status,omitempty"`
	ApprovedAt  string   `json:"approved_at,omitempty"`
	ApprovedBy  *User    `json:"approved_by,omitempty"`
	CompletedAt string   `json:"completed_at,omitempty"`
}

// RemediationBulkRejectionParameters is the request body for rejecting all
// remediations on an incident report.
type RemediationBulkRejectionParameters struct {
	Comment     string `json:"comment"`
	Useful      bool   `json:"useful"`
	Name        string `json:"name,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Email       string `json:"email,omitempty"`
}

// IncidentReportsResponse is a page of incident reports.
type IncidentReportsResponse struct {
	IncidentReports []IncidentReport `json:"incident_reports"`
	Pagination      Pagination       `json:"pagination"`
}

type incidentReportEnvelope struct {
	IncidentReport IncidentReport `json:"incident_report"`
}

// ListIncidentReports returns a single page of incident reports.
func (c *Client) ListIncidentReports(ctx context.Context, params map[string]string) (*IncidentReportsResponse, error) {
	result, err := get[IncidentReportsResponse](ctx, c, endpointURL("incident_reports"), params)
	if err != nil {
		return nil, fmt.Errorf("list incident reports: %w", err)
	}
	return result, nil
}

// GetIncidentReport returns a single incident report by ID.
func (c *Client) GetIncidentReport(ctx context.Context, id int64) (*IncidentReport, error) {
	result, err := get[incidentReportEnvelope](ctx, c, endpointURL(fmt.Sprintf("incident_reports/%d", id)), nil)
	if err != nil {
		return nil, fmt.Errorf("get incident report: %w", err)
	}
	return &result.IncidentReport, nil
}

// ResolveIncidentReport resolves an incident report.
func (c *Client) ResolveIncidentReport(ctx context.Context, id int64) (*IncidentReport, error) {
	result, err := post[IncidentReport](ctx, c, endpointURL(fmt.Sprintf("incident_reports/%d/resolution", id)), nil)
	if err != nil {
		return nil, fmt.Errorf("resolve incident report: %w", err)
	}
	return result, nil
}

// ListRemediations returns the remediations for an incident report.
func (c *Client) ListRemediations(ctx context.Context, incidentReportID int64) ([]Remediation, error) {
	result, err := get[[]Remediation](ctx, c, endpointURL(fmt.Sprintf("incident_reports/%d/remediations", incidentReportID)), nil)
	if err != nil {
		return nil, fmt.Errorf("list remediations: %w", err)
	}
	return *result, nil
}

// GetRemediation returns a single remediation for an incident report.
func (c *Client) GetRemediation(ctx context.Context, incidentReportID, remediationID int64) (*Remediation, error) {
	result, err := get[Remediation](ctx, c, endpointURL(fmt.Sprintf("incident_reports/%d/remediations/%d", incidentReportID, remediationID)), nil)
	if err != nil {
		return nil, fmt.Errorf("get remediation: %w", err)
	}
	return result, nil
}

// ApproveRemediations approves all remediations on an incident report.
func (c *Client) ApproveRemediations(ctx context.Context, incidentReportID int64) (*IncidentReport, error) {
	result, err := post[IncidentReport](ctx, c, endpointURL(fmt.Sprintf("incident_reports/%d/remediations/bulk_approval", incidentReportID)), nil)
	if err != nil {
		return nil, fmt.Errorf("approve remediations: %w", err)
	}
	return result, nil
}

// RejectRemediations rejects all remediations on an incident report.
func (c *Client) RejectRemediations(ctx context.Context, incidentReportID int64, body RemediationBulkRejectionParameters) error {
	if _, err := post[struct{}](ctx, c, endpointURL(fmt.Sprintf("incident_reports/%d/remediations/bulk_rejection", incidentReportID)), body); err != nil {
		return fmt.Errorf("reject remediations: %w", err)
	}
	return nil
}

// ListAccountIncidentReports returns a single page of incident reports for an
// account (Reseller credentials only).
func (c *Client) ListAccountIncidentReports(ctx context.Context, accountID int64, params map[string]string) (*IncidentReportsResponse, error) {
	result, err := get[IncidentReportsResponse](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/incident_reports", accountID)), params)
	if err != nil {
		return nil, fmt.Errorf("list account incident reports: %w", err)
	}
	return result, nil
}

// GetAccountIncidentReport returns a single incident report for an account
// (Reseller credentials only).
func (c *Client) GetAccountIncidentReport(ctx context.Context, accountID, id int64) (*IncidentReport, error) {
	result, err := get[incidentReportEnvelope](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/incident_reports/%d", accountID, id)), nil)
	if err != nil {
		return nil, fmt.Errorf("get account incident report: %w", err)
	}
	return &result.IncidentReport, nil
}

// ResolveAccountIncidentReport resolves an incident report for an account
// (Reseller credentials only).
func (c *Client) ResolveAccountIncidentReport(ctx context.Context, accountID, id int64) (*IncidentReport, error) {
	result, err := post[IncidentReport](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/incident_reports/%d/resolution", accountID, id)), nil)
	if err != nil {
		return nil, fmt.Errorf("resolve account incident report: %w", err)
	}
	return result, nil
}

// ListAccountRemediations returns the remediations for an account's incident
// report (Reseller credentials only).
func (c *Client) ListAccountRemediations(ctx context.Context, accountID, incidentReportID int64) ([]Remediation, error) {
	result, err := get[[]Remediation](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/incident_reports/%d/remediations", accountID, incidentReportID)), nil)
	if err != nil {
		return nil, fmt.Errorf("list account remediations: %w", err)
	}
	return *result, nil
}

// GetAccountRemediation returns a single remediation for an account's incident
// report (Reseller credentials only).
func (c *Client) GetAccountRemediation(ctx context.Context, accountID, incidentReportID, remediationID int64) (*Remediation, error) {
	result, err := get[Remediation](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/incident_reports/%d/remediations/%d", accountID, incidentReportID, remediationID)), nil)
	if err != nil {
		return nil, fmt.Errorf("get account remediation: %w", err)
	}
	return result, nil
}

// ApproveAccountRemediations approves all remediations on an account's incident
// report (Reseller credentials only).
func (c *Client) ApproveAccountRemediations(ctx context.Context, accountID, incidentReportID int64) (*IncidentReport, error) {
	result, err := post[IncidentReport](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/incident_reports/%d/remediations/bulk_approval", accountID, incidentReportID)), nil)
	if err != nil {
		return nil, fmt.Errorf("approve account remediations: %w", err)
	}
	return result, nil
}

// RejectAccountRemediations rejects all remediations on an account's incident
// report (Reseller credentials only).
func (c *Client) RejectAccountRemediations(ctx context.Context, accountID, incidentReportID int64, body RemediationBulkRejectionParameters) error {
	if _, err := post[struct{}](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/incident_reports/%d/remediations/bulk_rejection", accountID, incidentReportID)), body); err != nil {
		return fmt.Errorf("reject account remediations: %w", err)
	}
	return nil
}
