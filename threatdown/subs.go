package threatdown

import (
	"context"
	"errors"
	"fmt"
)

type SiteSubscription struct {
	Allocations SiteSubscriptionAllocations `json:"allocations"`
	MachineType MachineType                 `json:"machine_type"`
	Product     ProductType                 `json:"product"`
	TermLength  string                      `json:"term_length"`
	TermType    TermType                    `json:"term_type"`
}

type SiteSubscriptionAllocations struct {
	EndpointDetectionResponseForServers int `json:"edrs"`
	EndpointProtection                  int `json:"ep"`
	EndpointProtectionForServers        int `json:"eps"`
	EndpointProtectionResponse          int `json:"edr"`
	IncidentResponse                    int `json:"ir"`
	MobileSecurityForBusiness           int `json:"mob"`
}

type (
	ProductType string
	MachineType string
	TermType    string
)

const (
	ProductTypeIR                ProductType = "ir"
	ProductTypeEP                ProductType = "ep"
	ProductTypeEDR               ProductType = "edr"
	ProductTypeMOB               ProductType = "mob"
	MachineTypeWorkstation       MachineType = "ws"
	MachineTypeWorkstationServer MachineType = "ws-ser"
	TermTypePaid                 TermType    = "paid"
	TermTypeTrial                TermType    = "trial"
)

var ErrSubNotFound = errors.New("subscription not found")

func (c *Client) GetSiteSubscriptions(ctx context.Context, siteID string) ([]SiteSubscription, error) {
	subs, err := get[[]SiteSubscription](ctx, c, endpointURLV2(fmt.Sprintf("sites/%s/subscriptions", siteID)), nil)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrSubNotFound
		}
		return nil, fmt.Errorf("get site subscriptions: %w", err)
	}
	return *subs, nil
}

func (c *Client) CreateSiteSubscription(ctx context.Context, siteID string, sub SiteSubscription) ([]SiteSubscription, error) {
	subs, err := post[[]SiteSubscription](ctx, c, endpointURLV2(fmt.Sprintf("sites/%s/subscriptions", siteID)), sub)
	if err != nil {
		return nil, fmt.Errorf("create site subscription: %w", err)
	}
	return *subs, nil
}

func (c *Client) UpdateSiteSubscriptions(ctx context.Context, siteID string, subs []SiteSubscription) ([]SiteSubscription, error) {
	result, err := put[[]SiteSubscription](ctx, c, endpointURLV2(fmt.Sprintf("sites/%s/subscriptions", siteID)), subs)
	if err != nil {
		return nil, fmt.Errorf("update site subscriptions: %w", err)
	}
	return *result, nil
}
