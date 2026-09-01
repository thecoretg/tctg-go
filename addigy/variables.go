package addigy

import (
	"context"
	"fmt"
)

// Variable is a global variable (variable_entities.Variable). DefaultValue is
// untyped in the API and may be a string or other JSON value.
type Variable struct {
	Key          string `json:"key,omitempty"`
	Type         string `json:"type,omitempty"`
	DefaultValue any    `json:"default_value,omitempty"`
	CreatedDate  string `json:"created_date,omitempty"`
	UpdatedDate  string `json:"updated_date,omitempty"`
}

// NewVariableRequest is the body for creating a variable
// (variable_entities.NewVariableRequest). Type must be "string" or "secret".
type NewVariableRequest struct {
	Key          string `json:"key"`
	Type         string `json:"type"`
	DefaultValue any    `json:"default_value,omitempty"`
}

// VariableUpdateRequest is the body for updating a variable
// (variable_entities.VariableUpdateRequest).
type VariableUpdateRequest struct {
	Key          string `json:"key"`
	DefaultValue any    `json:"default_value,omitempty"`
}

// VariableValueResponse holds a resolved variable value
// (variable_entities.VariableValueResponse).
type VariableValueResponse struct {
	Value any `json:"value,omitempty"`
}

// PolicyValue is a per-policy variable value (variable_entities.PolicyValue).
type PolicyValue struct {
	PolicyID string `json:"policy_id"`
	Value    any    `json:"value,omitempty"`
}

// VariablePolicies lists the policy-level values set for a variable
// (variable_entities.VariablePolicies).
type VariablePolicies struct {
	VariableKey  string        `json:"variable_key"`
	PolicyValues []PolicyValue `json:"policy_values"`
	UpdatedDate  string        `json:"updated_date,omitempty"`
}

// VariablePolicy is the body for assigning a policy value to a variable
// (variable_entities.VariablePolicy).
type VariablePolicy struct {
	VariableKey string `json:"variable_key"`
	PolicyID    string `json:"policy_id"`
	Value       any    `json:"value,omitempty"`
}

// AssetVariableUsage reports where a variable is used
// (global_variables_service.AssetVariableUsage).
type AssetVariableUsage struct {
	VariableKey    string   `json:"variable_key,omitempty"`
	OrganizationID string   `json:"organization_id,omitempty"`
	AssetType      string   `json:"asset_type,omitempty"`
	AssetIDs       []string `json:"asset_ids,omitempty"`
}

// VariableFilter narrows a variable query (variable_entities.Filter).
type VariableFilter struct {
	Keys        []string `json:"keys,omitempty"`
	KeyContains string   `json:"key_contains,omitempty"`
}

// VariablesQueryRequest is the body for querying variables
// (variable_entities.VariablesQueryRequest).
type VariablesQueryRequest struct {
	Page          int             `json:"page,omitempty"`
	PerPage       int             `json:"per_page,omitempty"`
	Query         *VariableFilter `json:"query,omitempty"`
	SortDirection string          `json:"sort_direction,omitempty"`
	SortField     string          `json:"sort_field,omitempty"`
}

// VariablesQueryResponse is the paginated list of variables
// (variables.QueryVariablesResponse).
type VariablesQueryResponse struct {
	Items    []Variable `json:"items"`
	Metadata Metadata   `json:"metadata"`
}

// CreateVariable creates a variable (POST /o/{organization_id}/variables).
func (c *Client) CreateVariable(ctx context.Context, req NewVariableRequest) (*Variable, error) {
	url, err := c.orgURL("variables")
	if err != nil {
		return nil, err
	}
	result, err := c.Post[Variable](ctx, url, req)
	if err != nil {
		return nil, fmt.Errorf("create variable: %w", err)
	}
	return result, nil
}

// UpdateVariable updates a variable (PUT /o/{organization_id}/variables).
func (c *Client) UpdateVariable(ctx context.Context, req VariableUpdateRequest) (*Variable, error) {
	url, err := c.orgURL("variables")
	if err != nil {
		return nil, err
	}
	result, err := c.Put[Variable](ctx, url, req)
	if err != nil {
		return nil, fmt.Errorf("update variable: %w", err)
	}
	return result, nil
}

// DeleteVariable deletes a variable by key
// (DELETE /o/{organization_id}/variables).
func (c *Client) DeleteVariable(ctx context.Context, key string) error {
	url, err := c.orgURL("variables")
	if err != nil {
		return err
	}
	if err := c.Delete(ctx, url, map[string]string{"key": key}); err != nil {
		return fmt.Errorf("delete variable: %w", err)
	}
	return nil
}

// GetVariableValue returns a variable's value
// (GET /o/{organization_id}/variables/value).
func (c *Client) GetVariableValue(ctx context.Context, key string) (*VariableValueResponse, error) {
	url, err := c.orgURL("variables/value")
	if err != nil {
		return nil, err
	}
	result, err := c.Get[VariableValueResponse](ctx, url, map[string]string{"variable_key": key})
	if err != nil {
		return nil, fmt.Errorf("get variable value: %w", err)
	}
	return result, nil
}

// GetVariableUsage returns where a variable is used
// (GET /o/{organization_id}/variables/usage).
func (c *Client) GetVariableUsage(ctx context.Context, key string) ([]AssetVariableUsage, error) {
	url, err := c.orgURL("variables/usage")
	if err != nil {
		return nil, err
	}
	result, err := c.Get[[]AssetVariableUsage](ctx, url, map[string]string{"variable_key": key})
	if err != nil {
		return nil, fmt.Errorf("get variable usage: %w", err)
	}
	return *result, nil
}

// GetVariablePolicies returns the policy-level values for variables, optionally
// filtered by policy_id and variable_key
// (GET /o/{organization_id}/variables/policies).
func (c *Client) GetVariablePolicies(ctx context.Context, policyID, variableKey string) ([]VariablePolicies, error) {
	url, err := c.orgURL("variables/policies")
	if err != nil {
		return nil, err
	}
	params := map[string]string{}
	if policyID != "" {
		params["policy_id"] = policyID
	}
	if variableKey != "" {
		params["variable_key"] = variableKey
	}
	result, err := c.Get[[]VariablePolicies](ctx, url, params)
	if err != nil {
		return nil, fmt.Errorf("get variable policies: %w", err)
	}
	return *result, nil
}

// AssignVariablePolicy assigns a policy value to a variable
// (POST /o/{organization_id}/variables/policies).
func (c *Client) AssignVariablePolicy(ctx context.Context, req VariablePolicy) error {
	url, err := c.orgURL("variables/policies")
	if err != nil {
		return err
	}
	if _, err := c.Post[struct{}](ctx, url, req); err != nil {
		return fmt.Errorf("assign variable policy: %w", err)
	}
	return nil
}

// RemoveVariablePolicy removes a policy value from a variable
// (DELETE /o/{organization_id}/variables/policies).
func (c *Client) RemoveVariablePolicy(ctx context.Context, policyID, variableKey string) error {
	url, err := c.orgURL("variables/policies")
	if err != nil {
		return err
	}
	params := map[string]string{"policy_id": policyID, "variable_key": variableKey}
	if err := c.Delete(ctx, url, params); err != nil {
		return fmt.Errorf("remove variable policy: %w", err)
	}
	return nil
}

// GetVariablePolicyValue returns a variable's value for a specific policy
// (GET /o/{organization_id}/variables/policies/value).
func (c *Client) GetVariablePolicyValue(ctx context.Context, variableKey, policyID string) (*VariableValueResponse, error) {
	url, err := c.orgURL("variables/policies/value")
	if err != nil {
		return nil, err
	}
	params := map[string]string{"variable_key": variableKey, "policy_id": policyID}
	result, err := c.Get[VariableValueResponse](ctx, url, params)
	if err != nil {
		return nil, fmt.Errorf("get variable policy value: %w", err)
	}
	return result, nil
}

// QueryVariables returns variables filtered by key (POST /oa/variables/query).
func (c *Client) QueryVariables(ctx context.Context, req VariablesQueryRequest) (*VariablesQueryResponse, error) {
	result, err := c.Post[VariablesQueryResponse](ctx, c.url("oa/variables/query"), req)
	if err != nil {
		return nil, fmt.Errorf("query variables: %w", err)
	}
	return result, nil
}
