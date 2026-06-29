package addigy

import "context"

// collectPaged accumulates every page of a paginated Addigy endpoint. fetch is
// called once per page (1-based) and returns that page's items and the response
// Metadata. Paging stops at Metadata.PageCount, falling back to an empty page
// when the count is unknown (zero).
func collectPaged[T any](startPage int, fetch func(page int) ([]T, Metadata, error)) ([]T, error) {
	page := max(startPage, 1)
	var all []T
	for {
		items, meta, err := fetch(page)
		if err != nil {
			return nil, err
		}
		all = append(all, items...)
		if len(items) == 0 || (meta.PageCount > 0 && page >= meta.PageCount) {
			break
		}
		page++
	}
	return all, nil
}

// SearchAllDevices runs SearchDevices across every page and returns all matching
// devices. The filter's Page field sets the starting page (default 1); PerPage
// controls the page size.
func (c *Client) SearchAllDevices(ctx context.Context, filter DeviceFilter) ([]DeviceAudit, error) {
	return collectPaged(filter.Page, func(page int) ([]DeviceAudit, Metadata, error) {
		filter.Page = page
		resp, err := c.SearchDevices(ctx, filter)
		if err != nil {
			return nil, Metadata{}, err
		}
		return resp.Items, resp.Metadata, nil
	})
}

// QueryAllCustomFacts runs QueryCustomFacts across every page and returns all
// matching custom facts. The request's Page field sets the starting page.
func (c *Client) QueryAllCustomFacts(ctx context.Context, req FactQuery) ([]Fact, error) {
	return collectPaged(req.Page, func(page int) ([]Fact, Metadata, error) {
		req.Page = page
		resp, err := c.QueryCustomFacts(ctx, req)
		if err != nil {
			return nil, Metadata{}, err
		}
		return resp.Items, resp.Metadata, nil
	})
}

// QueryAllVariables runs QueryVariables across every page and returns all
// matching variables. The request's Page field sets the starting page.
func (c *Client) QueryAllVariables(ctx context.Context, req VariablesQueryRequest) ([]Variable, error) {
	return collectPaged(req.Page, func(page int) ([]Variable, Metadata, error) {
		req.Page = page
		resp, err := c.QueryVariables(ctx, req)
		if err != nil {
			return nil, Metadata{}, err
		}
		return resp.Items, resp.Metadata, nil
	})
}
