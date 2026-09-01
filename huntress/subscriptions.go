package huntress

import (
	"context"
	"fmt"
)

// SubscriptionSchedule is a scheduled change to a subscription.
type SubscriptionSchedule struct {
	ID          int64  `json:"id,omitempty"`
	Minimum     int64  `json:"minimum,omitempty"`
	Maximum     int64  `json:"maximum,omitempty"`
	Status      string `json:"status,omitempty"`
	TargetPrice int64  `json:"target_price,omitempty"`
	Months      int64  `json:"months,omitempty"`
	PromoUnits  int64  `json:"promo_units,omitempty"`
	StartsAt    string `json:"starts_at,omitempty"`
	EndsAt      string `json:"ends_at,omitempty"`
}

// Subscription represents a reseller subscription.
type Subscription struct {
	ID                   int64                  `json:"id,omitempty"`
	Account              map[string]any         `json:"account,omitempty"`
	Product              string                 `json:"product,omitempty"`
	Status               string                 `json:"status,omitempty"`
	MinimumUsage         int64                  `json:"minimum_usage,omitempty"`
	TerminalMinimumUsage int64                  `json:"terminal_minimum_usage,omitempty"`
	BillingInterval      string                 `json:"billing_interval,omitempty"`
	EffectiveDate        string                 `json:"effective_date,omitempty"`
	RenewalDate          string                 `json:"renewal_date,omitempty"`
	AutoRenew            bool                   `json:"auto_renew,omitempty"`
	Schedules            []SubscriptionSchedule `json:"schedules,omitempty"`
}

// SubscriptionCreationParameters is the request body for creating a subscription.
type SubscriptionCreationParameters struct {
	AccountID       int    `json:"account_id"`
	Product         string `json:"product"`
	Minimum         int    `json:"minimum"`
	PurchaseOrder   string `json:"purchase_order"`
	BillingInterval string `json:"billing_interval,omitempty"`
}

// SubscriptionUpdateParameters is the request body for updating a subscription.
type SubscriptionUpdateParameters struct {
	Minimum         int    `json:"minimum,omitempty"`
	PurchaseOrder   string `json:"purchase_order,omitempty"`
	BillingInterval string `json:"billing_interval,omitempty"`
	AdditionalUnits int    `json:"additional_units,omitempty"`
	AutoRenew       *bool  `json:"auto_renew,omitempty"`
}

// SubscriptionUpgradeParameters is the request body for upgrading a subscription.
type SubscriptionUpgradeParameters struct {
	Minimum       int    `json:"minimum"`
	PurchaseOrder string `json:"purchase_order,omitempty"`
	OnRenewal     *bool  `json:"on_renewal,omitempty"`
}

// SubscriptionsResponse is a page of subscriptions.
type SubscriptionsResponse struct {
	Subscriptions []Subscription `json:"subscriptions"`
	Pagination    Pagination     `json:"pagination"`
}

type subscriptionEnvelope struct {
	Subscription Subscription `json:"subscription"`
}

// ListSubscriptions returns a single page of reseller subscriptions (Reseller
// credentials only).
func (c *Client) ListSubscriptions(ctx context.Context, params map[string]string) (*SubscriptionsResponse, error) {
	result, err := c.Get[SubscriptionsResponse](ctx, endpointURL("reseller/subscriptions"), params)
	if err != nil {
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}
	return result, nil
}

// CreateSubscription creates a reseller subscription (Reseller credentials only).
func (c *Client) CreateSubscription(ctx context.Context, body SubscriptionCreationParameters) (*Subscription, error) {
	result, err := c.Post[Subscription](ctx, endpointURL("reseller/subscriptions"), body)
	if err != nil {
		return nil, fmt.Errorf("create subscription: %w", err)
	}
	return result, nil
}

// GetSubscription returns a single reseller subscription by ID (Reseller
// credentials only).
func (c *Client) GetSubscription(ctx context.Context, id int64) (*Subscription, error) {
	result, err := c.Get[subscriptionEnvelope](ctx, endpointURL(fmt.Sprintf("reseller/subscriptions/%d", id)), nil)
	if err != nil {
		return nil, fmt.Errorf("get subscription: %w", err)
	}
	return &result.Subscription, nil
}

// UpdateSubscription updates a reseller subscription (Reseller credentials only).
func (c *Client) UpdateSubscription(ctx context.Context, id int64, body SubscriptionUpdateParameters) (*Subscription, error) {
	result, err := c.Patch[Subscription](ctx, endpointURL(fmt.Sprintf("reseller/subscriptions/%d", id)), body)
	if err != nil {
		return nil, fmt.Errorf("update subscription: %w", err)
	}
	return result, nil
}

// UpgradeSubscription upgrades a reseller subscription (Reseller credentials only).
func (c *Client) UpgradeSubscription(ctx context.Context, id int64, body SubscriptionUpgradeParameters) (*Subscription, error) {
	result, err := c.Post[Subscription](ctx, endpointURL(fmt.Sprintf("reseller/subscriptions/%d/upgrade", id)), body)
	if err != nil {
		return nil, fmt.Errorf("upgrade subscription: %w", err)
	}
	return result, nil
}
