package iru

import (
	"context"
	"fmt"
	"net/http"

	"github.com/thecoretg/tctg-go/internal/httpx"
)

var ErrNotFound = fmt.Errorf("404 status returned")

// Get issues a GET request and decodes the JSON response into T.
func (c *Client) Get[T any](ctx context.Context, path string, params map[string]string) (*T, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodGet, c.baseURL+path, params, nil)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		if res.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error response from iru: %s", res.Body)
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

// Post issues a POST request with body and decodes the JSON response into T.
func (c *Client) Post[T any](ctx context.Context, path string, body any) (*T, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodPost, c.baseURL+path, nil, body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("error response from iru: %s", res.Body)
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

// Patch issues a PATCH request and decodes the JSON response into T.
func (c *Client) Patch[T any](ctx context.Context, path string, body any) (*T, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodPatch, c.baseURL+path, nil, body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("error response from iru: %s", res.Body)
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

// Delete issues a DELETE request, discarding any response body.
func (c *Client) Delete(ctx context.Context, path string) error {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodDelete, c.baseURL+path, nil, nil)
	if err != nil {
		return err
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("error response from iru: %s", res.Body)
	}
	return nil
}
