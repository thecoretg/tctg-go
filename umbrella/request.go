package umbrella

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/thecoretg/tctg-go/internal/httpx"
)

var ErrNotFound = errors.New("404 status returned")

// apiError formats an error response from the Umbrella API, decoding the
// standard {message} body when present for a clearer message.
func apiError(res *httpx.Response) error {
	var er ErrorResponse
	if err := json.Unmarshal(res.Body, &er); err == nil && er.Message != "" {
		return fmt.Errorf("error response from Umbrella API: %s [%d]", er.Message, res.StatusCode)
	}
	return fmt.Errorf("error response from Umbrella API: %s [%d]", res.Body, res.StatusCode)
}

// Get issues a GET request and decodes the JSON response into T.
func (c *Client) Get[T any](ctx context.Context, url string, params map[string]string) (*T, error) {
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

// Post issues a POST request with body and decodes the JSON response into T.
func (c *Client) Post[T any](ctx context.Context, url string, body any) (*T, error) {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodPost, url, nil, body)
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

// Put issues a PUT request with body and decodes the JSON response into T.
func (c *Client) Put[T any](ctx context.Context, url string, body any) (*T, error) {
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

// Delete issues a DELETE request, discarding any response body.
func (c *Client) Delete(ctx context.Context, url string) error {
	res, err := httpx.Do(ctx, c.httpClient, http.MethodDelete, url, nil, nil)
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
