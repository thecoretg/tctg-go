package psa

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	baseUrl = "https://api-na.myconnectwise.net/v4_6_release/apis/3.0"
)

var ErrNotFound = errors.New("404 status returned")

func GetOne[T any](ctx context.Context, c *Client, endpoint string, params map[string]string) (*T, error) {
	var target T
	res, err := c.restClient.R().
		SetContext(ctx).
		SetQueryParams(params).
		SetResult(&target).
		Get(fullURL(baseUrl, endpoint))
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		if res.StatusCode() == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error response from ConnectWise API: %s", res.String())
	}

	return &target, nil
}

func GetMany[T any](ctx context.Context, c *Client, endpoint string, params map[string]string) ([]T, error) {
	var allItems []T

	endpoint = fullURL(baseUrl, endpoint)
	for endpoint != "" {
		var target []T
		res, err := c.restClient.R().
			SetContext(ctx).
			SetQueryParams(params).
			SetResult(&target).
			Get(endpoint)
		if err != nil {
			return nil, err
		}

		if res.IsError() {
			if res.StatusCode() == http.StatusNotFound {
				return nil, ErrNotFound
			}
			return nil, fmt.Errorf("error response from ConnectWise API: %s", res.String())
		}

		allItems = append(allItems, target...)
		params = nil
		endpoint = parseLinkHeader(res.Header().Get("Link"), "next")
	}

	return allItems, nil
}

func Post[T any](ctx context.Context, c *Client, endpoint string, body any) (*T, error) {
	var target T
	res, err := c.restClient.R().
		SetContext(ctx).
		SetBody(body).
		SetResult(&target).
		Post(fullURL(baseUrl, endpoint))
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("error response from ConnectWise API: %s", res.String())
	}

	return &target, nil
}

func Put[T any](ctx context.Context, c *Client, endpoint string, body any) (*T, error) {
	var target T
	res, err := c.restClient.R().
		SetContext(ctx).
		SetBody(body).
		SetResult(&target).
		Put(fullURL(baseUrl, endpoint))
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		if res.StatusCode() == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error response from ConnectWise API: %s", res.String())
	}

	return &target, nil
}

func Patch[T any](ctx context.Context, c *Client, endpoint string, patchOps []PatchOp) (*T, error) {
	var target T
	res, err := c.restClient.R().
		SetContext(ctx).
		SetBody(patchOps).
		SetResult(&target).
		Patch(fullURL(baseUrl, endpoint))
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("error response from ConnectWise API: %s", res.String())
	}

	return &target, nil
}

func Delete(ctx context.Context, c *Client, endpoint string) error {
	res, err := c.restClient.R().
		SetContext(ctx).
		Delete(fullURL(baseUrl, endpoint))
	if err != nil {
		return err
	}

	if res.IsError() {
		if res.StatusCode() == http.StatusNotFound {
			return ErrNotFound
		}
		return fmt.Errorf("error response from ConnectWise API: %s", res.String())
	}

	return nil
}

func fullURL(base, endpoint string) string {
	return fmt.Sprintf("%s/%s", base, endpoint)
}

func parseLinkHeader(linkHeader, rel string) string {
	for link := range strings.SplitSeq(linkHeader, ",") {
		parts := strings.Split(strings.TrimSpace(link), ";")
		if len(parts) < 2 {
			continue
		}
		urlPart := strings.Trim(parts[0], "<>")
		relPart := strings.TrimSpace(parts[1])
		if relPart == fmt.Sprintf(`rel="%s"`, rel) {
			return urlPart
		}
	}

	return ""
}
