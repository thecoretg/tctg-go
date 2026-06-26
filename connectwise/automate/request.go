package automate

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/thecoretg/tctg-go/internal/httpx"
)

var ErrNotFound = errors.New("404 status returned")

func get[T any](ctx context.Context, c *Client, endpoint string, params map[string]string) (*T, error) {
	target, _, err := getWithHeaders[T](ctx, c, endpoint, params)
	return target, err
}

// getWithHeaders is get that also returns the response headers, used by the
// ListAll* pagination helpers to read Automate's Total-Count header.
func getWithHeaders[T any](ctx context.Context, c *Client, endpoint string, params map[string]string) (*T, http.Header, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodGet, c.fullURL(endpoint), params, nil)
	if err != nil {
		return nil, nil, err
	}

	if res.StatusCode >= 400 {
		if res.StatusCode == http.StatusNotFound {
			return nil, nil, ErrNotFound
		}
		return nil, nil, fmt.Errorf("error response from ConnectWise Automate API: %s [%d]", res.Body, res.StatusCode)
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, nil, err
	}
	return &target, res.Header, nil
}

func post[T any](ctx context.Context, c *Client, endpoint string, body any) (*T, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodPost, c.fullURL(endpoint), nil, body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("error response from ConnectWise Automate API: %s [%d]", res.Body, res.StatusCode)
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

func put[T any](ctx context.Context, c *Client, endpoint string, body any) (*T, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodPut, c.fullURL(endpoint), nil, body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		if res.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error response from ConnectWise Automate API: %s [%d]", res.Body, res.StatusCode)
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

func patch[T any](ctx context.Context, c *Client, endpoint string, patchOps []PatchOp) (*T, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodPatch, c.fullURL(endpoint), nil, patchOps)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("error response from ConnectWise Automate API: %s [%d]", res.Body, res.StatusCode)
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

func del(ctx context.Context, c *Client, endpoint string) error {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodDelete, c.fullURL(endpoint), nil, nil)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		if res.StatusCode == http.StatusNotFound {
			return ErrNotFound
		}
		return fmt.Errorf("error response from ConnectWise Automate API: %s [%d]", res.Body, res.StatusCode)
	}

	return nil
}

// delReturn issues a DELETE for endpoints that return a body (e.g. resetting an
// extra field returns the reset field).
func delReturn[T any](ctx context.Context, c *Client, endpoint string) (*T, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodDelete, c.fullURL(endpoint), nil, nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		if res.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error response from ConnectWise Automate API: %s [%d]", res.Body, res.StatusCode)
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}
