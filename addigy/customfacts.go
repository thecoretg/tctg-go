package addigy

import (
	"context"
	"fmt"
)

// Fact is a custom fact. The Addigy API returns slightly different shapes across
// endpoints (auditor_service.Fact and fact_entities.Fact); this struct is the
// superset of both, with optional fields.
type Fact struct {
	ID               string           `json:"id,omitempty"`
	Identifier       string           `json:"identifier,omitempty"`
	Name             string           `json:"name,omitempty"`
	Notes            string           `json:"notes,omitempty"`
	ReturnType       string           `json:"return_type,omitempty"`
	Version          int              `json:"version,omitempty"`
	CommunityFactID  string           `json:"community_fact_id,omitempty"`
	CommunityVersion int              `json:"community_version,omitempty"`
	InstructionID    string           `json:"instruction_id,omitempty"`
	OrgID            string           `json:"orgid,omitempty"`
	Provider         string           `json:"provider,omitempty"`
	Source           string           `json:"source,omitempty"`
	OSArchitectures  *OSArchitectures `json:"os_architectures,omitempty"`
}

// OSArchitectures holds the per-OS scripts for a custom fact
// (fact_entities.OSArchitectures).
type OSArchitectures struct {
	Linux *FactOSArchitecture `json:"linux,omitempty"`
	MacOS *FactOSArchitecture `json:"macOS,omitempty"`
}

// FactOSArchitecture is the script and metadata for one OS
// (fact_entities.FactOSArchitecture).
type FactOSArchitecture struct {
	IsSupported bool   `json:"is_supported,omitempty"`
	Language    string `json:"language,omitempty"`
	MD5Hash     string `json:"md5_hash,omitempty"`
	Script      string `json:"script,omitempty"`
	Shebang     string `json:"shebang,omitempty"`
}

// Instruction is the instruction backing a custom fact (fact_entities.Instruction).
type Instruction struct {
	Condition        string `json:"condition,omitempty"`
	Identifier       string `json:"identifier,omitempty"`
	InstructionID    string `json:"instruction_id,omitempty"`
	Label            string `json:"label,omitempty"`
	Name             string `json:"name,omitempty"`
	PolicyRestricted bool   `json:"policy_restricted,omitempty"`
	Provider         string `json:"provider,omitempty"`
	Public           bool   `json:"public,omitempty"`
	RemoveScript     string `json:"remove_script,omitempty"`
	RunOnSuccess     bool   `json:"run_on_success,omitempty"`
	StatusOnSkipped  string `json:"status_on_skipped,omitempty"`
	UserEmail        string `json:"user_email,omitempty"`
}

// FactsResponse is the paginated list of custom facts (auditor_service.FactsResponse).
type FactsResponse struct {
	Items    []Fact   `json:"items"`
	Metadata Metadata `json:"metadata"`
}

// FactResponse is the result of creating a custom fact (fact_entities.FactResponse).
type FactResponse struct {
	Fact        Fact        `json:"fact"`
	Instruction Instruction `json:"instruction"`
}

// FactPostRequest is the body for creating a custom fact (fact_entities.FactPostRequest).
type FactPostRequest struct {
	Name            string           `json:"name,omitempty"`
	Notes           string           `json:"notes,omitempty"`
	ReturnType      string           `json:"return_type,omitempty"`
	OSArchitectures *OSArchitectures `json:"os_architectures,omitempty"`
}

// FactPutRequest is the body for updating a custom fact (fact_entities.FactPutRequest).
type FactPutRequest struct {
	ID              string           `json:"id,omitempty"`
	Name            string           `json:"name,omitempty"`
	Notes           string           `json:"notes,omitempty"`
	ReturnType      string           `json:"return_type,omitempty"`
	OSArchitectures *OSArchitectures `json:"os_architectures,omitempty"`
}

// FactQuery filters a custom fact query (fact_entities.RequestQuery).
type FactQuery struct {
	Page          int         `json:"page,omitempty"`
	PerPage       int         `json:"per_page,omitempty"`
	Query         *FactFilter `json:"query,omitempty"`
	SortDirection string      `json:"sort_direction,omitempty"`
	SortField     string      `json:"sort_field,omitempty"`
}

// FactFilter narrows a custom fact query (fact_entities.Filter).
type FactFilter struct {
	ID           string `json:"id,omitempty"`
	NameContains string `json:"name_contains,omitempty"`
}

// AssignFactToPolicy assigns a custom fact to policies (facts.AssignToPolicy).
type AssignFactToPolicy struct {
	ID       string   `json:"id,omitempty"`
	Policies []string `json:"policies,omitempty"`
}

// FactPolicyResult reports which assignments succeeded/failed (fact_entities.Response).
type FactPolicyResult struct {
	Failed    []string `json:"failed,omitempty"`
	Succeeded []string `json:"succeeded,omitempty"`
}

// FactUsage lists where a custom fact is referenced (facts.FactUsage).
type FactUsage struct {
	Alerts       []UsageItem `json:"alerts,omitempty"`
	Integrations []UsageItem `json:"integrations,omitempty"`
	Policies     []UsageItem `json:"policies,omitempty"`
	Reports      []UsageItem `json:"reports,omitempty"`
	UserConfigs  []UsageItem `json:"user_configs,omitempty"`
}

// UsageItem is a single reference to a custom fact (facts.UsageItem).
type UsageItem struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	ParentID string `json:"parent_id,omitempty"`
}

// ListCustomFacts returns the organization's custom facts
// (GET /o/{organization_id}/facts/custom).
func (c *Client) ListCustomFacts(ctx context.Context) (*FactsResponse, error) {
	url, err := c.orgURL("facts/custom")
	if err != nil {
		return nil, err
	}
	result, err := c.Get[FactsResponse](ctx, url, nil)
	if err != nil {
		return nil, fmt.Errorf("list custom facts: %w", err)
	}
	return result, nil
}

// GetCustomFact returns a single custom fact by ID
// (GET /o/{organization_id}/facts/custom/{id}).
func (c *Client) GetCustomFact(ctx context.Context, id string) (*Fact, error) {
	url, err := c.orgURL("facts/custom/" + id)
	if err != nil {
		return nil, err
	}
	result, err := c.Get[Fact](ctx, url, nil)
	if err != nil {
		return nil, fmt.Errorf("get custom fact: %w", err)
	}
	return result, nil
}

// GetCustomFactUsage returns where a custom fact is referenced
// (GET /o/{organization_id}/facts/custom/{id}/usage).
func (c *Client) GetCustomFactUsage(ctx context.Context, id string) (*FactUsage, error) {
	url, err := c.orgURL("facts/custom/" + id + "/usage")
	if err != nil {
		return nil, err
	}
	result, err := c.Get[FactUsage](ctx, url, nil)
	if err != nil {
		return nil, fmt.Errorf("get custom fact usage: %w", err)
	}
	return result, nil
}

// CreateCustomFact creates a custom fact (POST /o/{organization_id}/facts/custom).
func (c *Client) CreateCustomFact(ctx context.Context, req FactPostRequest) (*FactResponse, error) {
	url, err := c.orgURL("facts/custom")
	if err != nil {
		return nil, err
	}
	result, err := c.Post[FactResponse](ctx, url, req)
	if err != nil {
		return nil, fmt.Errorf("create custom fact: %w", err)
	}
	return result, nil
}

// UpdateCustomFact updates a custom fact (PUT /o/{organization_id}/facts/custom).
func (c *Client) UpdateCustomFact(ctx context.Context, req FactPutRequest) error {
	url, err := c.orgURL("facts/custom")
	if err != nil {
		return err
	}
	if _, err := c.Put[struct{}](ctx, url, req); err != nil {
		return fmt.Errorf("update custom fact: %w", err)
	}
	return nil
}

// DeleteCustomFact deletes a custom fact by ID
// (DELETE /o/{organization_id}/facts/custom).
func (c *Client) DeleteCustomFact(ctx context.Context, id string) error {
	url, err := c.orgURL("facts/custom")
	if err != nil {
		return err
	}
	if err := c.Delete(ctx, url, map[string]string{"id": id}); err != nil {
		return fmt.Errorf("delete custom fact: %w", err)
	}
	return nil
}

// QueryCustomFacts returns custom facts filtered by id or name
// (POST /o/{organization_id}/facts/custom/query).
func (c *Client) QueryCustomFacts(ctx context.Context, req FactQuery) (*FactsResponse, error) {
	url, err := c.orgURL("facts/custom/query")
	if err != nil {
		return nil, err
	}
	result, err := c.Post[FactsResponse](ctx, url, req)
	if err != nil {
		return nil, fmt.Errorf("query custom facts: %w", err)
	}
	return result, nil
}

// AssignCustomFactToPolicy assigns a custom fact to one or more policies
// (POST /o/{organization_id}/facts/custom/policy).
func (c *Client) AssignCustomFactToPolicy(ctx context.Context, req AssignFactToPolicy) (*FactPolicyResult, error) {
	url, err := c.orgURL("facts/custom/policy")
	if err != nil {
		return nil, err
	}
	result, err := c.Post[FactPolicyResult](ctx, url, req)
	if err != nil {
		return nil, fmt.Errorf("assign custom fact to policy: %w", err)
	}
	return result, nil
}

// UnassignCustomFactFromPolicy removes a custom fact from a policy
// (DELETE /o/{organization_id}/facts/custom/policy).
func (c *Client) UnassignCustomFactFromPolicy(ctx context.Context, id, policyID string) error {
	url, err := c.orgURL("facts/custom/policy")
	if err != nil {
		return err
	}
	if err := c.Delete(ctx, url, map[string]string{"id": id, "policy_id": policyID}); err != nil {
		return fmt.Errorf("unassign custom fact from policy: %w", err)
	}
	return nil
}
