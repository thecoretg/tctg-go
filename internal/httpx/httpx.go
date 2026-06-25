// Package httpx provides the shared HTTP plumbing for tctg-go's API clients:
// a retrying RoundTripper that injects static headers, plus small helpers for
// issuing JSON requests and decoding responses. It exists so every client gets
// identical, tested retry and encoding behavior without depending on a
// third-party HTTP library.
package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

const (
	backoffBase = 100 * time.Millisecond
	backoffMax  = 2 * time.Second
)

// NewClient builds an *http.Client whose transport injects headers on every
// request and retries safe failures up to retries times. base may be nil, in
// which case http.DefaultTransport is used; pass a non-nil base (e.g. an
// oauth2 transport) to layer retry on top of it.
func NewClient(base http.RoundTripper, headers map[string]string, retries int) *http.Client {
	return &http.Client{
		Transport: &RetryTransport{Base: base, Headers: headers, Retries: retries},
	}
}

// RetryTransport wraps a base RoundTripper. It sets any configured Headers that
// the request does not already carry, then retries transport errors, 429, and
// 5xx responses up to Retries times using exponential backoff with jitter. It
// honors a Retry-After header when present and aborts as soon as the request
// context is cancelled.
type RetryTransport struct {
	Base    http.RoundTripper
	Headers map[string]string
	Retries int
}

func (t *RetryTransport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}

func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Inject static headers on a clone so the caller's request is untouched.
	if len(t.Headers) > 0 {
		req = req.Clone(req.Context())
		for k, v := range t.Headers {
			if req.Header.Get(k) == "" {
				req.Header.Set(k, v)
			}
		}
	}

	var (
		resp *http.Response
		err  error
	)
	for attempt := 0; ; attempt++ {
		resp, err = t.base().RoundTrip(req)
		if attempt >= t.Retries || !shouldRetry(resp, err) {
			return resp, err
		}

		wait := retryAfter(resp)
		if wait <= 0 {
			wait = backoff(attempt)
		}

		// A request with a body can only be replayed if GetBody is available.
		var nextBody io.ReadCloser
		if req.Body != nil && req.Body != http.NoBody {
			if req.GetBody == nil {
				return resp, err
			}
			b, gerr := req.GetBody()
			if gerr != nil {
				return resp, err
			}
			nextBody = b
		}

		// Drain and close the failed response so the connection is reusable.
		if resp != nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
		if nextBody != nil {
			req.Body = nextBody
		}

		timer := time.NewTimer(wait)
		select {
		case <-req.Context().Done():
			timer.Stop()
			return nil, req.Context().Err()
		case <-timer.C:
		}
	}
}

func shouldRetry(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}
	return resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500
}

func retryAfter(resp *http.Response) time.Duration {
	if resp == nil {
		return 0
	}
	v := resp.Header.Get("Retry-After")
	if v == "" {
		return 0
	}
	if secs, err := strconv.Atoi(v); err == nil {
		return time.Duration(secs) * time.Second
	}
	if when, err := http.ParseTime(v); err == nil {
		if d := time.Until(when); d > 0 {
			return d
		}
	}
	return 0
}

// backoff returns an exponentially growing delay with half jitter, capped at
// backoffMax: roughly 100ms, 200ms, 400ms, ... each randomized within [d/2, d].
func backoff(attempt int) time.Duration {
	d := backoffBase << attempt
	if d > backoffMax || d <= 0 {
		d = backoffMax
	}
	half := int64(d) / 2
	return time.Duration(half + rand.Int64N(half+1))
}

// Response is the decoded result of a request: the HTTP status, the full body,
// and the response headers.
type Response struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

// Do issues an HTTP request against client. A non-nil body is JSON encoded and
// sent with Content-Type: application/json; non-empty params are added as query
// parameters. The returned error is non-nil only for request construction,
// transport, or body-read failures — HTTP error statuses are reported through
// Response.StatusCode and left for the caller to interpret.
func Do(ctx context.Context, client *http.Client, method, rawURL string, params map[string]string, body any) (*Response, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, rawURL, reqBody)
	if err != nil {
		return nil, err
	}

	if len(params) > 0 {
		q := req.URL.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{StatusCode: resp.StatusCode, Body: data, Header: resp.Header}, nil
}

// DecodeJSON unmarshals data into target, treating an empty body as a no-op so
// that responses without content (e.g. 204) do not produce a decode error.
func DecodeJSON(data []byte, target any) error {
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, target)
}
