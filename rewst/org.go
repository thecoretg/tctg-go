package rewst

import (
	"context"
	"errors"
	"fmt"
)

type Org struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	OrgSlug       string `json:"org_slug"`
	Domain        string `json:"domain"`
	ManagingOrgID string `json:"managing_org_id"`
	ManagingOrg   struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"managingOrg"`
	RocSiteID any  `json:"roc_site_id"`
	IsEnabled bool `json:"is_enabled"`
	Tags      []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"tags"`
}

type listOrgsResp struct {
	Items []Org   `json:"items"`
	Error *string `json:"error"`
}

func (c *Client) ListOrgs(ctx context.Context) ([]Org, error) {
	result, err := get[listOrgsResp](ctx, c, c.listOrgsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("list orgs: %w", err)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("list orgs: %w", errors.New(*result.Error))
	}

	return result.Items, nil
}
