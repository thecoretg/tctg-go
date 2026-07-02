package umbrella

import (
	"context"
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

// ListAllCustomers fetches every page of ListCustomers and returns all customers.
func (c *Client) ListAllCustomers(ctx context.Context) ([]Customer, error) {
	return collectPaged(1, defaultPageLimit, func(page, limit int) ([]Customer, error) {
		return c.ListCustomers(ctx, map[string]string{
			"page":  strconv.Itoa(page),
			"limit": strconv.Itoa(limit),
		})
	})
}

// ListAllCustomerAddresses fetches every page of ListCustomerAddresses and
// returns all addresses.
func (c *Client) ListAllCustomerAddresses(ctx context.Context) ([]CustomerAddress, error) {
	return collectPaged(1, defaultPageLimit, func(page, limit int) ([]CustomerAddress, error) {
		return c.ListCustomerAddresses(ctx, map[string]string{
			"page":  strconv.Itoa(page),
			"limit": strconv.Itoa(limit),
		})
	})
}

// ListAllPolicies fetches every page of ListPolicies and returns all policies.
// policyType filters by "dns" or "web"; pass "" to use the API default (dns).
func (c *Client) ListAllPolicies(ctx context.Context, policyType string) ([]Policy, error) {
	return collectPaged(1, defaultPageLimit, func(page, limit int) ([]Policy, error) {
		params := map[string]string{
			"page":  strconv.Itoa(page),
			"limit": strconv.Itoa(limit),
		}
		if policyType != "" {
			params["type"] = policyType
		}
		return c.ListPolicies(ctx, params)
	})
}

// ListAllRoamingComputers fetches every page of ListRoamingComputers and
// returns all roaming computers. filters is merged into each page request and
// may contain any of the ListRoamingComputers filter params ("name", "status",
// "swgStatus", "lastSyncBefore", "lastSyncAfter"); pass nil for none.
func (c *Client) ListAllRoamingComputers(ctx context.Context, filters map[string]string) ([]RoamingComputer, error) {
	return collectPaged(1, defaultPageLimit, func(page, limit int) ([]RoamingComputer, error) {
		params := map[string]string{
			"page":  strconv.Itoa(page),
			"limit": strconv.Itoa(limit),
		}
		maps.Copy(params, filters)
		return c.ListRoamingComputers(ctx, params)
	})
}
