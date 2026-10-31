package ratelimit

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestErrorMessageAndMatching(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{"vendor with Retry-After", &Error{Service: "huntress", StatusCode: 429, RetryAfter: 30 * time.Second, Message: "slow down"},
			"huntress rate limit exceeded [429]; retry after 30s: slow down"},
		{"vendor without Retry-After", &Error{Service: "salesforce", StatusCode: 403, Message: "REQUEST_LIMIT_EXCEEDED"},
			"salesforce rate limit exceeded [403]: REQUEST_LIMIT_EXCEEDED"},
		{"client limiter", &Error{Message: "would exceed context deadline"},
			"rate limit exceeded: would exceed context deadline"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
			wrapped := fmt.Errorf("listing: %w", tt.err)
			if !errors.Is(wrapped, ErrRateLimited) {
				t.Error("errors.Is(wrapped, ErrRateLimited) = false, want true")
			}
			var e *Error
			if !errors.As(wrapped, &e) || e != tt.err {
				t.Error("errors.As did not recover the *Error")
			}
		})
	}
}
