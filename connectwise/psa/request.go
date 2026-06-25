package psa

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/thecoretg/tctg-go/internal/httpx"
)

const (
	baseURL = "https://api-na.myconnectwise.net/v4_6_release/apis/3.0"
)

var ErrNotFound = errors.New("404 status returned")

func get[T any](ctx context.Context, c *Client, endpoint string, params map[string]string) (*T, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodGet, fullURL(baseURL, endpoint), params, nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		if res.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error response from ConnectWise API: %s", res.Body)
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

func getMany[T any](ctx context.Context, c *Client, endpoint string, params map[string]string) ([]T, error) {
	var allItems []T

	endpoint = fullURL(baseURL, endpoint)
	for endpoint != "" {
		res, err := httpx.Do(ctx, c.httpClient, http.MethodGet, endpoint, params, nil)
		if err != nil {
			return nil, err
		}

		if res.StatusCode >= 400 {
			if res.StatusCode == http.StatusNotFound {
				return nil, ErrNotFound
			}
			return nil, fmt.Errorf("error response from ConnectWise API: %s", res.Body)
		}

		var target []T
		if err := httpx.DecodeJSON(res.Body, &target); err != nil {
			return nil, err
		}

		allItems = append(allItems, target...)
		params = nil
		endpoint = parseLinkHeader(res.Header.Get("Link"), "next")
	}

	return allItems, nil
}

func post[T any](ctx context.Context, c *Client, endpoint string, body any) (*T, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodPost, fullURL(baseURL, endpoint), nil, body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("error response from ConnectWise API: %s", res.Body)
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

func put[T any](ctx context.Context, c *Client, endpoint string, body any) (*T, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodPut, fullURL(baseURL, endpoint), nil, body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		if res.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error response from ConnectWise API: %s", res.Body)
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

func patch[T any](ctx context.Context, c *Client, endpoint string, patchOps []PatchOp) (*T, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodPatch, fullURL(baseURL, endpoint), nil, patchOps)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("error response from ConnectWise API: %s", res.Body)
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

func del(ctx context.Context, c *Client, endpoint string) error {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodDelete, fullURL(baseURL, endpoint), nil, nil)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		if res.StatusCode == http.StatusNotFound {
			return ErrNotFound
		}
		return fmt.Errorf("error response from ConnectWise API: %s", res.Body)
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
