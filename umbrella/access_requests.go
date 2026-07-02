package umbrella

import (
	"context"
	"fmt"
)

// AccessRequest is the information about an organization's access request. The
// timestamp fields are epoch milliseconds.
type AccessRequest struct {
	ID               int    `json:"id"`
	OrganizationID   int    `json:"organizationId"`
	State            string `json:"state"`
	OrganizationName string `json:"organizationName"`
	DisplayAt        int64  `json:"displayAt,omitempty"`
	CreatedAt        int64  `json:"createdAt,omitempty"`
	ModifiedAt       int64  `json:"modifiedAt,omitempty"`
}

// CreateAccessRequest creates an access request for the customer's organization.
func (c *Client) CreateAccessRequest(ctx context.Context, customerID int) (*AccessRequest, error) {
	url := endpointURL(fmt.Sprintf("providers/customers/%d/accessRequests", customerID))
	result, err := post[AccessRequest](ctx, c, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create access request: %w", err)
	}
	return result, nil
}

// GetAccessRequest gets the access request details for the customer's organization.
func (c *Client) GetAccessRequest(ctx context.Context, customerID, accessRequestID int) (*AccessRequest, error) {
	url := endpointURL(fmt.Sprintf("providers/customers/%d/accessRequests/%d", customerID, accessRequestID))
	result, err := get[AccessRequest](ctx, c, url, nil)
	if err != nil {
		return nil, fmt.Errorf("get access request: %w", err)
	}
	return result, nil
}

// UpdateAccessRequest updates the access request for the customer's organization.
func (c *Client) UpdateAccessRequest(ctx context.Context, customerID, accessRequestID int) (*AccessRequest, error) {
	url := endpointURL(fmt.Sprintf("providers/customers/%d/accessRequests/%d", customerID, accessRequestID))
	result, err := put[AccessRequest](ctx, c, url, nil)
	if err != nil {
		return nil, fmt.Errorf("update access request: %w", err)
	}
	return result, nil
}
