package threatdown

import (
	"context"
	"fmt"
)

type Site struct {
	ID                     string `json:"id"`
	AccountID              string `json:"accountid,omitempty"`
	CompanyName            string `json:"company_name"`
	AccountOwner           []User `json:"account_owner"`
	FirstName              string `json:"firstname"`
	LastName               string `json:"lastname"`
	Email                  string `json:"email"`
	IsRemoved              bool   `json:"is_removed"`
	CreatedDate            string `json:"createddate,omitempty"`
	CreatedByID            string `json:"createdbyid,omitempty"`
	LastModifiedByID       string `json:"lastmodifiedbyid,omitempty"`
	LastModifiedDate       string `json:"lastmodifieddate,omitempty"`
	NebulaAccountStatus    string `json:"nebula_account_status,omitempty"`
	NebulaAccountID        string `json:"nebula_account_id,omitempty"`
	NebulaAccountToken     string `json:"nebula_account_token,omitempty"`
	NoSubscription         bool   `json:"no_subscription"`
	CloudEvaluation        bool   `json:"cloud_evaluation"`
	Utility                bool   `json:"utility"`
	NFR                    bool   `json:"nfr"`
	BillingDuration        string `json:"billing_duration,omitempty"`
	BillingDate            int    `json:"billing_date"`
	WorkstationInstalled   bool   `json:"workstation_installed"`
	ServerInstalled        bool   `json:"server_installed"`
	MDREnabled             bool   `json:"mdr_enabled"`
	AutoConvertTrialToPaid bool   `json:"auto_convert_trial_to_paid"`
	BillingType            string `json:"billing_type,omitempty"`
	AccountStatus          string `json:"account_status,omitempty"`
}

// SiteInput is used for both creating and updating a site.
type SiteInput struct {
	ID           string   `json:"id,omitempty"`
	CompanyName  string   `json:"company_name,omitempty"`
	FirstName    string   `json:"firstname,omitempty"`
	LastName     string   `json:"lastname,omitempty"`
	Email        string   `json:"email,omitempty"`
	AccountOwner []string `json:"account_owner,omitempty"`
	SiteEndDate  string   `json:"site_end_date,omitempty"`
}

type sitesResp struct {
	Sites []Site `json:"sites"`
}

func (c *Client) CreateSite(ctx context.Context, input SiteInput) (*Site, error) {
	result, err := post[SiteInput](ctx, c, endpointURLV1("sites"), input)
	if err != nil {
		return nil, fmt.Errorf("create site: %w", err)
	}

	site, err := c.GetSite(ctx, result.ID)
	if err != nil {
		return nil, fmt.Errorf("getting site from id: %w", err)
	}
	return site, nil
}

func (c *Client) ListSites(ctx context.Context) ([]Site, error) {
	result, err := get[sitesResp](ctx, c, endpointURLV1("sites"), nil)
	if err != nil {
		return nil, fmt.Errorf("list sites: %w", err)
	}
	return result.Sites, nil
}

func (c *Client) GetSite(ctx context.Context, id string) (*Site, error) {
	site, err := get[Site](ctx, c, endpointURLV1("sites/"+id), nil)
	if err != nil {
		return nil, fmt.Errorf("get site: %w", err)
	}
	site.ID = id
	return site, nil
}

func (c *Client) UpdateSite(ctx context.Context, id string, input SiteInput) (*Site, error) {
	existing, err := c.GetSite(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting existing site for update: %w", err)
	}
	if input.CompanyName == "" {
		input.CompanyName = existing.CompanyName
	}
	if input.FirstName == "" {
		input.FirstName = existing.FirstName
	}
	if input.LastName == "" {
		input.LastName = existing.LastName
	}
	if input.Email == "" {
		input.Email = existing.Email
	}
	if input.AccountOwner == nil {
		owners := make([]string, len(existing.AccountOwner))
		for i, u := range existing.AccountOwner {
			owners[i] = u.ID
		}
		input.AccountOwner = owners
	}

	_, err = put[SiteInput](ctx, c, endpointURLV1("sites/"+id), input)
	if err != nil {
		return nil, fmt.Errorf("update site: %w", err)
	}

	site, err := c.GetSite(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting site from id: %w", err)
	}
	return site, nil
}

func (c *Client) DeleteSite(ctx context.Context, id string) error {
	if err := del(ctx, c, endpointURLV1("sites/"+id)); err != nil {
		return fmt.Errorf("delete site: %w", err)
	}
	return nil
}

func (c *Client) GetSiteByNebulaAccountID(ctx context.Context, accountID string) (*Site, error) {
	site, err := get[Site](ctx, c, endpointURLV1("sites/nebula-accounts/"+accountID), nil)
	if err != nil {
		return nil, fmt.Errorf("get site by nebula account id: %w", err)
	}
	return site, nil
}
