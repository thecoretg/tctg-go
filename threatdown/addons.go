package threatdown

import (
	"context"
	"fmt"
)

type AddOn struct {
	Product                 string `json:"product"`
	TermType                string `json:"term_type"`
	TermLength              int    `json:"term_length"`
	TrialExtensionAvailable bool   `json:"trial_extension_available"`
	AddOnAllocation         int    `json:"add_on_allocation"`
}

func (c *Client) ListAddOns(ctx context.Context, siteID string) ([]AddOn, error) {
	type resp []AddOn
	e := fmt.Sprintf("sites/%s/addOns", siteID)
	result, err := get[resp](ctx, c, endpointURLV2(e), nil)
	if err != nil {
		return nil, fmt.Errorf("list add ons: %w", err)
	}

	return *result, nil
}
