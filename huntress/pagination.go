package huntress

import (
	"context"
	"maps"
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

// ListAllOrganizations fetches every page of organizations.
func (c *Client) ListAllOrganizations(ctx context.Context, params map[string]string) ([]Organization, error) {
	return collectPaged(func(token string) ([]Organization, Pagination, error) {
		resp, err := c.ListOrganizations(ctx, pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.Organizations, resp.Pagination, nil
	})
}

// ListAllMemberships fetches every page of memberships.
func (c *Client) ListAllMemberships(ctx context.Context, params map[string]string) ([]Membership, error) {
	return collectPaged(func(token string) ([]Membership, Pagination, error) {
		resp, err := c.ListMemberships(ctx, pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.Memberships, resp.Pagination, nil
	})
}

// ListAllAgents fetches every page of agents.
func (c *Client) ListAllAgents(ctx context.Context, params map[string]string) ([]Agent, error) {
	return collectPaged(func(token string) ([]Agent, Pagination, error) {
		resp, err := c.ListAgents(ctx, pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.Agents, resp.Pagination, nil
	})
}

// ListAllEscalations fetches every page of escalations.
func (c *Client) ListAllEscalations(ctx context.Context, params map[string]string) ([]Escalation, error) {
	return collectPaged(func(token string) ([]Escalation, Pagination, error) {
		resp, err := c.ListEscalations(ctx, pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.Escalations, resp.Pagination, nil
	})
}

// ListAllIdentities fetches every page of identities.
func (c *Client) ListAllIdentities(ctx context.Context, params map[string]string) ([]Identity, error) {
	return collectPaged(func(token string) ([]Identity, Pagination, error) {
		resp, err := c.ListIdentities(ctx, pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.Identities, resp.Pagination, nil
	})
}

// ListAllAccounts fetches every page of accounts (Reseller credentials only).
func (c *Client) ListAllAccounts(ctx context.Context, params map[string]string) ([]Account, error) {
	return collectPaged(func(token string) ([]Account, Pagination, error) {
		resp, err := c.ListAccounts(ctx, pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.Accounts, resp.Pagination, nil
	})
}

// ListAllIncidentReports fetches every page of incident reports.
func (c *Client) ListAllIncidentReports(ctx context.Context, params map[string]string) ([]IncidentReport, error) {
	return collectPaged(func(token string) ([]IncidentReport, Pagination, error) {
		resp, err := c.ListIncidentReports(ctx, pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.IncidentReports, resp.Pagination, nil
	})
}

// ListAllExternalPorts fetches every page of external ports.
func (c *Client) ListAllExternalPorts(ctx context.Context, params map[string]string) ([]ExternalPort, error) {
	return collectPaged(func(token string) ([]ExternalPort, Pagination, error) {
		resp, err := c.ListExternalPorts(ctx, pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.ExternalPorts, resp.Pagination, nil
	})
}

// ListAllInvoices fetches every page of invoices.
func (c *Client) ListAllInvoices(ctx context.Context, params map[string]string) ([]Invoice, error) {
	return collectPaged(func(token string) ([]Invoice, Pagination, error) {
		resp, err := c.ListInvoices(ctx, pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.Invoices, resp.Pagination, nil
	})
}

// ListAllReports fetches every page of summary reports.
func (c *Client) ListAllReports(ctx context.Context, params map[string]string) ([]SummaryReport, error) {
	return collectPaged(func(token string) ([]SummaryReport, Pagination, error) {
		resp, err := c.ListReports(ctx, pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.Reports, resp.Pagination, nil
	})
}

// ListAllSignals fetches every page of signals.
func (c *Client) ListAllSignals(ctx context.Context, params map[string]string) ([]Signal, error) {
	return collectPaged(func(token string) ([]Signal, Pagination, error) {
		resp, err := c.ListSignals(ctx, pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.Signals, resp.Pagination, nil
	})
}

// ListAllSubscriptions fetches every page of reseller subscriptions (Reseller
// credentials only).
func (c *Client) ListAllSubscriptions(ctx context.Context, params map[string]string) ([]Subscription, error) {
	return collectPaged(func(token string) ([]Subscription, Pagination, error) {
		resp, err := c.ListSubscriptions(ctx, pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.Subscriptions, resp.Pagination, nil
	})
}

// ListAllResellerInvoices fetches every page of reseller invoices (Reseller
// credentials only).
func (c *Client) ListAllResellerInvoices(ctx context.Context, params map[string]string) ([]Invoice, error) {
	return collectPaged(func(token string) ([]Invoice, Pagination, error) {
		resp, err := c.ListResellerInvoices(ctx, pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.Invoices, resp.Pagination, nil
	})
}

// ListAllInvoiceAccountUsageLineItems fetches every page of account usage line
// items for a reseller invoice (Reseller credentials only).
func (c *Client) ListAllInvoiceAccountUsageLineItems(ctx context.Context, invoiceID int64, params map[string]string) ([]AccountUsageLineItem, error) {
	return collectPaged(func(token string) ([]AccountUsageLineItem, Pagination, error) {
		resp, err := c.ListInvoiceAccountUsageLineItems(ctx, invoiceID, pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.AccountUsageLineItems, resp.Pagination, nil
	})
}

// ListAllInvoiceOrganizationUsageLineItems fetches every page of organization
// usage line items for a reseller invoice (Reseller credentials only).
func (c *Client) ListAllInvoiceOrganizationUsageLineItems(ctx context.Context, invoiceID int64, params map[string]string) ([]OrganizationUsageLineItem, error) {
	return collectPaged(func(token string) ([]OrganizationUsageLineItem, Pagination, error) {
		resp, err := c.ListInvoiceOrganizationUsageLineItems(ctx, invoiceID, pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.OrganizationUsageLineItems, resp.Pagination, nil
	})
}

// ListAllAccountOrganizations fetches every page of an account's organizations
// (Reseller credentials only).
func (c *Client) ListAllAccountOrganizations(ctx context.Context, accountID int64, params map[string]string) ([]Organization, error) {
	return collectPaged(func(token string) ([]Organization, Pagination, error) {
		resp, err := c.ListAccountOrganizations(ctx, accountID, pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.Organizations, resp.Pagination, nil
	})
}

// ListAllAccountMemberships fetches every page of an account's memberships
// (Reseller credentials only).
func (c *Client) ListAllAccountMemberships(ctx context.Context, accountID int64, params map[string]string) ([]Membership, error) {
	return collectPaged(func(token string) ([]Membership, Pagination, error) {
		resp, err := c.ListAccountMemberships(ctx, accountID, pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.Memberships, resp.Pagination, nil
	})
}

// ListAllAccountAgents fetches every page of an account's agents (Reseller
// credentials only).
func (c *Client) ListAllAccountAgents(ctx context.Context, accountID int64, params map[string]string) ([]Agent, error) {
	return collectPaged(func(token string) ([]Agent, Pagination, error) {
		resp, err := c.ListAccountAgents(ctx, accountID, pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.Agents, resp.Pagination, nil
	})
}

// ListAllAccountIncidentReports fetches every page of an account's incident
// reports (Reseller credentials only).
func (c *Client) ListAllAccountIncidentReports(ctx context.Context, accountID int64, params map[string]string) ([]IncidentReport, error) {
	return collectPaged(func(token string) ([]IncidentReport, Pagination, error) {
		resp, err := c.ListAccountIncidentReports(ctx, accountID, pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.IncidentReports, resp.Pagination, nil
	})
}

// ListAllAccountExternalPorts fetches every page of an account's external ports
// (Reseller credentials only).
func (c *Client) ListAllAccountExternalPorts(ctx context.Context, accountID int64, params map[string]string) ([]ExternalPort, error) {
	return collectPaged(func(token string) ([]ExternalPort, Pagination, error) {
		resp, err := c.ListAccountExternalPorts(ctx, accountID, pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.ExternalPorts, resp.Pagination, nil
	})
}

// ListAllAccountInvoices fetches every page of an account's invoices (Reseller
// credentials only).
func (c *Client) ListAllAccountInvoices(ctx context.Context, accountID int64, params map[string]string) ([]Invoice, error) {
	return collectPaged(func(token string) ([]Invoice, Pagination, error) {
		resp, err := c.ListAccountInvoices(ctx, accountID, pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.Invoices, resp.Pagination, nil
	})
}

// ListAllAccountReports fetches every page of an account's summary reports
// (Reseller credentials only).
func (c *Client) ListAllAccountReports(ctx context.Context, accountID int64, params map[string]string) ([]SummaryReport, error) {
	return collectPaged(func(token string) ([]SummaryReport, Pagination, error) {
		resp, err := c.ListAccountReports(ctx, accountID, pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.Reports, resp.Pagination, nil
	})
}

// ListAllAccountSignals fetches every page of an account's signals (Reseller
// credentials only).
func (c *Client) ListAllAccountSignals(ctx context.Context, accountID int64, params map[string]string) ([]Signal, error) {
	return collectPaged(func(token string) ([]Signal, Pagination, error) {
		resp, err := c.ListAccountSignals(ctx, accountID, pageParams(params, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.Signals, resp.Pagination, nil
	})
}
