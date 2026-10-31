package psa

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/thecoretg/tctg-go/internal/httpx"
)

const (
	defaultBaseURL = "https://api-na.myconnectwise.net/v4_6_release/apis/3.0"
)

var (
	ErrNotFound = errors.New("404 status returned")
	// ErrForbidden is a 403. ConnectWise also answers 403 for a record the API can no longer reach
	// although it still exists, such as a ticket merged into another, so it is not always a
	// permissions problem.
	ErrForbidden = errors.New("403 status returned")
)

// checkStatus turns an error response into an error, mapping 404 to ErrNotFound, 403 to
// ErrForbidden and 429 to a *ratelimit.Error so callers can match them with errors.Is whatever
// the method. A 403 keeps the
// response body, which says why ConnectWise refused.
func checkStatus(res *httpx.Response) error {
	switch {
	case res.StatusCode < 400:
		return nil
	case res.StatusCode == http.StatusNotFound:
		return ErrNotFound
	case res.StatusCode == http.StatusTooManyRequests:
		return httpx.RateLimitError("connectwise psa", res)
	case res.StatusCode == http.StatusForbidden:
		return fmt.Errorf("error response from ConnectWise API: %w: %s", ErrForbidden, res.Body)
	default:
		return fmt.Errorf("error response from ConnectWise API: %s", res.Body)
	}
}

// Get issues a GET request and decodes the JSON response into T.
func (c *Client) Get[T any](ctx context.Context, endpoint string, params map[string]string) (*T, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodGet, c.endpointURL(endpoint), params, nil)
	if err != nil {
		return nil, err
	}

	if err := checkStatus(res); err != nil {
		return nil, err
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

// GetMany follows ConnectWise's Link header until every page is collected, or
// until a WithLimit cap is reached.
func (c *Client) GetMany[T any](ctx context.Context, endpoint string, params map[string]string, opts ...ListOption) ([]T, error) {
	cfg := newListConfig(opts)
	if cfg.limit > 0 {
		params = withPageSize(params, min(cfg.limit, maxPageSize))
	}

	var allItems []T

	next := c.endpointURL(endpoint)
	for next != "" {
		res, err := httpx.Do(ctx, c.httpClient, http.MethodGet, next, params, nil)
		if err != nil {
			return nil, err
		}

		if err := checkStatus(res); err != nil {
			return nil, err
		}

		var target []T
		if err := httpx.DecodeJSON(res.Body, &target); err != nil {
			return nil, err
		}

		allItems = append(allItems, target...)
		if cfg.limit > 0 && len(allItems) >= cfg.limit {
			return allItems[:cfg.limit], nil
		}

		params = nil
		next = parseLinkHeader(res.Header.Get("Link"), "next")
	}

	return allItems, nil
}

// Post issues a POST request with body and decodes the JSON response into T.
func (c *Client) Post[T any](ctx context.Context, endpoint string, body any) (*T, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodPost, c.endpointURL(endpoint), nil, body)
	if err != nil {
		return nil, err
	}

	if err := checkStatus(res); err != nil {
		return nil, err
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

// Put issues a PUT request with body and decodes the JSON response into T.
func (c *Client) Put[T any](ctx context.Context, endpoint string, body any) (*T, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodPut, c.endpointURL(endpoint), nil, body)
	if err != nil {
		return nil, err
	}

	if err := checkStatus(res); err != nil {
		return nil, err
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

// Patch issues a PATCH request and decodes the JSON response into T.
func (c *Client) Patch[T any](ctx context.Context, endpoint string, patchOps []PatchOp) (*T, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodPatch, c.endpointURL(endpoint), nil, patchOps)
	if err != nil {
		return nil, err
	}

	if err := checkStatus(res); err != nil {
		return nil, err
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

// Delete issues a DELETE request, discarding any response body.
func (c *Client) Delete(ctx context.Context, endpoint string) error {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodDelete, c.endpointURL(endpoint), nil, nil)
	if err != nil {
		return err
	}

	return checkStatus(res)
}

func (c *Client) endpointURL(endpoint string) string {
	return fmt.Sprintf("%s/%s", c.baseURL, endpoint)
}

func parseLinkHeader(linkHeader, rel string) string {
	for link := range strings.SplitSeq(linkHeader, ",") {
		parts := strings.Split(strings.TrimSpace(link), ";")
		if len(parts) < 2 {
			continue
		}
		urlPart := strings.Trim(parts[0], "<>")
		relPart := strings.TrimSpace(parts[1])
		if relPart == fmt.Sprintf(`rel="%s"`, rel) {
			return urlPart
		}
	}

	return ""
}
