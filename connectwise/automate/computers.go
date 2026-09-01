package automate

import (
	"context"
	"encoding/json"
	"fmt"
)

func computerEndpoint(computerID int, sub string) string {
	return fmt.Sprintf("cwa/api/v1/Computers/%d/%s", computerID, sub)
}

// ListComputers returns managed computers (GET /cwa/api/v1/Computers). The API
// has no single-computer GET; filter by id with the options.ids query param.
func (c *Client) ListComputers(ctx context.Context, params map[string]string) ([]Computer, error) {
	result, err := c.Get[[]Computer](ctx, "cwa/api/v1/Computers", params)
	if err != nil {
		return nil, fmt.Errorf("list computers: %w", err)
	}
	return *result, nil
}

// GetComputerDevices lists a computer's hardware devices
// (GET /cwa/api/v1/Computers/{computerId}/Devices).
func (c *Client) GetComputerDevices(ctx context.Context, computerID int, params map[string]string) ([]ComputerDevice, error) {
	result, err := c.Get[[]ComputerDevice](ctx, computerEndpoint(computerID, "Devices"), params)
	if err != nil {
		return nil, fmt.Errorf("get computer devices: %w", err)
	}
	return *result, nil
}

// GetComputerServices lists a computer's Windows services
// (GET /cwa/api/v1/Computers/{computerId}/Services).
func (c *Client) GetComputerServices(ctx context.Context, computerID int, params map[string]string) ([]ComputerService, error) {
	result, err := c.Get[[]ComputerService](ctx, computerEndpoint(computerID, "Services"), params)
	if err != nil {
		return nil, fmt.Errorf("get computer services: %w", err)
	}
	return *result, nil
}

// GetComputerSoftware lists a computer's installed software
// (GET /cwa/api/v1/Computers/{computerId}/Software).
func (c *Client) GetComputerSoftware(ctx context.Context, computerID int, params map[string]string) ([]ComputerSoftware, error) {
	result, err := c.Get[[]ComputerSoftware](ctx, computerEndpoint(computerID, "Software"), params)
	if err != nil {
		return nil, fmt.Errorf("get computer software: %w", err)
	}
	return *result, nil
}

// GetComputerOperatingSystem returns a computer's OS details
// (GET /cwa/api/v1/Computers/{computerId}/OperatingSystem).
func (c *Client) GetComputerOperatingSystem(ctx context.Context, computerID int) (*ComputerOperatingSystem, error) {
	result, err := c.Get[ComputerOperatingSystem](ctx, computerEndpoint(computerID, "OperatingSystem"), nil)
	if err != nil {
		return nil, fmt.Errorf("get computer operating system: %w", err)
	}
	return result, nil
}

// GetComputerBios returns a computer's BIOS details
// (GET /cwa/api/v1/Computers/{computerId}/bios).
func (c *Client) GetComputerBios(ctx context.Context, computerID int) (*ComputerBios, error) {
	result, err := c.Get[ComputerBios](ctx, computerEndpoint(computerID, "bios"), nil)
	if err != nil {
		return nil, fmt.Errorf("get computer bios: %w", err)
	}
	return result, nil
}

// GetComputerPatchingStats returns a computer's patch compliance summary
// (GET /cwa/api/v1/Computers/{computerId}/PatchingStats).
func (c *Client) GetComputerPatchingStats(ctx context.Context, computerID int) (*ComputerPatchingStats, error) {
	result, err := c.Get[ComputerPatchingStats](ctx, computerEndpoint(computerID, "PatchingStats"), nil)
	if err != nil {
		return nil, fmt.Errorf("get computer patching stats: %w", err)
	}
	return result, nil
}

// GetComputerAlerts returns a computer's alerts
// (GET /cwa/api/v1/Computers/{computerId}/Alerts). The spec types the response
// as a generic object, so it is returned as raw JSON.
func (c *Client) GetComputerAlerts(ctx context.Context, computerID int, params map[string]string) (json.RawMessage, error) {
	result, err := c.Get[json.RawMessage](ctx, computerEndpoint(computerID, "Alerts"), params)
	if err != nil {
		return nil, fmt.Errorf("get computer alerts: %w", err)
	}
	return *result, nil
}

// GetComputerPatchJobs returns a computer's patch jobs
// (GET /cwa/api/v1/Computers/{computerId}/PatchJobs). The spec types the
// response as a generic object, so it is returned as raw JSON.
func (c *Client) GetComputerPatchJobs(ctx context.Context, computerID int, params map[string]string) (json.RawMessage, error) {
	result, err := c.Get[json.RawMessage](ctx, computerEndpoint(computerID, "PatchJobs"), params)
	if err != nil {
		return nil, fmt.Errorf("get computer patch jobs: %w", err)
	}
	return *result, nil
}

// GetComputerExtraFields returns the custom fields configured on a computer
// (GET /cwa/api/v1/Computers/{computerId}/ExtraFields). This endpoint is
// documented by ConnectWise but absent from the OpenAPI spec.
func (c *Client) GetComputerExtraFields(ctx context.Context, computerID int) ([]ExtraField, error) {
	result, err := c.Get[[]ExtraField](ctx, computerEndpoint(computerID, "ExtraFields"), nil)
	if err != nil {
		return nil, fmt.Errorf("get computer extra fields: %w", err)
	}
	return *result, nil
}

// PatchComputerExtraField updates a single extra field on a computer via JSON
// Patch operations, returning the updated field
// (PATCH /cwa/api/v1/Computers/{computerId}/ExtraFields/{extraFieldDefinitionId}).
// See PatchLocationExtraField for body examples. This endpoint is documented by
// ConnectWise but absent from the OpenAPI spec.
func (c *Client) PatchComputerExtraField(ctx context.Context, computerID, extraFieldDefinitionID int, patchOps []PatchOp) (*ExtraField, error) {
	endpoint := fmt.Sprintf("%s/%d", computerEndpoint(computerID, "ExtraFields"), extraFieldDefinitionID)
	result, err := c.Patch[ExtraField](ctx, endpoint, patchOps)
	if err != nil {
		return nil, fmt.Errorf("patch computer extra field: %w", err)
	}
	return result, nil
}

// ListComputerDrives lists disk drives across computers
// (GET /cwa/api/v1/Computers/Drives). Filter to a single computer with the
// options.condition query param.
func (c *Client) ListComputerDrives(ctx context.Context, params map[string]string) ([]ComputerDrive, error) {
	result, err := c.Get[[]ComputerDrive](ctx, "cwa/api/v1/Computers/Drives", params)
	if err != nil {
		return nil, fmt.Errorf("list computer drives: %w", err)
	}
	return *result, nil
}
