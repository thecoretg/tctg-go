package iru

import (
	"context"
	"fmt"
	"net/http"
)

var ErrNotFound = fmt.Errorf("404 status returned")

func Get[T any](ctx context.Context, c *Client, path string, params map[string]string) (*T, error) {
	var target T
	res, err := c.restyClient.R().
		SetContext(ctx).
		SetQueryParams(params).
		SetResult(&target).
		Get(path)
	if err != nil {
		return nil, err
	}
	if res.IsStatusFailure() {
		if res.StatusCode() == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error response from iru: %s", res.String())
	}
	return &target, nil
}

func Post[T any](ctx context.Context, c *Client, path string, body any) (*T, error) {
	var target T
	res, err := c.restyClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&target).
		Post(path)
	if err != nil {
		return nil, err
	}
	if res.IsStatusFailure() {
		return nil, fmt.Errorf("error response from iru: %s", res.String())
	}
	return &target, nil
}

func Patch[T any](ctx context.Context, c *Client, path string, body any) (*T, error) {
	var target T
	res, err := c.restyClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&target).
		Patch(path)
	if err != nil {
		return nil, err
	}
	if res.IsStatusFailure() {
		return nil, fmt.Errorf("error response from iru: %s", res.String())
	}
	return &target, nil
}

func Delete(ctx context.Context, c *Client, path string) error {
	res, err := c.restyClient.R().
		SetContext(ctx).
		Delete(path)
	if err != nil {
		return err
	}
	if res.IsStatusFailure() {
		return fmt.Errorf("error response from iru: %s", res.String())
	}
	return nil
}
