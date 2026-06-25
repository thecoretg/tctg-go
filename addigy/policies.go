package addigy

import (
	"context"
	"encoding/json"
	"fmt"
)

// Policy is an Addigy policy (definitions.Policy). The nested settings objects
// (splashtop, ssh, system updates, etc.) are kept as raw JSON in this scaffold;
// define typed structs for them as needed.
type Policy struct {
	AgentPath            string   `json:"agent_path,omitempty"`
	AgentVersion         string   `json:"agent_version,omitempty"`
	AutotaskAccountID    int      `json:"autotask_account_id,omitempty"`
	Color                string   `json:"color,omitempty"`
	ConnectwiseAccountID int      `json:"connectwise_account_id,omitempty"`
	CreationTime         float64  `json:"creation_time,omitempty"`
	DownloadPath         string   `json:"download_path,omitempty"`
	Icon                 string   `json:"icon,omitempty"`
	IgnoreUpdates        bool     `json:"ignore_updates,omitempty"`
	Instructions         []string `json:"instructions,omitempty"`
	ITGlueAccountID      string   `json:"itglue_account_id,omitempty"`
	LastDeployed         string   `json:"last_deployed,omitempty"`
	Name                 string   `json:"name,omitempty"`
	OrgID                string   `json:"orgid,omitempty"`
	Parent               string   `json:"parent,omitempty"`
	PolicyID             string   `json:"policyId,omitempty"`

	AddigyShieldSettings      json.RawMessage `json:"addigy_shield_settings,omitempty"`
	AddigySync                json.RawMessage `json:"addigy_sync,omitempty"`
	CollectorSettings         json.RawMessage `json:"collector_settings,omitempty"`
	PrebuiltAppSettings       json.RawMessage `json:"prebuilt_app_settings,omitempty"`
	Schedules                 json.RawMessage `json:"schedules,omitempty"`
	SelfServiceInstructionIDs map[string]bool `json:"self_service_instruction_ids,omitempty"`
	SplashtopSettings         json.RawMessage `json:"splashtop_settings,omitempty"`
	SSHSettings               json.RawMessage `json:"ssh_settings,omitempty"`
	SystemUpdatesSettings     json.RawMessage `json:"system_updates_settings,omitempty"`
	VNCSettings               json.RawMessage `json:"vnc_settings,omitempty"`
}

// RemoteControlSettings configures splashtop/ssh remote control on policy
// create/update (definitions.RemoteControlSettings).
type RemoteControlSettings struct {
	Enabled bool `json:"enabled,omitempty"`
}

// CreatePolicyRequest is the body for creating a policy (definitions.CreatePolicyRequest).
type CreatePolicyRequest struct {
	Name              string                 `json:"name"`
	Color             string                 `json:"color,omitempty"`
	Icon              string                 `json:"icon,omitempty"`
	ParentPolicyID    string                 `json:"parent_policy_id,omitempty"`
	SplashtopSettings *RemoteControlSettings `json:"splashtop_settings,omitempty"`
	SSHSettings       *RemoteControlSettings `json:"ssh_settings,omitempty"`
}

// PolicyUpdateRequest is the body for updating a policy (definitions.PolicyUpdateRequest).
type PolicyUpdateRequest struct {
	PolicyID string `json:"policy_id"`
	Name     string `json:"name,omitempty"`
	Color    string `json:"color,omitempty"`
	Icon     string `json:"icon,omitempty"`
}

// PolicyParentUpdateRequest sets a policy's parent (definitions.PolicyParentUpdateRequest).
type PolicyParentUpdateRequest struct {
	PolicyID       string `json:"policy_id"`
	ParentPolicyID string `json:"parent_policy_id"`
}

// PolicyQueryRequest filters a policy query by policy ID (definitions.PolicyQueryRequest).
type PolicyQueryRequest struct {
	Policies []string `json:"policies,omitempty"`
}

// CreatePolicy creates a policy (POST /o/{organization_id}/policies).
func (c *Client) CreatePolicy(ctx context.Context, req CreatePolicyRequest) (*Policy, error) {
	url, err := c.orgURL("policies")
	if err != nil {
		return nil, err
	}
	result, err := post[Policy](ctx, c, url, req)
	if err != nil {
		return nil, fmt.Errorf("create policy: %w", err)
	}
	return result, nil
}

// UpdatePolicy updates a policy (PUT /o/{organization_id}/policies).
func (c *Client) UpdatePolicy(ctx context.Context, req PolicyUpdateRequest) (*Policy, error) {
	url, err := c.orgURL("policies")
	if err != nil {
		return nil, err
	}
	result, err := put[Policy](ctx, c, url, req)
	if err != nil {
		return nil, fmt.Errorf("update policy: %w", err)
	}
	return result, nil
}

// DeletePolicy deletes a policy by ID (DELETE /o/{organization_id}/policies).
func (c *Client) DeletePolicy(ctx context.Context, policyID string) error {
	url, err := c.orgURL("policies")
	if err != nil {
		return err
	}
	if err := del(ctx, c, url, map[string]string{"id": policyID}); err != nil {
		return fmt.Errorf("delete policy: %w", err)
	}
	return nil
}

// UpdatePolicyParent sets a policy's parent (PUT /o/{organization_id}/policies/parent).
func (c *Client) UpdatePolicyParent(ctx context.Context, req PolicyParentUpdateRequest) error {
	url, err := c.orgURL("policies/parent")
	if err != nil {
		return err
	}
	if _, err := put[struct{}](ctx, c, url, req); err != nil {
		return fmt.Errorf("update policy parent: %w", err)
	}
	return nil
}

// DeletePolicyParent removes a policy's parent (DELETE /o/{organization_id}/policies/parent).
func (c *Client) DeletePolicyParent(ctx context.Context, policyID string) error {
	url, err := c.orgURL("policies/parent")
	if err != nil {
		return err
	}
	if err := del(ctx, c, url, map[string]string{"policy_id": policyID}); err != nil {
		return fmt.Errorf("delete policy parent: %w", err)
	}
	return nil
}

// QueryPolicies returns policy info, optionally filtered by policy ID
// (POST /oa/policies/query).
func (c *Client) QueryPolicies(ctx context.Context, req PolicyQueryRequest) ([]Policy, error) {
	result, err := post[[]Policy](ctx, c, c.url("oa/policies/query"), req)
	if err != nil {
		return nil, fmt.Errorf("query policies: %w", err)
	}
	return *result, nil
}
