package rewst

import (
	"context"
	"fmt"
	"net/http"
)

var ErrNotFound = fmt.Errorf("404 status returned")

func get[T any](ctx context.Context, c *Client, url string, params map[string]string) (*T, error) {
	var target T
	res, err := c.restClient.R().
		SetContext(ctx).
		SetQueryParams(params).
		SetResult(&target).
		Get(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		if res.StatusCode() == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error response from Rewst: %s", res.String())
	}

	return &target, nil
}

func post[T any](ctx context.Context, c *Client, url string, body any) (*T, error) {
	var target T
	res, err := c.restClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&target).
		Post(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("error response from Rewst: %s", res.String())
	}

	return &target, nil
}
