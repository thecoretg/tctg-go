package psa

import (
	"context"
	"fmt"
)

func companyIDEndpoint(companyID int) string {
	return fmt.Sprintf("company/companies/%d", companyID)
}

func companyTypeIDEndpoint(typeID int) string {
	return fmt.Sprintf("company/companies/types/%d", typeID)
}

func companyTypeAssociationEndpoint(companyID int) string {
	return fmt.Sprintf("%s/typeAssociations", companyIDEndpoint(companyID))
}

func companyTypeAssociationIDEndpoint(companyID, associationID int) string {
	return fmt.Sprintf("%s/%d", companyTypeAssociationEndpoint(companyID), associationID)
}

func companyTypeAssociationsIDEndpoint(associationID int) string {
	return fmt.Sprintf("company/companyTypeAssociations/%d", associationID)
}

func (c *Client) PostCompany(ctx context.Context, company *Company) (*Company, error) {
	return c.Post[Company](ctx, "company/companies", company)
}

// ListCompanies returns companies across every page unless WithLimit caps the
// result. On a large tenant an uncapped call is many round trips, so pass
// ConnectWise conditions in params to narrow it, for example
// {"conditions": `status/name="Active" and deletedFlag=false`}.
func (c *Client) ListCompanies(ctx context.Context, params map[string]string, opts ...ListOption) ([]Company, error) {
	return c.GetMany[Company](ctx, "company/companies", params, opts...)
}

func (c *Client) GetCompany(ctx context.Context, companyID int, params map[string]string) (*Company, error) {
	return c.Get[Company](ctx, companyIDEndpoint(companyID), params)
}

func (c *Client) PutCompany(ctx context.Context, companyID int, company *Company) (*Company, error) {
	return c.Put[Company](ctx, companyIDEndpoint(companyID), company)
}

func (c *Client) PatchCompany(ctx context.Context, companyID int, patchOps []PatchOp) (*Company, error) {
	return c.Patch[Company](ctx, companyIDEndpoint(companyID), patchOps)
}

func (c *Client) DeleteCompany(ctx context.Context, companyID int) error {
	return c.Delete(ctx, companyIDEndpoint(companyID))
}

func (c *Client) PostCompanyType(ctx context.Context, companyType *CompanyType) (*CompanyType, error) {
	return c.Post[CompanyType](ctx, "company/companies/types", companyType)
}

func (c *Client) ListCompanyTypes(ctx context.Context, params map[string]string, opts ...ListOption) ([]CompanyType, error) {
	return c.GetMany[CompanyType](ctx, "company/companies/types", params, opts...)
}

func (c *Client) GetCompanyType(ctx context.Context, typeID int, params map[string]string) (*CompanyType, error) {
	return c.Get[CompanyType](ctx, companyTypeIDEndpoint(typeID), params)
}

func (c *Client) PutCompanyType(ctx context.Context, typeID int, companyType *CompanyType) (*CompanyType, error) {
	return c.Put[CompanyType](ctx, companyTypeIDEndpoint(typeID), companyType)
}

func (c *Client) PatchCompanyType(ctx context.Context, typeID int, patchOps []PatchOp) (*CompanyType, error) {
	return c.Patch[CompanyType](ctx, companyTypeIDEndpoint(typeID), patchOps)
}

func (c *Client) DeleteCompanyType(ctx context.Context, typeID int) error {
	return c.Delete(ctx, companyTypeIDEndpoint(typeID))
}

func (c *Client) PostCompanyTypeAssociation(ctx context.Context, association *CompanyTypeAssociation) (*CompanyTypeAssociation, error) {
	return c.Post[CompanyTypeAssociation](ctx, "company/companyTypeAssociations", association)
}

func (c *Client) ListCompanyTypeAssociations(ctx context.Context, params map[string]string, opts ...ListOption) ([]CompanyTypeAssociation, error) {
	return c.GetMany[CompanyTypeAssociation](ctx, "company/companyTypeAssociations", params, opts...)
}

func (c *Client) GetCompanyTypeAssociation(ctx context.Context, associationID int, params map[string]string) (*CompanyTypeAssociation, error) {
	return c.Get[CompanyTypeAssociation](ctx, companyTypeAssociationsIDEndpoint(associationID), params)
}

func (c *Client) PutCompanyTypeAssociation(ctx context.Context, associationID int, association *CompanyTypeAssociation) (*CompanyTypeAssociation, error) {
	return c.Put[CompanyTypeAssociation](ctx, companyTypeAssociationsIDEndpoint(associationID), association)
}

func (c *Client) PatchCompanyTypeAssociation(ctx context.Context, associationID int, patchOps []PatchOp) (*CompanyTypeAssociation, error) {
	return c.Patch[CompanyTypeAssociation](ctx, companyTypeAssociationsIDEndpoint(associationID), patchOps)
}

func (c *Client) DeleteCompanyTypeAssociation(ctx context.Context, associationID int) error {
	return c.Delete(ctx, companyTypeAssociationsIDEndpoint(associationID))
}

func (c *Client) PostCompanyTypeAssociationForCompany(ctx context.Context, association *CompanyTypeAssociation, companyID int) (*CompanyTypeAssociation, error) {
	return c.Post[CompanyTypeAssociation](ctx, companyTypeAssociationEndpoint(companyID), association)
}

func (c *Client) ListCompanyTypeAssociationsForCompany(ctx context.Context, params map[string]string, companyID int, opts ...ListOption) ([]CompanyTypeAssociation, error) {
	return c.GetMany[CompanyTypeAssociation](ctx, companyTypeAssociationEndpoint(companyID), params, opts...)
}

func (c *Client) GetCompanyTypeAssociationForCompany(ctx context.Context, associationID int, params map[string]string, companyID int) (*CompanyTypeAssociation, error) {
	return c.Get[CompanyTypeAssociation](ctx, companyTypeAssociationIDEndpoint(companyID, associationID), params)
}

func (c *Client) PutCompanyTypeAssociationForCompany(ctx context.Context, associationID int, association *CompanyTypeAssociation, companyID int) (*CompanyTypeAssociation, error) {
	return c.Put[CompanyTypeAssociation](ctx, companyTypeAssociationIDEndpoint(companyID, associationID), association)
}

func (c *Client) PatchCompanyTypeAssociationForCompany(ctx context.Context, associationID int, patchOps []PatchOp, companyID int) (*CompanyTypeAssociation, error) {
	return c.Patch[CompanyTypeAssociation](ctx, companyTypeAssociationIDEndpoint(companyID, associationID), patchOps)
}

func (c *Client) DeleteCompanyTypeAssociationForCompany(ctx context.Context, associationID int, companyID int) error {
	return c.Delete(ctx, companyTypeAssociationIDEndpoint(companyID, associationID))
}
