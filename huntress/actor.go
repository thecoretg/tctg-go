package huntress

import (
	"context"
	"fmt"
)

// Actor describes the reseller, account, and user associated with the API
// credentials making the request.
type Actor struct {
	Reseller map[string]any `json:"reseller,omitempty"`
	Account  map[string]any `json:"account,omitempty"`
	User     map[string]any `json:"user,omitempty"`
}

// GetActor returns the actor (reseller/account/user) for the API credentials.
func (c *Client) GetActor(ctx context.Context) (*Actor, error) {
	result, err := c.Get[Actor](ctx, endpointURL("actor"), nil)
	if err != nil {
		return nil, fmt.Errorf("get actor: %w", err)
	}
	return result, nil
}
