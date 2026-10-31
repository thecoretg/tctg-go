package addigy

import (
	"context"
	"encoding/json"
	"fmt"
)

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

// ListAll runs a paginated Addigy query endpoint across every page and
// decodes its items into T. Pass json.RawMessage to keep each item exactly as
// the API returned it. path is relative to the API root, e.g.
// "oa/variables/query"; org-scoped endpoints sit under
// "o/{organization_id}/", e.g. "o/<org id>/devices". body is the query (a
// struct or map, sent as JSON) and is resent for every page with its "page"
// field advanced; a "page" in body sets the starting page and "per_page" the
// page size.
func (c *Client) ListAll[T any](ctx context.Context, path string, body any) ([]T, error) {
	b, err := toObject(body)
	if err != nil {
		return nil, fmt.Errorf("list all %s: encoding body: %w", path, err)
	}
	start, _ := b["page"].(float64)
	all, err := collectPaged(int(start), func(page int) ([]T, Metadata, error) {
		b["page"] = page
		resp, err := c.Post[struct {
			Items    []T      `json:"items"`
			Metadata Metadata `json:"metadata"`
		}](ctx, c.url(path), b)
		if err != nil {
			return nil, Metadata{}, err
		}
		return resp.Items, resp.Metadata, nil
	})
	if err != nil {
		return nil, fmt.Errorf("list all %s: %w", path, err)
	}
	return all, nil
}

// toObject turns a struct or map into a JSON object the pager can set "page"
// on. A nil body becomes an empty object.
func toObject(body any) (map[string]any, error) {
	obj := map[string]any{}
	if body == nil {
		return obj, nil
	}
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}
	return obj, nil
}

// SearchAllDevices runs SearchDevices across every page and returns all matching
// devices. The filter's Page field sets the starting page (default 1); PerPage
// controls the page size.
func (c *Client) SearchAllDevices(ctx context.Context, filter DeviceFilter) ([]DeviceAudit, error) {
	path, err := c.orgPath("devices")
	if err != nil {
		return nil, err
	}
	return c.ListAll[DeviceAudit](ctx, path, filter)
}

// QueryAllCustomFacts runs QueryCustomFacts across every page and returns all
// matching custom facts. The request's Page field sets the starting page.
func (c *Client) QueryAllCustomFacts(ctx context.Context, req FactQuery) ([]Fact, error) {
	path, err := c.orgPath("facts/custom/query")
	if err != nil {
		return nil, err
	}
	return c.ListAll[Fact](ctx, path, req)
}

// QueryAllVariables runs QueryVariables across every page and returns all
// matching variables. The request's Page field sets the starting page.
func (c *Client) QueryAllVariables(ctx context.Context, req VariablesQueryRequest) ([]Variable, error) {
	return c.ListAll[Variable](ctx, "oa/variables/query", req)
}
