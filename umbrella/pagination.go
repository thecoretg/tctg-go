package umbrella

import (
	"context"
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
