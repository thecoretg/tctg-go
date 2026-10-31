package salesforce

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thecoretg/tctg-go/ratelimit"
)

func TestDailyQuotaIsRateLimited(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		want   bool
	}{
		{"daily quota", http.StatusForbidden, `[{"message":"TotalRequests Limit exceeded.","errorCode":"REQUEST_LIMIT_EXCEEDED"}]`, true},
		{"429", http.StatusTooManyRequests, ``, true},
		{"permission denied", http.StatusForbidden, `[{"message":"insufficient access","errorCode":"INSUFFICIENT_ACCESS"}]`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			_, err := newStubClient(srv).Get[struct{}](context.Background(), srv.URL, nil)
			if err == nil {
				t.Fatal("Get() error = nil, want an error")
			}
			if got := errors.Is(err, ratelimit.ErrRateLimited); got != tt.want {
				t.Errorf("errors.Is(err, ErrRateLimited) = %v, want %v (err: %v)", got, tt.want, err)
			}
		})
	}
}
