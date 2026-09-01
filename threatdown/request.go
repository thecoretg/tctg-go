package threatdown

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"net/http"

	"github.com/thecoretg/tctg-go/internal/httpx"
)

var ErrNotFound = errors.New("404 status returned")

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
		return nil, fmt.Errorf("error response from Threatdown API: %s [%d]", res.Body, res.StatusCode)
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
		return nil, fmt.Errorf("error response from Threatdown API: %s [%d]", res.Body, res.StatusCode)
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
		return nil, fmt.Errorf("error response from Threatdown API: %s [%d]", res.Body, res.StatusCode)
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

// getAll fetches all pages of a cursor-paginated endpoint. The extract function
// returns the items and the next cursor from each response. If the endpoint does
// not paginate, extract should return an empty string for the cursor and getAll
// will return after the single request.
func (c *Client) getAll[T, R any](ctx context.Context, url string, params map[string]string, extract func(R) ([]T, string)) ([]T, error) {
	var all []T
	cursor := ""
	for {
		p := make(map[string]string, len(params)+1)
		maps.Copy(p, params)
		if cursor != "" {
			p["cursor"] = cursor
		}

		result, err := c.Get[R](ctx, url, p)
		if err != nil {
			return nil, err
		}

		items, nextCursor := extract(*result)
		all = append(all, items...)

		if nextCursor == "" {
			return all, nil
		}
		cursor = nextCursor
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
		return fmt.Errorf("error response from Threatdown API: %s [%d]", res.Body, res.StatusCode)
	}

	return nil
}
