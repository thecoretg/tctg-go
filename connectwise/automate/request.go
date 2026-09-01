package automate

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/thecoretg/tctg-go/internal/httpx"
)

var ErrNotFound = errors.New("404 status returned")

// send issues one request and, when the response is a 401 the client can
// recover from, logs in again and replays it exactly once. httpx.Do re-encodes
// the body on every call, so a replay needs no buffering of its own.
func (c *Client) send(ctx context.Context, method, endpoint string, params map[string]string, body any) (*httpx.Response, error) {
	stale := c.currentToken()

	res, err := httpx.Do(ctx, c.httpClient, method, c.fullURL(endpoint), params, body)
	if err != nil || res.StatusCode != http.StatusUnauthorized || !c.canRelogin(endpoint) {
		return res, err
	}

	if err := c.relogin(ctx, stale); err != nil {
		return nil, err
	}

	return httpx.Do(ctx, c.httpClient, method, c.fullURL(endpoint), params, body)
}

// checkStatus turns an error response into an error, mapping 404 to ErrNotFound
// so callers can match it with errors.Is regardless of the method used.
func checkStatus(res *httpx.Response) error {
	switch {
	case res.StatusCode < 400:
		return nil
	case res.StatusCode == http.StatusNotFound:
		return ErrNotFound
	default:
		return fmt.Errorf("error response from ConnectWise Automate API: %s [%d]", res.Body, res.StatusCode)
	}
}

func decode[T any](res *httpx.Response) (*T, error) {
	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

// Get issues a GET request and decodes the JSON response into T.
func (c *Client) Get[T any](ctx context.Context, endpoint string, params map[string]string) (*T, error) {
	target, _, err := c.getWithHeaders[T](ctx, endpoint, params)
	return target, err
}

// getWithHeaders is Get that also returns the response headers, used by the
// ListAll* pagination helpers to read Automate's Total-Count header.
func (c *Client) getWithHeaders[T any](ctx context.Context, endpoint string, params map[string]string) (*T, http.Header, error) {
	res, err := c.send(ctx, http.MethodGet, endpoint, params, nil)
	if err != nil {
		return nil, nil, err
	}
	if err := checkStatus(res); err != nil {
		return nil, nil, err
	}

	target, err := decode[T](res)
	if err != nil {
		return nil, nil, err
	}
	return target, res.Header, nil
}

// Post issues a POST request with body and decodes the JSON response into T.
func (c *Client) Post[T any](ctx context.Context, endpoint string, body any) (*T, error) {
	res, err := c.send(ctx, http.MethodPost, endpoint, nil, body)
	if err != nil {
		return nil, err
	}
	if err := checkStatus(res); err != nil {
		return nil, err
	}
	return decode[T](res)
}

// Put issues a PUT request with body and decodes the JSON response into T.
func (c *Client) Put[T any](ctx context.Context, endpoint string, body any) (*T, error) {
	res, err := c.send(ctx, http.MethodPut, endpoint, nil, body)
	if err != nil {
		return nil, err
	}
	if err := checkStatus(res); err != nil {
		return nil, err
	}
	return decode[T](res)
}

// Patch issues a PATCH request and decodes the JSON response into T.
func (c *Client) Patch[T any](ctx context.Context, endpoint string, patchOps []PatchOp) (*T, error) {
	res, err := c.send(ctx, http.MethodPatch, endpoint, nil, patchOps)
	if err != nil {
		return nil, err
	}
	if err := checkStatus(res); err != nil {
		return nil, err
	}
	return decode[T](res)
}

// Delete issues a DELETE request, discarding any response body.
func (c *Client) Delete(ctx context.Context, endpoint string) error {
	res, err := c.send(ctx, http.MethodDelete, endpoint, nil, nil)
	if err != nil {
		return err
	}
	return checkStatus(res)
}

// deleteReturn issues a DELETE for endpoints that return a body (e.g. resetting an
// extra field returns the reset field).
func (c *Client) deleteReturn[T any](ctx context.Context, endpoint string) (*T, error) {
	res, err := c.send(ctx, http.MethodDelete, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}
	if err := checkStatus(res); err != nil {
		return nil, err
	}
	return decode[T](res)
}
