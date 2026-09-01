package umbrella

import (
	"context"
	"fmt"
)

// Policy is a stub view of an Umbrella policy as returned by ListPolicies.
type Policy struct {
	PolicyID       int    `json:"policyId"`
	OrganizationID int    `json:"organizationId"`
	Name           string `json:"name"`
	Priority       int    `json:"priority"`
	CreatedAt      string `json:"createdAt"`
	IsDefault      bool   `json:"isDefault"`
}

// ListPolicies lists a single page of Umbrella policies. Supported params are
// "page", "limit" (max 100), and "type" ("dns" or "web"; defaults to dns when
// omitted).
func (c *Client) ListPolicies(ctx context.Context, params map[string]string) ([]Policy, error) {
	result, err := c.Get[[]Policy](ctx, deploymentsURL("policies"), params)
	if err != nil {
		return nil, fmt.Errorf("list policies: %w", err)
	}
	return *result, nil
}

// AddPolicyIdentity adds an identity (by origin ID) to a policy and returns the
// origin ID that was added. Policy changes may take up to 20 minutes to take
// effect globally.
func (c *Client) AddPolicyIdentity(ctx context.Context, policyID, originID int) (int, error) {
	url := deploymentsURL(fmt.Sprintf("policies/%d/identities/%d", policyID, originID))
	result, err := c.Put[int](ctx, url, nil)
	if err != nil {
		return 0, fmt.Errorf("add policy identity: %w", err)
	}
	return *result, nil
}

// DeletePolicyIdentity removes an identity (by origin ID) from a policy. Policy
// changes may take up to 20 minutes to take effect globally.
func (c *Client) DeletePolicyIdentity(ctx context.Context, policyID, originID int) error {
	url := deploymentsURL(fmt.Sprintf("policies/%d/identities/%d", policyID, originID))
	if err := c.Delete(ctx, url); err != nil {
		return fmt.Errorf("delete policy identity: %w", err)
	}
	return nil
}
