package huntress

import (
	"context"
	"fmt"
)

// Escalation represents a Huntress escalation.
type Escalation struct {
	ID            int64    `json:"id,omitempty"`
	Account       *Account `json:"account,omitempty"`
	Organizations []string `json:"organizations,omitempty"`
	CreatedAt     string   `json:"created_at,omitempty"`
	ResolvedAt    string   `json:"resolved_at,omitempty"`
	Severity      string   `json:"severity,omitempty"`
	Status        string   `json:"status,omitempty"`
	Subject       string   `json:"subject,omitempty"`
	Subtype       string   `json:"subtype,omitempty"`
	Type          string   `json:"type,omitempty"`
	UpdatedAt     string   `json:"updated_at,omitempty"`
}

// EscalationWithEntities is an escalation with its associated entities,
// returned when fetching a single escalation.
type EscalationWithEntities struct {
	Escalation
	Entities map[string]any `json:"entities,omitempty"`
}

// EscalationResolution is the result of resolving an escalation.
type EscalationResolution struct {
	Escalation       EscalationWithEntities `json:"escalation"`
	ResolutionMethod string                 `json:"resolution_method"`
}

// EscalationResolutionParameters is the request body for resolving an escalation.
type EscalationResolutionParameters struct {
	Determination string `json:"determination,omitempty"`
	Scope         string `json:"scope,omitempty"`
}

// EscalationsResponse is a page of escalations.
type EscalationsResponse struct {
	Escalations []Escalation `json:"escalations"`
	Pagination  Pagination   `json:"pagination"`
}

// ListEscalations returns a single page of escalations.
func (c *Client) ListEscalations(ctx context.Context, params map[string]string) (*EscalationsResponse, error) {
	result, err := c.Get[EscalationsResponse](ctx, endpointURL("escalations"), params)
	if err != nil {
		return nil, fmt.Errorf("list escalations: %w", err)
	}
	return result, nil
}

// GetEscalation returns a single escalation by ID, including its entities.
func (c *Client) GetEscalation(ctx context.Context, id int64) (*EscalationWithEntities, error) {
	result, err := c.Get[EscalationWithEntities](ctx, endpointURL(fmt.Sprintf("escalations/%d", id)), nil)
	if err != nil {
		return nil, fmt.Errorf("get escalation: %w", err)
	}
	return result, nil
}

// ResolveEscalation creates a resolution for an escalation.
func (c *Client) ResolveEscalation(ctx context.Context, id int64, body EscalationResolutionParameters) (*EscalationResolution, error) {
	result, err := c.Post[EscalationResolution](ctx, endpointURL(fmt.Sprintf("escalations/%d/resolution", id)), body)
	if err != nil {
		return nil, fmt.Errorf("resolve escalation: %w", err)
	}
	return result, nil
}
