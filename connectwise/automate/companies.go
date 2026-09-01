package automate

import (
	"context"
	"fmt"
)

func companyEndpoint(companyID int) string {
	return fmt.Sprintf("cwa/api/v1/Clients/%d", companyID)
}

func (c *Client) ListCompanies(ctx context.Context, params map[string]string) ([]Company, error) {
	result, err := c.Get[[]Company](ctx, "cwa/api/v1/Clients", params)
	if err != nil {
		return nil, fmt.Errorf("list companies: %w", err)
	}
	return *result, nil
}

func (c *Client) GetCompany(ctx context.Context, companyID int, params map[string]string) (*Company, error) {
	result, err := c.Get[Company](ctx, companyEndpoint(companyID), params)
	if err != nil {
		return nil, fmt.Errorf("get company: %w", err)
	}
	return result, nil
}

func (c *Client) PostCompany(ctx context.Context, company *Company) (*Company, error) {
	result, err := c.Post[Company](ctx, "cwa/api/v1/Clients", company)
	if err != nil {
		return nil, fmt.Errorf("post company: %w", err)
	}
	return result, nil
}

func (c *Client) PutCompany(ctx context.Context, companyID int, company *Company) (*Company, error) {
	result, err := c.Put[Company](ctx, companyEndpoint(companyID), company)
	if err != nil {
		return nil, fmt.Errorf("put company: %w", err)
	}
	return result, nil
}

func (c *Client) PatchCompany(ctx context.Context, companyID int, patchOps []PatchOp) (*Company, error) {
	result, err := c.Patch[Company](ctx, companyEndpoint(companyID), patchOps)
	if err != nil {
		return nil, fmt.Errorf("patch company: %w", err)
	}
	return result, nil
}

func (c *Client) DeleteCompany(ctx context.Context, companyID int) error {
	if err := c.Delete(ctx, companyEndpoint(companyID)); err != nil {
		return fmt.Errorf("delete company: %w", err)
	}
	return nil
}

// GetCompanyExtraFields returns the custom fields configured on a company
// (GET /cwa/api/v1/Clients/{clientId}/ExtraFields).
func (c *Client) GetCompanyExtraFields(ctx context.Context, companyID int) ([]ExtraField, error) {
	result, err := c.Get[[]ExtraField](ctx, companyEndpoint(companyID)+"/ExtraFields", nil)
	if err != nil {
		return nil, fmt.Errorf("get company extra fields: %w", err)
	}
	return *result, nil
}

// PatchCompanyExtraField updates a single extra field on a company via JSON
// Patch operations, returning the updated field
// (PATCH /cwa/api/v1/Clients/{clientId}/ExtraFields/{extraFieldDefinitionId}).
// For example, to set a text field's value:
//
//	[]PatchOp{{Op: "replace", Path: "/TextFieldSettings/Value", Value: "new value"}}
//
// This endpoint is documented by ConnectWise but absent from the OpenAPI spec.
func (c *Client) PatchCompanyExtraField(ctx context.Context, companyID, extraFieldDefinitionID int, patchOps []PatchOp) (*ExtraField, error) {
	endpoint := fmt.Sprintf("%s/ExtraFields/%d", companyEndpoint(companyID), extraFieldDefinitionID)
	result, err := c.Patch[ExtraField](ctx, endpoint, patchOps)
	if err != nil {
		return nil, fmt.Errorf("patch company extra field: %w", err)
	}
	return result, nil
}

// ListCompanyDocuments returns the documents attached to a company
// (GET /cwa/api/v1/clients/{clientId}/documents).
func (c *Client) ListCompanyDocuments(ctx context.Context, companyID int, params map[string]string) ([]Document, error) {
	result, err := c.Get[[]Document](ctx, companyEndpoint(companyID)+"/documents", params)
	if err != nil {
		return nil, fmt.Errorf("list company documents: %w", err)
	}
	return *result, nil
}

// ListCompanyLicenses returns the managed licenses on a company
// (GET /cwa/api/v1/clients/{clientId}/licenses).
func (c *Client) ListCompanyLicenses(ctx context.Context, companyID int, params map[string]string) ([]ManagedLicense, error) {
	result, err := c.Get[[]ManagedLicense](ctx, companyEndpoint(companyID)+"/licenses", params)
	if err != nil {
		return nil, fmt.Errorf("list company licenses: %w", err)
	}
	return *result, nil
}

// PostCompanyLicense creates a managed license on a company
// (POST /cwa/api/v1/clients/{clientId}/licenses).
func (c *Client) PostCompanyLicense(ctx context.Context, companyID int, license *ManagedLicense) (*ManagedLicense, error) {
	result, err := c.Post[ManagedLicense](ctx, companyEndpoint(companyID)+"/licenses", license)
	if err != nil {
		return nil, fmt.Errorf("post company license: %w", err)
	}
	return result, nil
}

// ListCompanyProductKeys returns the product keys on a company
// (GET /cwa/api/v1/clients/{clientId}/productkeys).
func (c *Client) ListCompanyProductKeys(ctx context.Context, companyID int, params map[string]string) ([]ProductKey, error) {
	result, err := c.Get[[]ProductKey](ctx, companyEndpoint(companyID)+"/productkeys", params)
	if err != nil {
		return nil, fmt.Errorf("list company product keys: %w", err)
	}
	return *result, nil
}

// PostCompanyProductKey creates a product key on a company
// (POST /cwa/api/v1/clients/{clientId}/productkeys).
func (c *Client) PostCompanyProductKey(ctx context.Context, companyID int, key *ProductKey) (*ProductKey, error) {
	result, err := c.Post[ProductKey](ctx, companyEndpoint(companyID)+"/productkeys", key)
	if err != nil {
		return nil, fmt.Errorf("post company product key: %w", err)
	}
	return result, nil
}

func companyPermissionsEndpoint(companyID, userClassID int) string {
	return fmt.Sprintf("cwa/api/v1/clients/%d/permissions/%d", companyID, userClassID)
}

// GetCompanyPermissions returns the permissions a user class has on a company
// (GET /cwa/api/v1/clients/{clientId}/permissions/{userClassId}).
func (c *Client) GetCompanyPermissions(ctx context.Context, companyID, userClassID int) ([]string, error) {
	result, err := c.Get[[]string](ctx, companyPermissionsEndpoint(companyID, userClassID), nil)
	if err != nil {
		return nil, fmt.Errorf("get company permissions: %w", err)
	}
	return *result, nil
}

// PutCompanyPermissions replaces the permissions a user class has on a company
// (PUT /cwa/api/v1/clients/{clientId}/permissions/{userClassId}).
func (c *Client) PutCompanyPermissions(ctx context.Context, companyID, userClassID int, permissions []string) ([]string, error) {
	result, err := c.Put[[]string](ctx, companyPermissionsEndpoint(companyID, userClassID), permissions)
	if err != nil {
		return nil, fmt.Errorf("put company permissions: %w", err)
	}
	return *result, nil
}

// PostCompanyPermissions adds permissions for a user class on a company
// (POST /cwa/api/v1/clients/{clientId}/permissions/{userClassId}).
func (c *Client) PostCompanyPermissions(ctx context.Context, companyID, userClassID int, permissions []string) ([]string, error) {
	result, err := c.Post[[]string](ctx, companyPermissionsEndpoint(companyID, userClassID), permissions)
	if err != nil {
		return nil, fmt.Errorf("post company permissions: %w", err)
	}
	return *result, nil
}

// DeleteCompanyPermissions removes a user class's permissions on a company
// (DELETE /cwa/api/v1/clients/{clientId}/permissions/{userClassId}).
func (c *Client) DeleteCompanyPermissions(ctx context.Context, companyID, userClassID int) error {
	if err := c.Delete(ctx, companyPermissionsEndpoint(companyID, userClassID)); err != nil {
		return fmt.Errorf("delete company permissions: %w", err)
	}
	return nil
}
