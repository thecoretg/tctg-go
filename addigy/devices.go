package addigy

import (
	"context"
	"fmt"
)

// DeviceFilter is the request body for device search (device_entities.DeviceFilter).
type DeviceFilter struct {
	DesiredFactIdentifiers []string     `json:"desired_fact_identifiers,omitempty"`
	Page                   int          `json:"page,omitempty"`
	PerPage                int          `json:"per_page,omitempty"`
	Query                  *QueryFilter `json:"query,omitempty"`
	SortDirection          string       `json:"sort_direction,omitempty"`
	SortField              string       `json:"sort_field,omitempty"`
}

// QueryFilter narrows a device search (device_entities.QueryFilter).
type QueryFilter struct {
	Filters   []AuditFilter `json:"filters,omitempty"`
	PolicyID  string        `json:"policy_id,omitempty"`
	SearchAny string        `json:"search_any,omitempty"`
}

// AuditFilter is a single device-fact filter (device_entities.AuditFilter).
type AuditFilter struct {
	AuditField string `json:"audit_field,omitempty"`
	Operation  string `json:"operation,omitempty"`
	RangeValue any    `json:"range_value,omitempty"`
	Type       string `json:"type,omitempty"`
	Value      any    `json:"value,omitempty"`
}

// DeviceAuditResponse is the paginated result of a device search
// (device_entities.DeviceAuditResponse).
type DeviceAuditResponse struct {
	Items    []DeviceAudit `json:"items"`
	Metadata Metadata      `json:"metadata"`
}

// DeviceAudit is a single device with its requested facts (device_entities.DeviceAudit).
type DeviceAudit struct {
	AgentAuditDate string                `json:"agent_audit_date,omitempty"`
	AgentID        string                `json:"agentid,omitempty"`
	AuditDate      string                `json:"audit_date,omitempty"`
	Facts          map[string]DeviceFact `json:"facts,omitempty"`
	OrgID          string                `json:"orgid,omitempty"`
}

// DeviceFact is a single fact value on a device (device_entities.DeviceFact).
type DeviceFact struct {
	ErrorMsg string `json:"error_msg,omitempty"`
	Type     string `json:"type,omitempty"`
	Value    any    `json:"value,omitempty"`
}

// AssignRequest assigns or unassigns devices to/from policies
// (device_entities.AssignRequest).
type AssignRequest struct {
	AgentIDs  []string `json:"agent_ids"`
	PolicyIDs []string `json:"policy_ids"`
}

// ManagedUser is a managed admin account on a device
// (device_users_service.ManagedUserResponse).
type ManagedUser struct {
	AccountName      string `json:"accountName,omitempty"`
	AgentID          string `json:"agentId,omitempty"`
	CreatedDate      string `json:"createdDate,omitempty"`
	FullName         string `json:"fullName,omitempty"`
	IsRevealed       bool   `json:"isRevealed,omitempty"`
	LastRevealDate   string `json:"lastRevealDate,omitempty"`
	LastRotationDate string `json:"lastRotationDate,omitempty"`
	ManagementType   string `json:"managementType,omitempty"`
	Password         string `json:"password,omitempty"`
	UpdatedDate      string `json:"updatedDate,omitempty"`
}

// SearchDevices runs a universal device search scoped to the client's
// organization (POST /o/{organization_id}/devices).
func (c *Client) SearchDevices(ctx context.Context, filter DeviceFilter) (*DeviceAuditResponse, error) {
	url, err := c.orgURL("devices")
	if err != nil {
		return nil, err
	}
	result, err := c.Post[DeviceAuditResponse](ctx, url, filter)
	if err != nil {
		return nil, fmt.Errorf("search devices: %w", err)
	}
	return result, nil
}

// AssignDevices assigns devices to policies
// (POST /o/{organization_id}/devices/assign).
func (c *Client) AssignDevices(ctx context.Context, req AssignRequest) error {
	url, err := c.orgURL("devices/assign")
	if err != nil {
		return err
	}
	if _, err := c.Post[struct{}](ctx, url, req); err != nil {
		return fmt.Errorf("assign devices: %w", err)
	}
	return nil
}

// UnassignDevices unassigns devices from policies
// (POST /o/{organization_id}/devices/unassign).
func (c *Client) UnassignDevices(ctx context.Context, req AssignRequest) error {
	url, err := c.orgURL("devices/unassign")
	if err != nil {
		return err
	}
	if _, err := c.Post[struct{}](ctx, url, req); err != nil {
		return fmt.Errorf("unassign devices: %w", err)
	}
	return nil
}

// GetDevicePolicyAssignments returns the policy IDs assigned to a device
// (GET /o/{organization_id}/devices/{agent_id}/policy-assignments).
func (c *Client) GetDevicePolicyAssignments(ctx context.Context, agentID string) ([]string, error) {
	url, err := c.orgURL("devices/" + agentID + "/policy-assignments")
	if err != nil {
		return nil, err
	}
	result, err := c.Get[[]string](ctx, url, nil)
	if err != nil {
		return nil, fmt.Errorf("get device policy assignments: %w", err)
	}
	return *result, nil
}

// GetManagedUsers returns the managed admin accounts on a device
// (GET /o/{organization_id}/devices/{agent_id}/managed-users).
func (c *Client) GetManagedUsers(ctx context.Context, agentID string) ([]ManagedUser, error) {
	url, err := c.orgURL("devices/" + agentID + "/managed-users")
	if err != nil {
		return nil, err
	}
	result, err := c.Get[[]ManagedUser](ctx, url, nil)
	if err != nil {
		return nil, fmt.Errorf("get managed users: %w", err)
	}
	return *result, nil
}

// RotateManagedUserPassword rotates the password for a managed admin account on
// a device (PATCH /o/{organization_id}/devices/{agent_id}/managed-users).
func (c *Client) RotateManagedUserPassword(ctx context.Context, agentID, accountName string) (*ManagedUser, error) {
	url, err := c.orgURL("devices/" + agentID + "/managed-users")
	if err != nil {
		return nil, err
	}
	body := map[string]string{"accountName": accountName}
	result, err := c.Patch[ManagedUser](ctx, url, body)
	if err != nil {
		return nil, fmt.Errorf("rotate managed user password: %w", err)
	}
	return result, nil
}

// RevealManagedUserPassword reveals the password for a managed admin account on
// a device (GET /o/{organization_id}/devices/{agent_id}/managed-users/{account_name}/password).
func (c *Client) RevealManagedUserPassword(ctx context.Context, agentID, accountName string) (*ManagedUser, error) {
	url, err := c.orgURL("devices/" + agentID + "/managed-users/" + accountName + "/password")
	if err != nil {
		return nil, err
	}
	result, err := c.Get[ManagedUser](ctx, url, nil)
	if err != nil {
		return nil, fmt.Errorf("reveal managed user password: %w", err)
	}
	return result, nil
}

// RemoveDevice removes a device by serial number
// (DELETE /o/{organization_id}/devices/{sn}).
func (c *Client) RemoveDevice(ctx context.Context, serialNumber string) error {
	url, err := c.orgURL("devices/" + serialNumber)
	if err != nil {
		return err
	}
	if err := c.Delete(ctx, url, nil); err != nil {
		return fmt.Errorf("remove device: %w", err)
	}
	return nil
}
