package umbrella

import (
	"context"
	"fmt"
	"strconv"
)

// CustomerDeal is the information about a customer deal.
type CustomerDeal struct {
	DealID         string          `json:"dealId"`
	EndCustomer    string          `json:"endCustomer"`
	MajorLineItems []MajorLineItem `json:"majorLineItems"`
	CanStampDeal   bool            `json:"canStampDeal"`
	TrialIDs       []string        `json:"trialIds"`
}

// MajorLineItem is an essential component of a deal.
type MajorLineItem struct {
	ObjectID       int    `json:"objectId,omitempty"`
	SourceAppRefID string `json:"sourceAppRefId,omitempty"`
}

// DealUpdateRequest is the body for updating a customer deal.
type DealUpdateRequest struct {
	CcoID          int             `json:"ccoid"`
	CustomerID     int             `json:"customerId"`
	QuoteID        int             `json:"quoteId,omitempty"`
	MajorLineItems []MajorLineItem `json:"majorLineItems,omitempty"`
}

// GetCustomerDeals gets the deals available to a customer by deal ID. ccoID is
// the ID of the user querying the deal and is required.
func (c *Client) GetCustomerDeals(ctx context.Context, dealID string, ccoID int) ([]CustomerDeal, error) {
	url := endpointURL(fmt.Sprintf("providers/customerDeals/%s", dealID))
	params := map[string]string{"ccoId": strconv.Itoa(ccoID)}
	result, err := get[[]CustomerDeal](ctx, c, url, params)
	if err != nil {
		return nil, fmt.Errorf("get customer deals: %w", err)
	}
	return *result, nil
}

// UpdateCustomerDeals updates a customer deal by deal ID.
func (c *Client) UpdateCustomerDeals(ctx context.Context, dealID string, body DealUpdateRequest) (*CustomerDeal, error) {
	url := endpointURL(fmt.Sprintf("providers/customerDeals/%s", dealID))
	result, err := put[CustomerDeal](ctx, c, url, body)
	if err != nil {
		return nil, fmt.Errorf("update customer deals: %w", err)
	}
	return result, nil
}
