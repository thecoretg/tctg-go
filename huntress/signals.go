package huntress

import (
	"context"
	"fmt"
)

// Signal represents a Huntress signal.
type Signal struct {
	ID                   int64          `json:"id,omitempty"`
	CreatedAt            string         `json:"created_at,omitempty"`
	Details              map[string]any `json:"details,omitempty"`
	Entity               map[string]any `json:"entity,omitempty"`
	InvestigatedAt       string         `json:"investigated_at,omitempty"`
	InvestigationContext string         `json:"investigation_context,omitempty"`
	Name                 string         `json:"name,omitempty"`
	Organization         map[string]any `json:"organization,omitempty"`
	Status               string         `json:"status,omitempty"`
	Type                 string         `json:"type,omitempty"`
	UpdatedAt            string         `json:"updated_at,omitempty"`
}

// SignalsResponse is a page of signals.
type SignalsResponse struct {
	Signals    []Signal   `json:"signals"`
	Pagination Pagination `json:"pagination"`
}

// ListSignals returns a single page of signals.
func (c *Client) ListSignals(ctx context.Context, params map[string]string) (*SignalsResponse, error) {
	result, err := get[SignalsResponse](ctx, c, endpointURL("signals"), params)
	if err != nil {
		return nil, fmt.Errorf("list signals: %w", err)
	}
	return result, nil
}

// GetSignal returns a single signal by ID.
func (c *Client) GetSignal(ctx context.Context, id int64) (*Signal, error) {
	result, err := get[Signal](ctx, c, endpointURL(fmt.Sprintf("signals/%d", id)), nil)
	if err != nil {
		return nil, fmt.Errorf("get signal: %w", err)
	}
	return result, nil
}

// ListAccountSignals returns a single page of signals for an account (Reseller
// credentials only).
func (c *Client) ListAccountSignals(ctx context.Context, accountID int64, params map[string]string) (*SignalsResponse, error) {
	result, err := get[SignalsResponse](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/signals", accountID)), params)
	if err != nil {
		return nil, fmt.Errorf("list account signals: %w", err)
	}
	return result, nil
}

// GetAccountSignal returns a single signal for an account (Reseller credentials only).
func (c *Client) GetAccountSignal(ctx context.Context, accountID, id int64) (*Signal, error) {
	result, err := get[Signal](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/signals/%d", accountID, id)), nil)
	if err != nil {
		return nil, fmt.Errorf("get account signal: %w", err)
	}
	return result, nil
}
