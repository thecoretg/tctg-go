package automate

import (
	"context"
	"fmt"
)

func locationEndpoint(locationID int) string {
	return fmt.Sprintf("cwa/api/v1/Locations/%d", locationID)
}

func locationExtraFieldEndpoint(locationID, extraFieldDefinitionID int) string {
	return fmt.Sprintf("cwa/api/v1/Locations/%d/ExtraFields/%d", locationID, extraFieldDefinitionID)
}

// ListLocations returns all locations (GET /cwa/api/v1/Locations).
func (c *Client) ListLocations(ctx context.Context, params map[string]string) ([]Location, error) {
	result, err := get[[]Location](ctx, c, "cwa/api/v1/Locations", params)
	if err != nil {
		return nil, fmt.Errorf("list locations: %w", err)
	}
	return *result, nil
}

// PostLocation creates a location (POST /cwa/api/v1/Locations).
func (c *Client) PostLocation(ctx context.Context, location *Location) (*Location, error) {
	result, err := post[Location](ctx, c, "cwa/api/v1/Locations", location)
	if err != nil {
		return nil, fmt.Errorf("post location: %w", err)
	}
	return result, nil
}

// GetLocation returns a single location (GET /cwa/api/v1/Locations/{locationId}).
func (c *Client) GetLocation(ctx context.Context, locationID int, params map[string]string) (*Location, error) {
	result, err := get[Location](ctx, c, locationEndpoint(locationID), params)
	if err != nil {
		return nil, fmt.Errorf("get location: %w", err)
	}
	return result, nil
}

// PutLocation replaces a location (PUT /cwa/api/v1/Locations/{locationId}).
func (c *Client) PutLocation(ctx context.Context, locationID int, location *Location) (*Location, error) {
	result, err := put[Location](ctx, c, locationEndpoint(locationID), location)
	if err != nil {
		return nil, fmt.Errorf("put location: %w", err)
	}
	return result, nil
}

// PatchLocation applies JSON Patch operations to a location
// (PATCH /cwa/api/v1/Locations/{locationId}).
func (c *Client) PatchLocation(ctx context.Context, locationID int, patchOps []PatchOp) (*Location, error) {
	result, err := patch[Location](ctx, c, locationEndpoint(locationID), patchOps)
	if err != nil {
		return nil, fmt.Errorf("patch location: %w", err)
	}
	return result, nil
}

// DeleteLocation deletes a location (DELETE /cwa/api/v1/Locations/{locationId}).
func (c *Client) DeleteLocation(ctx context.Context, locationID int) error {
	if err := del(ctx, c, locationEndpoint(locationID)); err != nil {
		return fmt.Errorf("delete location: %w", err)
	}
	return nil
}

// GetLocationExtraFields returns the custom fields configured on a location
// (GET /cwa/api/v1/Locations/{locationId}/ExtraFields).
func (c *Client) GetLocationExtraFields(ctx context.Context, locationID int) ([]ExtraField, error) {
	result, err := get[[]ExtraField](ctx, c, locationEndpoint(locationID)+"/ExtraFields", nil)
	if err != nil {
		return nil, fmt.Errorf("get location extra fields: %w", err)
	}
	return *result, nil
}

// PatchLocationExtraField updates a single extra field on a location via JSON
// Patch operations, returning the updated field
// (PATCH /cwa/api/v1/Locations/{locationId}/ExtraFields/{extraFieldDefinitionId}).
// For example, to set a text field's value:
//
//	[]PatchOp{{Op: "replace", Path: "/TextFieldSettings/Value", Value: "new value"}}
func (c *Client) PatchLocationExtraField(ctx context.Context, locationID, extraFieldDefinitionID int, patchOps []PatchOp) (*ExtraField, error) {
	result, err := patch[ExtraField](ctx, c, locationExtraFieldEndpoint(locationID, extraFieldDefinitionID), patchOps)
	if err != nil {
		return nil, fmt.Errorf("patch location extra field: %w", err)
	}
	return result, nil
}

// ResetLocationExtraField resets a location's extra field to its default,
// returning the reset field
// (DELETE /cwa/api/v1/Locations/{locationId}/ExtraFields/{extraFieldDefinitionId}).
func (c *Client) ResetLocationExtraField(ctx context.Context, locationID, extraFieldDefinitionID int) (*ExtraField, error) {
	result, err := delReturn[ExtraField](ctx, c, locationExtraFieldEndpoint(locationID, extraFieldDefinitionID))
	if err != nil {
		return nil, fmt.Errorf("reset location extra field: %w", err)
	}
	return result, nil
}
