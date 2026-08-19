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
	return post[Company](ctx, c, "company/companies", company)
}

// ListCompanies returns companies across every page unless WithLimit caps the
// result. On a large tenant an uncapped call is many round trips, so pass
// ConnectWise conditions in params to narrow it, for example
// {"conditions": `status/name="Active" and deletedFlag=false`}.
func (c *Client) ListCompanies(ctx context.Context, params map[string]string, opts ...ListOption) ([]Company, error) {
	return getMany[Company](ctx, c, "company/companies", params, opts...)
}

func (c *Client) GetCompany(ctx context.Context, companyID int, params map[string]string) (*Company, error) {
	return get[Company](ctx, c, companyIDEndpoint(companyID), params)
}

func (c *Client) PutCompany(ctx context.Context, companyID int, company *Company) (*Company, error) {
	return put[Company](ctx, c, companyIDEndpoint(companyID), company)
}

func (c *Client) PatchCompany(ctx context.Context, companyID int, patchOps []PatchOp) (*Company, error) {
	return patch[Company](ctx, c, companyIDEndpoint(companyID), patchOps)
}

func (c *Client) DeleteCompany(ctx context.Context, companyID int) error {
	return del(ctx, c, companyIDEndpoint(companyID))
}

func (c *Client) PostCompanyType(ctx context.Context, companyType *CompanyType) (*CompanyType, error) {
	return post[CompanyType](ctx, c, "company/companies/types", companyType)
}

func (c *Client) ListCompanyTypes(ctx context.Context, params map[string]string, opts ...ListOption) ([]CompanyType, error) {
	return getMany[CompanyType](ctx, c, "company/companies/types", params, opts...)
}

func (c *Client) GetCompanyType(ctx context.Context, typeID int, params map[string]string) (*CompanyType, error) {
	return get[CompanyType](ctx, c, companyTypeIDEndpoint(typeID), params)
}

func (c *Client) PutCompanyType(ctx context.Context, typeID int, companyType *CompanyType) (*CompanyType, error) {
	return put[CompanyType](ctx, c, companyTypeIDEndpoint(typeID), companyType)
}

func (c *Client) PatchCompanyType(ctx context.Context, typeID int, patchOps []PatchOp) (*CompanyType, error) {
	return patch[CompanyType](ctx, c, companyTypeIDEndpoint(typeID), patchOps)
}

func (c *Client) DeleteCompanyType(ctx context.Context, typeID int) error {
	return del(ctx, c, companyTypeIDEndpoint(typeID))
}

func (c *Client) PostCompanyTypeAssociation(ctx context.Context, association *CompanyTypeAssociation) (*CompanyTypeAssociation, error) {
	return post[CompanyTypeAssociation](ctx, c, "company/companyTypeAssociations", association)
}

func (c *Client) ListCompanyTypeAssociations(ctx context.Context, params map[string]string, opts ...ListOption) ([]CompanyTypeAssociation, error) {
	return getMany[CompanyTypeAssociation](ctx, c, "company/companyTypeAssociations", params, opts...)
}

func (c *Client) GetCompanyTypeAssociation(ctx context.Context, associationID int, params map[string]string) (*CompanyTypeAssociation, error) {
	return get[CompanyTypeAssociation](ctx, c, companyTypeAssociationsIDEndpoint(associationID), params)
}

func (c *Client) PutCompanyTypeAssociation(ctx context.Context, associationID int, association *CompanyTypeAssociation) (*CompanyTypeAssociation, error) {
	return put[CompanyTypeAssociation](ctx, c, companyTypeAssociationsIDEndpoint(associationID), association)
}

func (c *Client) PatchCompanyTypeAssociation(ctx context.Context, associationID int, patchOps []PatchOp) (*CompanyTypeAssociation, error) {
	return patch[CompanyTypeAssociation](ctx, c, companyTypeAssociationsIDEndpoint(associationID), patchOps)
}

func (c *Client) DeleteCompanyTypeAssociation(ctx context.Context, associationID int) error {
	return del(ctx, c, companyTypeAssociationsIDEndpoint(associationID))
}

func (c *Client) PostCompanyTypeAssociationForCompany(ctx context.Context, association *CompanyTypeAssociation, companyID int) (*CompanyTypeAssociation, error) {
	return post[CompanyTypeAssociation](ctx, c, companyTypeAssociationEndpoint(companyID), association)
}

func (c *Client) ListCompanyTypeAssociationsForCompany(ctx context.Context, params map[string]string, companyID int, opts ...ListOption) ([]CompanyTypeAssociation, error) {
	return getMany[CompanyTypeAssociation](ctx, c, companyTypeAssociationEndpoint(companyID), params, opts...)
}

func (c *Client) GetCompanyTypeAssociationForCompany(ctx context.Context, associationID int, params map[string]string, companyID int) (*CompanyTypeAssociation, error) {
	return get[CompanyTypeAssociation](ctx, c, companyTypeAssociationIDEndpoint(companyID, associationID), params)
}

func (c *Client) PutCompanyTypeAssociationForCompany(ctx context.Context, associationID int, association *CompanyTypeAssociation, companyID int) (*CompanyTypeAssociation, error) {
	return put[CompanyTypeAssociation](ctx, c, companyTypeAssociationIDEndpoint(companyID, associationID), association)
}

func (c *Client) PatchCompanyTypeAssociationForCompany(ctx context.Context, associationID int, patchOps []PatchOp, companyID int) (*CompanyTypeAssociation, error) {
	return patch[CompanyTypeAssociation](ctx, c, companyTypeAssociationIDEndpoint(companyID, associationID), patchOps)
}

func (c *Client) DeleteCompanyTypeAssociationForCompany(ctx context.Context, associationID int, companyID int) error {
	return del(ctx, c, companyTypeAssociationIDEndpoint(companyID, associationID))
}
