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

func (c *Client) PostCompany(ctx context.Context, company *Company) (*Company, error) {
	return post[Company](ctx, c, "company/companies", company)
}

func (c *Client) ListCompanies(ctx context.Context, params map[string]string) ([]Company, error) {
	return getMany[Company](ctx, c, "company/companies", params)
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

func (c *Client) ListCompanyTypes(ctx context.Context, params map[string]string) ([]CompanyType, error) {
	return getMany[CompanyType](ctx, c, "company/companies/types", params)
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
