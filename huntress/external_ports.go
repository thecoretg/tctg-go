package huntress

import (
	"context"
	"fmt"
)

// ExternalPort represents an externally reachable port found by external recon.
type ExternalPort struct {
	ID                 int64   `json:"id,omitempty"`
	IPAddress          string  `json:"ip_address,omitempty"`
	Port               int     `json:"port,omitempty"`
	Protocol           string  `json:"protocol,omitempty"`
	Service            string  `json:"service,omitempty"`
	RiskyService       bool    `json:"risky_service,omitempty"`
	LastScanAt         string  `json:"last_scan_at,omitempty"`
	LastExternalScanAt string  `json:"last_external_scan_at,omitempty"`
	OrganizationIDs    []int64 `json:"organization_ids,omitempty"`
}

// ExternalPortsResponse is a page of external ports.
type ExternalPortsResponse struct {
	ExternalPorts []ExternalPort `json:"external_ports"`
	Pagination    Pagination     `json:"pagination"`
}

type externalPortEnvelope struct {
	ExternalPort ExternalPort `json:"external_port"`
}

// ListExternalPorts returns a single page of external ports.
func (c *Client) ListExternalPorts(ctx context.Context, params map[string]string) (*ExternalPortsResponse, error) {
	result, err := get[ExternalPortsResponse](ctx, c, endpointURL("external_ports"), params)
	if err != nil {
		return nil, fmt.Errorf("list external ports: %w", err)
	}
	return result, nil
}

// GetExternalPort returns a single external port by ID.
func (c *Client) GetExternalPort(ctx context.Context, id int64) (*ExternalPort, error) {
	result, err := get[externalPortEnvelope](ctx, c, endpointURL(fmt.Sprintf("external_ports/%d", id)), nil)
	if err != nil {
		return nil, fmt.Errorf("get external port: %w", err)
	}
	return &result.ExternalPort, nil
}

// ListAccountExternalPorts returns a single page of external ports for an account
// (Reseller credentials only).
func (c *Client) ListAccountExternalPorts(ctx context.Context, accountID int64, params map[string]string) (*ExternalPortsResponse, error) {
	result, err := get[ExternalPortsResponse](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/external_ports", accountID)), params)
	if err != nil {
		return nil, fmt.Errorf("list account external ports: %w", err)
	}
	return result, nil
}

// GetAccountExternalPort returns a single external port for an account (Reseller
// credentials only).
func (c *Client) GetAccountExternalPort(ctx context.Context, accountID, id int64) (*ExternalPort, error) {
	result, err := get[externalPortEnvelope](ctx, c, endpointURL(fmt.Sprintf("accounts/%d/external_ports/%d", accountID, id)), nil)
	if err != nil {
		return nil, fmt.Errorf("get account external port: %w", err)
	}
	return &result.ExternalPort, nil
}
