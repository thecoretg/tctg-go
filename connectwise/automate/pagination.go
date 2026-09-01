package automate

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"strconv"
)

// defaultListAllPageSize is the page size the ListAll* helpers request when the
// caller does not set one. Automate's own default (30) would make these helpers
// issue many more round trips than necessary.
const defaultListAllPageSize = 1000

// totalCount reads Automate's Total-Count response header (the total number of
// records across all pages), returning -1 when it is absent or unparseable.
func totalCount(h http.Header) int {
	v := h.Get("Total-Count")
	if v == "" {
		return -1
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return -1
	}
	return n
}

// listAll fetches every page of a list endpoint, advancing Automate's page /
// pagesize query params and using the Total-Count response header to know when
// to stop. It stops on a short page even when Total-Count is absent. The
// supplied params are not modified; a caller-provided pagesize/page is honored
// (page sets the starting page).
func (c *Client) listAll[T any](ctx context.Context, endpoint string, params map[string]string) ([]T, error) {
	p := maps.Clone(params)
	if p == nil {
		p = map[string]string{}
	}

	pageSize := defaultListAllPageSize
	if v, err := strconv.Atoi(p["pagesize"]); err == nil && v > 0 {
		pageSize = v
	}
	p["pagesize"] = strconv.Itoa(pageSize)

	page := 1
	if v, err := strconv.Atoi(p["page"]); err == nil && v > 0 {
		page = v
	}

	var all []T
	for {
		p["page"] = strconv.Itoa(page)
		items, hdr, err := c.getWithHeaders[[]T](ctx, endpoint, p)
		if err != nil {
			return nil, err
		}
		all = append(all, (*items)...)

		// A short page is always the last page. Otherwise stop once we have
		// collected the server-reported total (when it provides one).
		if len(*items) < pageSize {
			break
		}
		if total := totalCount(hdr); total >= 0 && len(all) >= total {
			break
		}
		page++
	}
	return all, nil
}

// ListAllCompanies returns every company across all pages (see ListCompanies).
func (c *Client) ListAllCompanies(ctx context.Context, params map[string]string) ([]Company, error) {
	result, err := c.listAll[Company](ctx, "cwa/api/v1/Clients", params)
	if err != nil {
		return nil, fmt.Errorf("list all companies: %w", err)
	}
	return result, nil
}

// ListAllLocations returns every location across all pages (see ListLocations).
func (c *Client) ListAllLocations(ctx context.Context, params map[string]string) ([]Location, error) {
	result, err := c.listAll[Location](ctx, "cwa/api/v1/Locations", params)
	if err != nil {
		return nil, fmt.Errorf("list all locations: %w", err)
	}
	return result, nil
}

// ListAllComputers returns every computer across all pages (see ListComputers).
func (c *Client) ListAllComputers(ctx context.Context, params map[string]string) ([]Computer, error) {
	result, err := c.listAll[Computer](ctx, "cwa/api/v1/Computers", params)
	if err != nil {
		return nil, fmt.Errorf("list all computers: %w", err)
	}
	return result, nil
}

// ListAllComputerDrives returns every computer drive across all pages (see
// ListComputerDrives).
func (c *Client) ListAllComputerDrives(ctx context.Context, params map[string]string) ([]ComputerDrive, error) {
	result, err := c.listAll[ComputerDrive](ctx, "cwa/api/v1/Computers/Drives", params)
	if err != nil {
		return nil, fmt.Errorf("list all computer drives: %w", err)
	}
	return result, nil
}
