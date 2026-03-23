package salesforce

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type queryResp[T any] struct {
	TotalSize      int    `json:"totalSize"`
	Done           bool   `json:"done"`
	NextRecordsURL string `json:"nextRecordsUrl"`
	Records        []T    `json:"records"`
}

var ErrNotFound = errors.New("404 status returned")

// Query executes a SOQL query and returns all records, following nextRecordsUrl
// pages until Salesforce signals done.
func Query[T any](ctx context.Context, c *Client, q string) ([]T, error) {
	var all []T
	url := c.endpointURL("query")
	params := map[string]string{"q": q}

	for {
		result, err := get[queryResp[T]](ctx, c, url, params)
		if err != nil {
			return nil, err
		}

		all = append(all, result.Records...)

		if result.Done {
			return all, nil
		}

		// subsequent pages use the full nextRecordsUrl path, no query params
		url = c.baseURL + result.NextRecordsURL
		params = nil
	}
}

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
		return nil, fmt.Errorf("error response from Salesforce: %s", res.String())
	}

	return &target, nil
}
