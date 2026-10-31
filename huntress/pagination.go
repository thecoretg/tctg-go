package huntress

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"strings"
)

// maxPageLimit is the largest page size the Huntress API allows.
const maxPageLimit = "500"

// collectPaged accumulates every page of a cursor-paginated Huntress endpoint.
// fetch is called once per page with the current page token ("" for the first
// page) and returns that page's items plus the pagination cursor. Paging stops
// when the returned NextPageToken is empty.
func collectPaged[T any](fetch func(pageToken string) ([]T, Pagination, error)) ([]T, error) {
	var all []T
	token := ""
	for {
		items, pg, err := fetch(token)
		if err != nil {
			return nil, err
		}
		all = append(all, items...)
		if pg.NextPageToken == "" {
			break
		}
		token = pg.NextPageToken
	}
	return all, nil
}

// pageParams clones base and layers on the page token and a default page size,
// leaving any caller-supplied filters (and an explicit limit) intact.
func pageParams(base map[string]string, token string) map[string]string {
	p := make(map[string]string, len(base)+2)
	maps.Copy(p, base)
	if _, ok := p["limit"]; !ok {
		p["limit"] = maxPageLimit
	}
	if token != "" {
		p["page_token"] = token
	} else {
		delete(p, "page_token")
	}
	return p
}

// ListAll fetches every page of a Huntress list endpoint and decodes its items
// into T. Pass json.RawMessage to keep each item exactly as the API returned
// it. path is relative to the API root, e.g. "organizations" or
// "accounts/12/agents". Huntress names a list's item array after the resource,
// so the items are read from the key matching path's last segment; a response
// without that key is an error. params are sent with every page, and limit
// defaults to the 500 maximum.
func (c *Client) ListAll[T any](ctx context.Context, path string, params map[string]string) ([]T, error) {
	key := path[strings.LastIndex(path, "/")+1:]
	all, err := collectPaged(func(token string) ([]T, Pagination, error) {
		page, err := c.Get[map[string]json.RawMessage](ctx, endpointURL(path), pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		raw, ok := (*page)[key]
		if !ok {
			return nil, Pagination{}, fmt.Errorf("response has no %q array", key)
		}
		var items []T
		if err := json.Unmarshal(raw, &items); err != nil {
			return nil, Pagination{}, fmt.Errorf("decoding %q: %w", key, err)
		}
		var pg Pagination
		if raw, ok := (*page)["pagination"]; ok {
			if err := json.Unmarshal(raw, &pg); err != nil {
				return nil, Pagination{}, fmt.Errorf("decoding pagination: %w", err)
			}
		}
		return items, pg, nil
	})
	if err != nil {
		return nil, fmt.Errorf("list all %s: %w", path, err)
	}
	return all, nil
}

// ListAllOrganizations fetches every page of organizations.
func (c *Client) ListAllOrganizations(ctx context.Context, params map[string]string) ([]Organization, error) {
	return c.ListAll[Organization](ctx, "organizations", params)
}

// ListAllMemberships fetches every page of memberships.
func (c *Client) ListAllMemberships(ctx context.Context, params map[string]string) ([]Membership, error) {
	return c.ListAll[Membership](ctx, "memberships", params)
}

// ListAllAgents fetches every page of agents.
func (c *Client) ListAllAgents(ctx context.Context, params map[string]string) ([]Agent, error) {
	return c.ListAll[Agent](ctx, "agents", params)
}

// ListAllEscalations fetches every page of escalations.
func (c *Client) ListAllEscalations(ctx context.Context, params map[string]string) ([]Escalation, error) {
	return c.ListAll[Escalation](ctx, "escalations", params)
}

// ListAllIdentities fetches every page of identities.
func (c *Client) ListAllIdentities(ctx context.Context, params map[string]string) ([]Identity, error) {
	return c.ListAll[Identity](ctx, "identities", params)
}

// ListAllAccounts fetches every page of accounts (Reseller credentials only).
func (c *Client) ListAllAccounts(ctx context.Context, params map[string]string) ([]Account, error) {
	return c.ListAll[Account](ctx, "accounts", params)
}

// ListAllIncidentReports fetches every page of incident reports.
func (c *Client) ListAllIncidentReports(ctx context.Context, params map[string]string) ([]IncidentReport, error) {
	return c.ListAll[IncidentReport](ctx, "incident_reports", params)
}

// ListAllExternalPorts fetches every page of external ports.
func (c *Client) ListAllExternalPorts(ctx context.Context, params map[string]string) ([]ExternalPort, error) {
	return c.ListAll[ExternalPort](ctx, "external_ports", params)
}

// ListAllInvoices fetches every page of invoices.
func (c *Client) ListAllInvoices(ctx context.Context, params map[string]string) ([]Invoice, error) {
	return c.ListAll[Invoice](ctx, "invoices", params)
}

// ListAllReports fetches every page of summary reports.
func (c *Client) ListAllReports(ctx context.Context, params map[string]string) ([]SummaryReport, error) {
	return c.ListAll[SummaryReport](ctx, "reports", params)
}

// ListAllSignals fetches every page of signals.
func (c *Client) ListAllSignals(ctx context.Context, params map[string]string) ([]Signal, error) {
	return c.ListAll[Signal](ctx, "signals", params)
}

// ListAllSubscriptions fetches every page of reseller subscriptions (Reseller
// credentials only).
func (c *Client) ListAllSubscriptions(ctx context.Context, params map[string]string) ([]Subscription, error) {
	return c.ListAll[Subscription](ctx, "reseller/subscriptions", params)
}

// ListAllResellerInvoices fetches every page of reseller invoices (Reseller
// credentials only).
func (c *Client) ListAllResellerInvoices(ctx context.Context, params map[string]string) ([]Invoice, error) {
	return c.ListAll[Invoice](ctx, "reseller/invoices", params)
}

// ListAllInvoiceAccountUsageLineItems fetches every page of account usage line
// items for a reseller invoice (Reseller credentials only).
func (c *Client) ListAllInvoiceAccountUsageLineItems(ctx context.Context, invoiceID int64, params map[string]string) ([]AccountUsageLineItem, error) {
	return c.ListAll[AccountUsageLineItem](ctx, fmt.Sprintf("reseller/invoices/%d/account_usage_line_items", invoiceID), params)
}

// ListAllInvoiceOrganizationUsageLineItems fetches every page of organization
// usage line items for a reseller invoice (Reseller credentials only).
func (c *Client) ListAllInvoiceOrganizationUsageLineItems(ctx context.Context, invoiceID int64, params map[string]string) ([]OrganizationUsageLineItem, error) {
	return c.ListAll[OrganizationUsageLineItem](ctx, fmt.Sprintf("reseller/invoices/%d/organization_usage_line_items", invoiceID), params)
}

// ListAllAccountOrganizations fetches every page of an account's organizations
// (Reseller credentials only).
func (c *Client) ListAllAccountOrganizations(ctx context.Context, accountID int64, params map[string]string) ([]Organization, error) {
	return c.ListAll[Organization](ctx, fmt.Sprintf("accounts/%d/organizations", accountID), params)
}

// ListAllAccountMemberships fetches every page of an account's memberships
// (Reseller credentials only).
func (c *Client) ListAllAccountMemberships(ctx context.Context, accountID int64, params map[string]string) ([]Membership, error) {
	return c.ListAll[Membership](ctx, fmt.Sprintf("accounts/%d/memberships", accountID), params)
}

// ListAllAccountAgents fetches every page of an account's agents (Reseller
// credentials only).
func (c *Client) ListAllAccountAgents(ctx context.Context, accountID int64, params map[string]string) ([]Agent, error) {
	return c.ListAll[Agent](ctx, fmt.Sprintf("accounts/%d/agents", accountID), params)
}

// ListAllAccountIncidentReports fetches every page of an account's incident
// reports (Reseller credentials only).
func (c *Client) ListAllAccountIncidentReports(ctx context.Context, accountID int64, params map[string]string) ([]IncidentReport, error) {
	return c.ListAll[IncidentReport](ctx, fmt.Sprintf("accounts/%d/incident_reports", accountID), params)
}

// ListAllAccountExternalPorts fetches every page of an account's external ports
// (Reseller credentials only).
func (c *Client) ListAllAccountExternalPorts(ctx context.Context, accountID int64, params map[string]string) ([]ExternalPort, error) {
	return c.ListAll[ExternalPort](ctx, fmt.Sprintf("accounts/%d/external_ports", accountID), params)
}

// ListAllAccountInvoices fetches every page of an account's invoices (Reseller
// credentials only).
func (c *Client) ListAllAccountInvoices(ctx context.Context, accountID int64, params map[string]string) ([]Invoice, error) {
	return c.ListAll[Invoice](ctx, fmt.Sprintf("accounts/%d/invoices", accountID), params)
}

// ListAllAccountReports fetches every page of an account's summary reports
// (Reseller credentials only).
func (c *Client) ListAllAccountReports(ctx context.Context, accountID int64, params map[string]string) ([]SummaryReport, error) {
	return c.ListAll[SummaryReport](ctx, fmt.Sprintf("accounts/%d/reports", accountID), params)
}

// ListAllAccountSignals fetches every page of an account's signals (Reseller
// credentials only).
func (c *Client) ListAllAccountSignals(ctx context.Context, accountID int64, params map[string]string) ([]Signal, error) {
	return c.ListAll[Signal](ctx, fmt.Sprintf("accounts/%d/signals", accountID), params)
}
