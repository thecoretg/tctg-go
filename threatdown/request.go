package threatdown

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var ErrNotFound = errors.New("404 status returned")

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
		return nil, fmt.Errorf("error response from Threatdown API: %s", res.String())
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
		return nil, fmt.Errorf("error response from Threatdown API: %s", res.String())
	}

	return &target, nil
}

func put[T any](ctx context.Context, c *Client, url string, body any) (*T, error) {
	var target T
	res, err := c.restClient.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&target).
		Put(url)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		if res.StatusCode() == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error response from Threatdown API: %s", res.String())
	}

	return &target, nil
}

// getAll fetches all pages of a cursor-paginated endpoint. The extract function
// returns the items and the next cursor from each response. If the endpoint does
// not paginate, extract should return an empty string for the cursor and getAll
// will return after the single request.
func getAll[T, R any](ctx context.Context, c *Client, url string, params map[string]string, extract func(R) ([]T, string)) ([]T, error) {
	var all []T
	cursor := ""
	for {
		p := make(map[string]string, len(params)+1)
		for k, v := range params {
			p[k] = v
		}
		if cursor != "" {
			p["cursor"] = cursor
		}

		result, err := get[R](ctx, c, url, p)
		if err != nil {
			return nil, err
		}

		items, nextCursor := extract(*result)
		all = append(all, items...)

		if nextCursor == "" {
			return all, nil
		}
		cursor = nextCursor
	}
}

func del(ctx context.Context, c *Client, url string) error {
	res, err := c.restClient.R().
		SetContext(ctx).
		Delete(url)
	if err != nil {
		return err
	}

	if res.IsError() {
		if res.StatusCode() == http.StatusNotFound {
			return ErrNotFound
		}
		return fmt.Errorf("error response from Threatdown API: %s", res.String())
	}

	return nil
}
