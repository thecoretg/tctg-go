package iru

import (
	"context"
	"fmt"
	"net/http"

	"github.com/thecoretg/tctg-go/internal/httpx"
)

var ErrNotFound = fmt.Errorf("404 status returned")

func Get[T any](ctx context.Context, c *Client, path string, params map[string]string) (*T, error) {
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

func Post[T any](ctx context.Context, c *Client, path string, body any) (*T, error) {
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

func Patch[T any](ctx context.Context, c *Client, path string, body any) (*T, error) {
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

func Delete(ctx context.Context, c *Client, path string) error {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodDelete, c.baseURL+path, nil, nil)
	if err != nil {
		return err
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("error response from iru: %s", res.Body)
	}
	return nil
}
