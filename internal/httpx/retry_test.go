package httpx

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/time/rate"

	"github.com/thecoretg/tctg-go/ratelimit"
)

// statusSequence serves the given statuses in order, then 200 for every later
// request, and counts how many requests arrived.
func statusSequence(t *testing.T, statuses ...int) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := int(hits.Add(1))
		if n <= len(statuses) {
			w.WriteHeader(statuses[n-1])
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)
	return srv, &hits
}

func TestRetryPolicy(t *testing.T) {
	tests := []struct {
		name       string
		check      func(*http.Response, error) bool
		statuses   []int
		wantStatus int
		wantHits   int32
	}{
		{"default retries 429", nil, []int{429}, 200, 2},
		{"default retries 5xx", nil, []int{503}, 200, 2},
		{"default leaves 4xx alone", nil, []int{400}, 400, 1},
		{"NoRetryOn429 stops at 429", NoRetryOn429, []int{429}, 429, 1},
		{"NoRetryOn429 still retries 5xx", NoRetryOn429, []int{503}, 200, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, hits := statusSequence(t, tt.statuses...)
			client := NewClient(nil, nil, Policy{Retries: 1, CheckRetry: tt.check})

			res, err := Do(context.Background(), client, http.MethodGet, srv.URL, nil, nil)
			if err != nil {
				t.Fatalf("Do() error = %v", err)
			}
			if res.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", res.StatusCode, tt.wantStatus)
			}
			if got := hits.Load(); got != tt.wantHits {
				t.Errorf("requests = %d, want %d", got, tt.wantHits)
			}
		})
	}
}

// countingLimiter records each Wait and fails once err is set.
type countingLimiter struct {
	waits atomic.Int32
	err   error
}

func (l *countingLimiter) Wait(context.Context) error {
	l.waits.Add(1)
	return l.err
}

// TestLimiterWaitsBeforeEveryAttempt pins that a retry spends quota like any
// other request: one 429 then a 200 is two sends, so two waits.
func TestLimiterWaitsBeforeEveryAttempt(t *testing.T) {
	srv, hits := statusSequence(t, http.StatusTooManyRequests)
	lim := &countingLimiter{}
	client := NewClient(nil, nil, Policy{Retries: 1, Limiter: lim})

	res, err := Do(context.Background(), client, http.MethodGet, srv.URL, nil, nil)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", res.StatusCode)
	}
	if got, want := lim.waits.Load(), hits.Load(); got != want || got != 2 {
		t.Errorf("waits = %d, requests = %d, want 2 of each", got, want)
	}
}

// TestLimiterErrorStopsRequest pins that a limiter refusing to wait (its
// context ended, or the wait would outlast the deadline) sends nothing.
func TestLimiterErrorStopsRequest(t *testing.T) {
	srv, hits := statusSequence(t)
	wantErr := errors.New("would exceed deadline")
	client := NewClient(nil, nil, Policy{Limiter: &countingLimiter{err: wantErr}})

	_, err := Do(context.Background(), client, http.MethodPost, srv.URL, nil, map[string]string{"a": "b"})
	if !errors.Is(err, ratelimit.ErrRateLimited) || !strings.Contains(err.Error(), wantErr.Error()) {
		t.Errorf("Do() error = %v, want ErrRateLimited carrying %q", err, wantErr)
	}
	if got := hits.Load(); got != 0 {
		t.Errorf("requests = %d, want 0", got)
	}
}

func TestRateLimitErrorReadsRetryAfter(t *testing.T) {
	res := &Response{StatusCode: 429, Body: []byte("slow down"), Header: http.Header{"Retry-After": {"30"}}}
	e := RateLimitError("huntress", res)
	if e.Service != "huntress" || e.StatusCode != 429 || e.RetryAfter != 30*time.Second || e.Message != "slow down" {
		t.Errorf("RateLimitError() = %+v", e)
	}
}

// TestLimiterDeadlineIsRateLimited pins that a limiter refusing because the
// wait would outlast the deadline surfaces as ErrRateLimited, while a
// cancelled context stays a context error.
func TestLimiterDeadlineIsRateLimited(t *testing.T) {
	srv, hits := statusSequence(t)
	lim := rate.NewLimiter(rate.Every(time.Hour), 1)
	lim.Allow() // spend the only token, so the next one is an hour away
	client := NewClient(nil, nil, Policy{Limiter: lim})

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	_, err := Do(ctx, client, http.MethodGet, srv.URL, nil, nil)
	if !errors.Is(err, ratelimit.ErrRateLimited) {
		t.Errorf("deadline: error = %v, want ErrRateLimited", err)
	}

	ctx, cancel = context.WithCancel(context.Background())
	cancel()
	_, err = Do(ctx, client, http.MethodGet, srv.URL, nil, nil)
	if !errors.Is(err, context.Canceled) || errors.Is(err, ratelimit.ErrRateLimited) {
		t.Errorf("cancelled: error = %v, want context.Canceled only", err)
	}

	if got := hits.Load(); got != 0 {
		t.Errorf("requests = %d, want 0", got)
	}
}
