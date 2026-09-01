package threatdown

import (
	"cmp"
	"context"
	"fmt"
)

type Site struct {
	ID          string `json:"id"`
	AccountID   string `json:"accountid,omitempty"`
	CompanyName string `json:"company_name"`

	// AccountOwner is the built-in list of owner users from ThreatDown.
	// It is a massive list of objects which is returned for each site. For a shortened
	// list of emails, set shortenOwner to true.
	AccountOwner []User `json:"account_owner"`

	// OwnerEmails is a list of owner emails; this is a helper within tctg-go and is not
	// an item returned by the API. It is meant to drastically reduce the size of ListSites
	// outputs since it is likely to include hundreds of the same massive list.
	OwnerEmails []string `json:"owner_emails,omitempty"`

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
	result, err := c.Post[SiteInput](ctx, endpointURLV1("sites"), input)
	if err != nil {
		return nil, fmt.Errorf("create site: %w", err)
	}

	site, err := c.GetSite(ctx, result.ID, false)
	if err != nil {
		return nil, fmt.Errorf("getting site from id: %w", err)
	}
	return site, nil
}

func (c *Client) ListSites(ctx context.Context, shortenOwner bool) ([]Site, error) {
	result, err := c.Get[sitesResp](ctx, endpointURLV1("sites"), nil)
	if err != nil {
		return nil, fmt.Errorf("list sites: %w", err)
	}
	if shortenOwner {
		for i := range result.Sites {
			s := &result.Sites[i]
			emails := make([]string, len(s.AccountOwner))
			for j, u := range s.AccountOwner {
				emails[j] = u.Email
			}
			s.OwnerEmails = emails
			s.AccountOwner = nil
		}
	}
	return result.Sites, nil
}

func (c *Client) GetSite(ctx context.Context, id string, shortenOwner bool) (*Site, error) {
	site, err := c.Get[Site](ctx, endpointURLV1("sites/"+id), nil)
	if err != nil {
		return nil, fmt.Errorf("get site: %w", err)
	}

	// The GET single endpoint returns a different ID format than POST/LIST.
	// Overwrite with the ID used in the request so callers always have the
	// hex-encoded form required for subsequent API calls.
	site.ID = id
	if shortenOwner {
		emails := make([]string, len(site.AccountOwner))
		for i, u := range site.AccountOwner {
			emails[i] = u.Email
		}
		site.OwnerEmails = emails
		site.AccountOwner = nil
	}
	return site, nil
}

func (c *Client) UpdateSite(ctx context.Context, id string, input SiteInput) (*Site, error) {
	existing, err := c.GetSite(ctx, id, false)
	if err != nil {
		return nil, fmt.Errorf("getting existing site for update: %w", err)
	}

	// the following values are required, even for a put, so use existing if
	// a zero value is in the input.
	input.CompanyName = cmp.Or(input.CompanyName, existing.CompanyName)
	input.FirstName = cmp.Or(input.FirstName, existing.FirstName)
	input.LastName = cmp.Or(input.LastName, existing.LastName)
	input.Email = cmp.Or(input.Email, existing.Email)

	// same as above, but can't use cmp with slices
	if input.AccountOwner == nil {
		owners := make([]string, len(existing.AccountOwner))
		for i, u := range existing.AccountOwner {
			owners[i] = u.ID
		}
		input.AccountOwner = owners
	}

	_, err = c.Put[SiteInput](ctx, endpointURLV1("sites/"+id), input)
	if err != nil {
		return nil, fmt.Errorf("update site: %w", err)
	}

	site, err := c.GetSite(ctx, id, false)
	if err != nil {
		return nil, fmt.Errorf("getting site from id: %w", err)
	}
	return site, nil
}

func (c *Client) DeleteSite(ctx context.Context, id string) error {
	if err := c.Delete(ctx, endpointURLV1("sites/"+id)); err != nil {
		return fmt.Errorf("delete site: %w", err)
	}
	return nil
}

func (c *Client) GetSiteByNebulaAccountID(ctx context.Context, accountID string) (*Site, error) {
	site, err := c.Get[Site](ctx, endpointURLV1("sites/nebula-accounts/"+accountID), nil)
	if err != nil {
		return nil, fmt.Errorf("get site by nebula account id: %w", err)
	}
	return site, nil
}
