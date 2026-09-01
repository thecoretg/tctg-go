package webex

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/thecoretg/tctg-go/internal/httpx"
)

const (
	baseURL = "https://webexapis.com/v1"
)

var ErrNotFound = errors.New("404 status returned")

// Get issues a GET request and decodes the JSON response into T.
func (c *Client) Get[T any](ctx context.Context, endpoint string, params map[string]string) (*T, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodGet, fullURL(baseURL, endpoint), params, nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		if res.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error response from Webex API: %s", res.Body)
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

// GetMany issues GET requests, following the Link header until every page
// has been collected.
func (c *Client) GetMany[T any](ctx context.Context, endpoint string, params map[string]string) ([]T, error) {
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
			return nil, fmt.Errorf("error response from Webex API: %s", res.Body)
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

// Post issues a POST request with body and decodes the JSON response into T.
func (c *Client) Post[T any](ctx context.Context, endpoint string, body any) (*T, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodPost, fullURL(baseURL, endpoint), nil, body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("error response from Webex API: %s", res.Body)
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

// Put issues a PUT request with body and decodes the JSON response into T.
func (c *Client) Put[T any](ctx context.Context, endpoint string, body any) (*T, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodPut, fullURL(baseURL, endpoint), nil, body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		if res.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error response from Webex API: %s", res.Body)
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

// Delete issues a DELETE request, discarding any response body.
func (c *Client) Delete(ctx context.Context, endpoint string) error {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodDelete, fullURL(baseURL, endpoint), nil, nil)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		if res.StatusCode == http.StatusNotFound {
			return ErrNotFound
		}
		return fmt.Errorf("error response from Webex API: %s", res.Body)
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
