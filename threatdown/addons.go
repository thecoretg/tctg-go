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

// AddOnInput is the minimal payload accepted when setting a site's add-ons.
type AddOnInput struct {
	Product  string `json:"product"`
	TermType string `json:"term_type"`
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

// UpdateSiteAddOns replaces the full set of add-ons for a site. The caller is
// responsible for preserving any add-ons that should remain by including them in
// addOns (this endpoint is a full PUT, not a patch).
func (c *Client) UpdateSiteAddOns(ctx context.Context, siteID string, addOns []AddOnInput) ([]AddOn, error) {
	e := fmt.Sprintf("sites/%s/addOns", siteID)
	result, err := put[[]AddOn](ctx, c, endpointURLV2(e), addOns)
	if err != nil {
		return nil, fmt.Errorf("update add ons: %w", err)
	}
	return *result, nil
}
