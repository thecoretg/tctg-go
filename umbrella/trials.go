package umbrella

import (
	"context"
	"fmt"
)

// CustomerSubscription is the rich subscription/trial view of a customer,
// returned by the subscription-details and trial-extension endpoints. The
// createdAt/modifiedAt fields are epoch milliseconds; startsAt/endsAt are dates.
type CustomerSubscription struct {
	OrganizationID              int            `json:"organizationId"`
	OrganizationTypeID          int            `json:"organizationTypeId,omitempty"`
	OriginID                    int            `json:"originId,omitempty"`
	OrganizationName            string         `json:"organizationName,omitempty"`
	CreatedAt                   int64          `json:"createdAt,omitempty"`
	ModifiedAt                  int64          `json:"modifiedAt,omitempty"`
	SubscriptionID              int            `json:"subscriptionId,omitempty"`
	PackageID                   int            `json:"packageId,omitempty"`
	PackageName                 string         `json:"packageName,omitempty"`
	PackageInternalName         string         `json:"packageInternalName,omitempty"`
	Users                       int            `json:"users,omitempty"`
	StartsAt                    string         `json:"startsAt,omitempty"`
	EndsAt                      string         `json:"endsAt,omitempty"`
	Strength                    *TrialStrength `json:"strength,omitempty"`
	AdminEmails                 []string       `json:"adminEmails,omitempty"`
	StreetAddress               string         `json:"streetAddress,omitempty"`
	StreetAddress2              string         `json:"streetAddress2,omitempty"`
	City                        string         `json:"city,omitempty"`
	State                       string         `json:"state,omitempty"`
	CountryCode                 string         `json:"countryCode,omitempty"`
	ZipCode                     string         `json:"zipCode,omitempty"`
	DealID                      string         `json:"dealId,omitempty"`
	PpovLifecycle               *PpovLifecycle `json:"ppovLifecycle,omitempty"`
	HasDistributorVisibility    bool           `json:"hasDistributorVisibility,omitempty"`
	IsOnboardingWizardCompleted bool           `json:"isOnboardingWizardCompleted,omitempty"`
	TrialID                     string         `json:"trialId,omitempty"`
	AccountManagerEmails        []string       `json:"accountManagerEmails,omitempty"`
	AccessRequestID             int            `json:"accessRequestId,omitempty"`
	AccessRequestState          string         `json:"accessRequestState,omitempty"`
	TrialPeriod                 string         `json:"trialPeriod,omitempty"`
	TrialExtensionCount         int            `json:"trialExtensionCount,omitempty"`
	TrialExtendedDays           int            `json:"trialExtendedDays,omitempty"`
}

// PpovLifecycle holds the email details about lifecycle events from a trial.
type PpovLifecycle struct {
	State                   string           `json:"state,omitempty"`
	Date                    string           `json:"date,omitempty"`
	Enabled                 bool             `json:"enabled,omitempty"`
	LastSentDate            string           `json:"lastSentDate,omitempty"`
	MailIdentifiers         *MailIdentifiers `json:"mailIdentifiers,omitempty"`
	ExcludedLifecycleEmails []string         `json:"excludedLifecycleEmails,omitempty"`
}

// MailIdentifiers holds the lifecycle email identifiers for a customer.
type MailIdentifiers struct {
	NoLoginDayFourMailIdentifier   string `json:"noLoginDayFourMailIdentifier,omitempty"`
	NoOriginDayThreeMailIdentifier string `json:"noOriginDayThreeMailIdentifier,omitempty"`
	NoOriginDaySevenMailIdentifier string `json:"noOriginDaySevenMailIdentifier,omitempty"`
}

// TrialStrength is the number of features consumed by a customer trial.
type TrialStrength struct {
	CustomerLoggedIn  bool   `json:"customerLoggedIn"`
	LastLoginDate     string `json:"lastLoginDate"`
	IdentitiesCreated int    `json:"identitiesCreated"`
	HasTraffic        bool   `json:"hasTraffic"`
	TrialStrength     string `json:"trialStrength"`
}

// PutTrialConversion converts an MSLA trial to an MSLA customer using the given
// package ID and returns the conversion status. Valid package IDs are the
// trial-conversion subset (246, 248, 250, 252).
func (c *Client) PutTrialConversion(ctx context.Context, customerID, packageID int) (string, error) {
	url := endpointURL(fmt.Sprintf("providers/customers/%d/trialconversions", customerID))
	body := map[string]int{"packageId": packageID}
	result, err := c.Put[struct {
		ConversionStatus string `json:"conversionStatus"`
	}](ctx, url, body)
	if err != nil {
		return "", fmt.Errorf("put trial conversion: %w", err)
	}
	return result.ConversionStatus, nil
}

// CreateTrialExtension extends the customer's trial by the given number of days
// (7 or 14) and returns the updated subscription.
func (c *Client) CreateTrialExtension(ctx context.Context, customerID, days int) (*CustomerSubscription, error) {
	url := endpointURL(fmt.Sprintf("providers/customers/%d/trialExtensions", customerID))
	body := map[string]int{"trialExtensionDays": days}
	result, err := c.Post[CustomerSubscription](ctx, url, body)
	if err != nil {
		return nil, fmt.Errorf("create trial extension: %w", err)
	}
	return result, nil
}

// GetSubscriptionDetails gets the subscription details for the customer's organization.
func (c *Client) GetSubscriptionDetails(ctx context.Context, customerID int) (*CustomerSubscription, error) {
	url := endpointURL(fmt.Sprintf("providers/customers/%d/subscriptionDetails", customerID))
	result, err := c.Get[CustomerSubscription](ctx, url, nil)
	if err != nil {
		return nil, fmt.Errorf("get subscription details: %w", err)
	}
	return result, nil
}

// GetTrialStrength gets the strength of a customer trial.
func (c *Client) GetTrialStrength(ctx context.Context, customerID int) (*TrialStrength, error) {
	url := endpointURL(fmt.Sprintf("providers/customers/%d/trialStrengths", customerID))
	result, err := c.Get[TrialStrength](ctx, url, nil)
	if err != nil {
		return nil, fmt.Errorf("get trial strength: %w", err)
	}
	return result, nil
}
