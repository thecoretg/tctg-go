package umbrella

import (
	"context"
	"fmt"
)

// RoamingComputer is the properties of a roaming computer in the organization.
type RoamingComputer struct {
	OriginID           int    `json:"originId"`
	Name               string `json:"name"`
	DeviceID           string `json:"deviceId"`
	Type               string `json:"type"`
	Status             string `json:"status"`
	SwgStatus          string `json:"swgStatus"`
	LastSyncStatus     string `json:"lastSyncStatus"`
	LastSyncSwgStatus  string `json:"lastSyncSwgStatus"`
	LastSync           string `json:"lastSync"`
	AppliedBundle      int    `json:"appliedBundle"`
	HasIPBlocking      bool   `json:"hasIpBlocking"`
	Version            string `json:"version"`
	OSVersion          string `json:"osVersion"`
	OSVersionName      string `json:"osVersionName"`
	AnyconnectDeviceID string `json:"anyconnectDeviceId,omitempty"`
}

// RoamingComputerUpdateRequest is the body for updating a roaming computer. The
// name must be 1–50 characters.
type RoamingComputerUpdateRequest struct {
	Name string `json:"name"`
}

// OrgInfo is the OrgInfo.json properties for deploying the Cisco Secure Client
// on user devices in the organization.
type OrgInfo struct {
	OrganizationID int    `json:"organizationId"`
	Fingerprint    string `json:"fingerprint"`
	UserID         int    `json:"userId"`
}

// ListRoamingComputers lists a single page of roaming computers. Supported
// params are "page", "limit" (max 100), "name", "status", "swgStatus",
// "lastSyncBefore", and "lastSyncAfter".
func (c *Client) ListRoamingComputers(ctx context.Context, params map[string]string) ([]RoamingComputer, error) {
	result, err := get[[]RoamingComputer](ctx, c, deploymentsURL("roamingcomputers"), params)
	if err != nil {
		return nil, fmt.Errorf("list roaming computers: %w", err)
	}
	return *result, nil
}

// GetRoamingComputer gets a roaming computer by device ID.
func (c *Client) GetRoamingComputer(ctx context.Context, deviceID string) (*RoamingComputer, error) {
	url := deploymentsURL(fmt.Sprintf("roamingcomputers/%s", deviceID))
	result, err := get[RoamingComputer](ctx, c, url, nil)
	if err != nil {
		return nil, fmt.Errorf("get roaming computer: %w", err)
	}
	return result, nil
}

// UpdateRoamingComputer updates the name of a roaming computer by device ID.
func (c *Client) UpdateRoamingComputer(ctx context.Context, deviceID string, body RoamingComputerUpdateRequest) (*RoamingComputer, error) {
	url := deploymentsURL(fmt.Sprintf("roamingcomputers/%s", deviceID))
	result, err := put[RoamingComputer](ctx, c, url, body)
	if err != nil {
		return nil, fmt.Errorf("update roaming computer: %w", err)
	}
	return result, nil
}

// DeleteRoamingComputer deletes a roaming computer by device ID.
func (c *Client) DeleteRoamingComputer(ctx context.Context, deviceID string) error {
	url := deploymentsURL(fmt.Sprintf("roamingcomputers/%s", deviceID))
	if err := del(ctx, c, url); err != nil {
		return fmt.Errorf("delete roaming computer: %w", err)
	}
	return nil
}

// GetOrganizationInfo gets the OrgInfo.json properties for deploying the Cisco
// Secure Client in the organization.
func (c *Client) GetOrganizationInfo(ctx context.Context) (*OrgInfo, error) {
	result, err := get[OrgInfo](ctx, c, deploymentsURL("roamingcomputers/orgInfo"), nil)
	if err != nil {
		return nil, fmt.Errorf("get organization info: %w", err)
	}
	return result, nil
}
