package psa

import (
	"maps"
	"strconv"
)

// maxPageSize is the largest page ConnectWise will return. When a caller sets a
// limit, requesting pages of this size keeps the number of round trips to the
// minimum needed to satisfy it.
const maxPageSize = 1000

// ListOption configures a list request.
type ListOption func(*listConfig)

type listConfig struct {
	limit int
}

// WithLimit caps the number of records a list call returns, stopping pagination
// as soon as the cap is reached. Without it a list call follows every page,
// which on a large tenant is many round trips and a large response. A limit of
// zero or less means no cap.
//
// When the returned slice length equals the limit the server may hold more
// records; narrow the request with ConnectWise conditions rather than raising
// the limit if you need a complete answer.
func WithLimit(n int) ListOption {
	return func(c *listConfig) { c.limit = n }
}

func newListConfig(opts []ListOption) listConfig {
	var cfg listConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	return cfg
}

// withPageSize returns params with pageSize set to size, leaving a
// caller-supplied pageSize and the caller's own map untouched.
func withPageSize(params map[string]string, size int) map[string]string {
	if _, ok := params["pageSize"]; ok {
		return params
	}

	out := maps.Clone(params)
	if out == nil {
		out = map[string]string{}
	}
	out["pageSize"] = strconv.Itoa(size)

	return out
}
