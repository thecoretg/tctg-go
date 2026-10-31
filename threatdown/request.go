package threatdown

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"strings"

	"github.com/thecoretg/tctg-go/internal/httpx"
)

var ErrNotFound = errors.New("404 status returned")

// apiError turns an error response into an error, mapping 429 to a
// *ratelimit.Error.
func apiError(res *httpx.Response) error {
	if res.StatusCode == http.StatusTooManyRequests {
		return httpx.RateLimitError("threatdown", res)
	}
	return fmt.Errorf("error response from Threatdown API: %s [%d]", res.Body, res.StatusCode)
}

// Get issues a GET request and decodes the JSON response into T.
func (c *Client) Get[T any](ctx context.Context, url string, params map[string]string) (*T, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodGet, url, params, nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		if res.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, apiError(res)
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

// Post issues a POST request with body and decodes the JSON response into T.
func (c *Client) Post[T any](ctx context.Context, url string, body any) (*T, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodPost, url, nil, body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, apiError(res)
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

// Put issues a PUT request with body and decodes the JSON response into T.
func (c *Client) Put[T any](ctx context.Context, url string, body any) (*T, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodPut, url, nil, body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		if res.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, apiError(res)
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

// ListAll fetches every page of a cursor-paginated ThreatDown GET endpoint and
// decodes its items into T. Pass json.RawMessage to keep each item exactly as
// the API returned it. path is relative to the OneView root and includes the
// API version, e.g. "v1/users". ThreatDown names a list's item array after the
// resource, so the items are read from the key matching path's last segment;
// a response without that key is an error. params are sent with every page,
// and the next_cursor of each response is sent back as the cursor param.
func (c *Client) ListAll[T any](ctx context.Context, path string, params map[string]string) ([]T, error) {
	all, err := collectCursor[T](path, func(cursor string) (map[string]json.RawMessage, error) {
		p := maps.Clone(params)
		if p == nil {
			p = make(map[string]string, 1)
		}
		if cursor != "" {
			p["cursor"] = cursor
		}
		page, err := c.Get[map[string]json.RawMessage](ctx, oneViewURL(path), p)
		if err != nil {
			return nil, err
		}
		return *page, nil
	})
	if err != nil {
		return nil, fmt.Errorf("list all %s: %w", path, err)
	}
	return all, nil
}

// SearchAll is ListAll for ThreatDown's POST search endpoints, such as
// "v1/accounts/{account_id}/endpoints": body is sent with every page, with the
// previous response's next_cursor added to it.
func (c *Client) SearchAll[T any](ctx context.Context, path string, body map[string]any) ([]T, error) {
	all, err := collectCursor[T](path, func(cursor string) (map[string]json.RawMessage, error) {
		b := maps.Clone(body)
		if b == nil {
			b = make(map[string]any, 1)
		}
		if cursor != "" {
			b["next_cursor"] = cursor
		}
		page, err := c.Post[map[string]json.RawMessage](ctx, oneViewURL(path), b)
		if err != nil {
			return nil, err
		}
		return *page, nil
	})
	if err != nil {
		return nil, fmt.Errorf("search all %s: %w", path, err)
	}
	return all, nil
}

// collectCursor requests pages from fetch until a response has no
// next_cursor, decoding each page's items from the key named by path's last
// segment.
func collectCursor[T any](path string, fetch func(cursor string) (map[string]json.RawMessage, error)) ([]T, error) {
	key := path[strings.LastIndex(path, "/")+1:]
	var all []T
	cursor := ""
	for {
		page, err := fetch(cursor)
		if err != nil {
			return nil, err
		}
		raw, ok := page[key]
		if !ok {
			return nil, fmt.Errorf("response has no %q array", key)
		}
		var items []T
		if err := json.Unmarshal(raw, &items); err != nil {
			return nil, fmt.Errorf("decoding %q: %w", key, err)
		}
		all = append(all, items...)

		cursor = ""
		if raw, ok := page["next_cursor"]; ok {
			if err := json.Unmarshal(raw, &cursor); err != nil {
				return nil, fmt.Errorf("decoding next_cursor: %w", err)
			}
		}
		if cursor == "" {
			return all, nil
		}
	}
}

// Delete issues a DELETE request, discarding any response body.
func (c *Client) Delete(ctx context.Context, url string) error {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodDelete, url, nil, nil)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		if res.StatusCode == http.StatusNotFound {
			return ErrNotFound
		}
		return apiError(res)
	}

	return nil
}
