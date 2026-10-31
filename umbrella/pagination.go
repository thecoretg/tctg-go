package umbrella

import (
	"context"
	"fmt"
	"maps"
	"strconv"
)

// defaultPageLimit is the maximum page size the API allows.
const defaultPageLimit = 100

// collectPaged accumulates every page of a page-numbered Umbrella endpoint.
// fetch is called once per page (1-based) with the page size and returns that
// page's items. Paging stops when a page returns fewer than limit items (the
// API returns bare arrays with no total/page-count metadata).
func collectPaged[T any](startPage, limit int, fetch func(page, limit int) ([]T, error)) ([]T, error) {
	page := max(startPage, 1)
	if limit <= 0 {
		limit = defaultPageLimit
	}
	var all []T
	for {
		items, err := fetch(page, limit)
		if err != nil {
			return nil, err
		}
		all = append(all, items...)
		if len(items) < limit {
			break
		}
		page++
	}
	return all, nil
}

// ListAll fetches every page of a page-numbered Umbrella list endpoint and
// decodes its items into T. Pass json.RawMessage to keep each item exactly as
// the API returned it. path is relative to the API host and includes the API's
// own prefix, e.g. "admin/v2/providers/customers" or
// "deployments/v2/roamingcomputers". params are sent with every page; a "page"
// sets the starting page and a "limit" the page size, which defaults to the
// 100 maximum. Umbrella returns bare arrays with no total, so paging stops at
// the first page shorter than the page size.
func (c *Client) ListAll[T any](ctx context.Context, path string, params map[string]string) ([]T, error) {
	start, _ := strconv.Atoi(params["page"])
	limit, _ := strconv.Atoi(params["limit"])
	all, err := collectPaged(start, limit, func(page, limit int) ([]T, error) {
		p := maps.Clone(params)
		if p == nil {
			p = make(map[string]string, 2)
		}
		p["page"] = strconv.Itoa(page)
		p["limit"] = strconv.Itoa(limit)
		items, err := c.Get[[]T](ctx, hostURL(path), p)
		if err != nil {
			return nil, err
		}
		return *items, nil
	})
	if err != nil {
		return nil, fmt.Errorf("list all %s: %w", path, err)
	}
	return all, nil
}

// ListAllCustomers fetches every page of ListCustomers and returns all customers.
func (c *Client) ListAllCustomers(ctx context.Context) ([]Customer, error) {
	return c.ListAll[Customer](ctx, "admin/v2/providers/customers", nil)
}

// ListAllCustomerAddresses fetches every page of ListCustomerAddresses and
// returns all addresses.
func (c *Client) ListAllCustomerAddresses(ctx context.Context) ([]CustomerAddress, error) {
	return c.ListAll[CustomerAddress](ctx, "admin/v2/providers/customerAddresses", nil)
}

// ListAllPolicies fetches every page of ListPolicies and returns all policies.
// policyType filters by "dns" or "web"; pass "" to use the API default (dns).
func (c *Client) ListAllPolicies(ctx context.Context, policyType string) ([]Policy, error) {
	var params map[string]string
	if policyType != "" {
		params = map[string]string{"type": policyType}
	}
	return c.ListAll[Policy](ctx, "deployments/v2/policies", params)
}

// ListAllRoamingComputers fetches every page of ListRoamingComputers and
// returns all roaming computers. filters is merged into each page request and
// may contain any of the ListRoamingComputers filter params ("name", "status",
// "swgStatus", "lastSyncBefore", "lastSyncAfter"); pass nil for none.
func (c *Client) ListAllRoamingComputers(ctx context.Context, filters map[string]string) ([]RoamingComputer, error) {
	return c.ListAll[RoamingComputer](ctx, "deployments/v2/roamingcomputers", filters)
}
