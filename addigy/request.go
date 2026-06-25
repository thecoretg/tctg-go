package addigy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/thecoretg/tctg-go/internal/httpx"
)

var ErrNotFound = errors.New("404 status returned")

// apiError formats an error response from the Addigy API, decoding the standard
// ErrorResponse body when present for a clearer message.
func apiError(res *httpx.Response) error {
	var er ErrorResponse
	if err := json.Unmarshal(res.Body, &er); err == nil && er.Message != "" {
		return fmt.Errorf("error response from Addigy API: %s [%d]", er.Message, res.StatusCode)
	}
	return fmt.Errorf("error response from Addigy API: %s [%d]", res.Body, res.StatusCode)
}

func get[T any](ctx context.Context, c *Client, url string, params map[string]string) (*T, error) {
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

func post[T any](ctx context.Context, c *Client, url string, body any) (*T, error) {
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

func put[T any](ctx context.Context, c *Client, url string, body any) (*T, error) {
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

func patch[T any](ctx context.Context, c *Client, url string, body any) (*T, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodPatch, url, nil, body)
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

// del issues a DELETE request. Several Addigy delete endpoints identify the
// target via query parameters rather than a path segment, so params is passed
// through to the request.
func del(ctx context.Context, c *Client, url string, params map[string]string) error {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodDelete, url, params, nil)
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
