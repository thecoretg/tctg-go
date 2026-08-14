package salesforce

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/thecoretg/tctg-go/internal/httpx"
)

type queryResp[T any] struct {
	TotalSize      int    `json:"totalSize"`
	Done           bool   `json:"done"`
	NextRecordsURL string `json:"nextRecordsUrl"`
	Records        []T    `json:"records"`
}

var ErrNotFound = errors.New("404 status returned")

// Query executes a SOQL query and returns all records, following nextRecordsUrl
// pages until Salesforce signals done. Records are decoded into T exactly as
// Salesforce returns them, including the attributes key; use QueryRecords for
// normalized keys.
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
			break
		}

		url = c.baseURL + result.NextRecordsURL
		params = nil
	}

	return all, nil
}

// QueryRecords executes a SOQL query and returns each record with its keys
// lowercased, __c suffixes stripped, and the attributes key dropped.
func QueryRecords(ctx context.Context, c *Client, q string) ([]map[string]any, error) {
	records, err := Query[map[string]any](ctx, c, q)
	if err != nil {
		return nil, err
	}

	return SimplifyRecords(records), nil
}

func SimplifyRecords(records []map[string]any) []map[string]any {
	for i, r := range records {
		records[i] = simplifyRecord(r)
	}

	return records
}

func simplifyRecord(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		if k == "attributes" {
			continue
		}
		out[strings.TrimSuffix(strings.ToLower(k), "__c")] = v
	}
	return out
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
		return nil, fmt.Errorf("error response from Salesforce: %s [%d]", res.Body, res.StatusCode)
	}

	var target T
	if err := httpx.DecodeJSON(res.Body, &target); err != nil {
		return nil, err
	}
	return &target, nil
}
