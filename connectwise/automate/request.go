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
func send(ctx context.Context, c *Client, method, endpoint string, params map[string]string, body any) (*httpx.Response, error) {
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

func get[T any](ctx context.Context, c *Client, endpoint string, params map[string]string) (*T, error) {
	target, _, err := getWithHeaders[T](ctx, c, endpoint, params)
	return target, err
}

// getWithHeaders is get that also returns the response headers, used by the
// ListAll* pagination helpers to read Automate's Total-Count header.
func getWithHeaders[T any](ctx context.Context, c *Client, endpoint string, params map[string]string) (*T, http.Header, error) {
	res, err := send(ctx, c, http.MethodGet, endpoint, params, nil)
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

func post[T any](ctx context.Context, c *Client, endpoint string, body any) (*T, error) {
	res, err := send(ctx, c, http.MethodPost, endpoint, nil, body)
	if err != nil {
		return nil, err
	}
	if err := checkStatus(res); err != nil {
		return nil, err
	}
	return decode[T](res)
}

func put[T any](ctx context.Context, c *Client, endpoint string, body any) (*T, error) {
	res, err := send(ctx, c, http.MethodPut, endpoint, nil, body)
	if err != nil {
		return nil, err
	}
	if err := checkStatus(res); err != nil {
		return nil, err
	}
	return decode[T](res)
}

func patch[T any](ctx context.Context, c *Client, endpoint string, patchOps []PatchOp) (*T, error) {
	res, err := send(ctx, c, http.MethodPatch, endpoint, nil, patchOps)
	if err != nil {
		return nil, err
	}
	if err := checkStatus(res); err != nil {
		return nil, err
	}
	return decode[T](res)
}

func del(ctx context.Context, c *Client, endpoint string) error {
	res, err := send(ctx, c, http.MethodDelete, endpoint, nil, nil)
	if err != nil {
		return err
	}
	return checkStatus(res)
}

// delReturn issues a DELETE for endpoints that return a body (e.g. resetting an
// extra field returns the reset field).
func delReturn[T any](ctx context.Context, c *Client, endpoint string) (*T, error) {
	res, err := send(ctx, c, http.MethodDelete, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}
	if err := checkStatus(res); err != nil {
		return nil, err
	}
	return decode[T](res)
}
