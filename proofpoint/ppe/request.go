package ppe

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/thecoretg/tctg-go/internal/httpx"
)

var (
	// ErrUnauthorized is returned for a 401: the username or password is
	// wrong. The API sends no further detail.
	ErrUnauthorized = errors.New("401 unauthorized")
	// ErrForbidden is returned for a 403, usually an organization outside the
	// account's hierarchy or a role without API access.
	ErrForbidden = errors.New("403 forbidden")
	// ErrNotFound is returned for a 404.
	ErrNotFound = errors.New("404 not found")
	// ErrConflict is returned for a 409: the resource being created already
	// exists.
	ErrConflict = errors.New("409 conflict")
)

// apiError maps an error response to an error. The sentinel statuses wrap
// their sentinel so errors.Is works, keeping any body the API sent as detail.
// A 422 or any other status wraps the body, which is where the API explains
// the problem.
func apiError(res *httpx.Response) error {
	var sentinel error
	switch res.StatusCode {
	case http.StatusTooManyRequests:
		return httpx.RateLimitError("proofpoint essentials", res)
	case http.StatusUnauthorized:
		sentinel = ErrUnauthorized
	case http.StatusForbidden:
		sentinel = ErrForbidden
	case http.StatusNotFound:
		sentinel = ErrNotFound
	case http.StatusConflict:
		sentinel = ErrConflict
	}

	body := strings.TrimSpace(string(res.Body))
	if sentinel != nil {
		if body == "" {
			return sentinel
		}
		return fmt.Errorf("%w: %s", sentinel, body)
	}
	return fmt.Errorf("error response from Proofpoint Essentials API: %s [%d]", body, res.StatusCode)
}

// do sends a request and returns the response, or an error for any status of
// 400 or above.
func (c *Client) do(ctx context.Context, method, url string, params map[string]string, body any) (*httpx.Response, error) {
	res, err := httpx.Do(ctx, c.httpClient, method, url, params, body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return nil, apiError(res)
	}
	return res, nil
}

// Get issues a GET request and decodes the JSON response into T.
func (c *Client) Get[T any](ctx context.Context, url string, params map[string]string) (*T, error) {
	res, err := c.do(ctx, http.MethodGet, url, params, nil)
	if err != nil {
		return nil, err
	}
	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

// Post issues a POST request with body and decodes the JSON response into T.
func (c *Client) Post[T any](ctx context.Context, url string, body any) (*T, error) {
	res, err := c.do(ctx, http.MethodPost, url, nil, body)
	if err != nil {
		return nil, err
	}
	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

// Put issues a PUT request with body. The API answers a successful PUT with
// 204 No Content, so any response body is discarded.
func (c *Client) Put(ctx context.Context, url string, body any) error {
	_, err := c.do(ctx, http.MethodPut, url, nil, body)
	return err
}

// Patch issues a PATCH request with body. The API answers a successful PATCH
// with 204 No Content, so any response body is discarded.
func (c *Client) Patch(ctx context.Context, url string, body any) error {
	_, err := c.do(ctx, http.MethodPatch, url, nil, body)
	return err
}

// Delete issues a DELETE request, discarding any response body.
func (c *Client) Delete(ctx context.Context, url string) error {
	_, err := c.do(ctx, http.MethodDelete, url, nil, nil)
	return err
}

// getList fetches a list endpoint and decodes its items into T. The spec
// documents every list response as a bare array, but most are wrapped in an
// object keyed by resource type ({"users": [...]}); billing alone returns a
// bare array. Both are accepted: a bare array, or an object holding the list
// under one of keys.
func (c *Client) getList[T any](ctx context.Context, url string, keys ...string) ([]T, error) {
	raw, err := c.Get[json.RawMessage](ctx, url, nil)
	if err != nil {
		return nil, err
	}
	items, err := decodeList[T](*raw, keys...)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", url, err)
	}
	return items, nil
}

func decodeList[T any](data json.RawMessage, keys ...string) ([]T, error) {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return nil, nil
	}

	if data[0] == '[' {
		var items []T
		if err := json.Unmarshal(data, &items); err != nil {
			return nil, err
		}
		return items, nil
	}

	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}
	for _, k := range keys {
		if v, ok := obj[k]; ok {
			var items []T
			if err := httpx.DecodeJSON(v, &items); err != nil {
				return nil, fmt.Errorf("key %q: %w", k, err)
			}
			return items, nil
		}
	}

	present := make([]string, 0, len(obj))
	for k := range obj {
		present = append(present, k)
	}
	slices.Sort(present)
	return nil, fmt.Errorf("list response has none of keys %v; got keys %v", keys, present)
}
