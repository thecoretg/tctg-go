package ratelimit

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErrRateLimited matches every *Error with errors.Is, so a caller can check
// for a rate limit without caring which vendor or which side refused.
var ErrRateLimited = errors.New("rate limited")

// Error reports a request refused for rate limiting. Either the vendor refused
// it (StatusCode is set: a 429 that outlasted the client's retries, or
// Salesforce's 403 REQUEST_LIMIT_EXCEEDED), or the client's own Limiter did
// (StatusCode is 0: the wait for quota would have outlasted the context's
// deadline, so nothing was sent).
type Error struct {
	// Service names the API, e.g. "huntress". Empty when the client's own
	// Limiter refused.
	Service string
	// StatusCode is the vendor's HTTP status, or 0 when nothing was sent.
	StatusCode int
	// RetryAfter is the wait the vendor asked for in a Retry-After header, or
	// 0 when it sent none.
	RetryAfter time.Duration
	// Message carries the vendor's response body or the limiter's reason.
	Message string
}

func (e *Error) Error() string {
	var b strings.Builder
	if e.Service != "" {
		b.WriteString(e.Service)
		b.WriteByte(' ')
	}
	b.WriteString("rate limit exceeded")
	if e.StatusCode != 0 {
		fmt.Fprintf(&b, " [%d]", e.StatusCode)
	}
	if e.RetryAfter > 0 {
		fmt.Fprintf(&b, "; retry after %s", e.RetryAfter)
	}
	if e.Message != "" {
		b.WriteString(": ")
		b.WriteString(e.Message)
	}
	return b.String()
}

// Is reports whether target is ErrRateLimited.
func (e *Error) Is(target error) bool { return target == ErrRateLimited }
