package httpx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// TestDoKeepsCompleteBodyAfterStreamReset pins the ConnectWise Automate 404:
// the server sends a complete response and then resets the HTTP/2 stream, so
// io.ReadAll fails even though every byte arrived. Dropping the response would
// hide the status behind an opaque transport error.
func TestDoKeepsCompleteBodyAfterStreamReset(t *testing.T) {
	const payload = `{"Message":"Not found."}`

	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(payload)))
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(payload))
		w.(http.Flusher).Flush()

		panic(http.ErrAbortHandler)
	}))
	srv.EnableHTTP2 = true
	srv.StartTLS()
	defer srv.Close()

	client := srv.Client()

	res, err := Do(context.Background(), client, http.MethodGet, srv.URL, nil, nil)
	if err != nil {
		t.Fatalf("Do() error = %v, want the 404 to survive the stream reset", err)
	}
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want %d", res.StatusCode, http.StatusNotFound)
	}
	if string(res.Body) != payload {
		t.Errorf("body = %q, want %q", res.Body, payload)
	}
}

// TestDoStillFailsOnTruncatedBody keeps a genuinely short read an error: the
// declared length is never reached, so the response cannot be trusted.
func TestDoStillFailsOnTruncatedBody(t *testing.T) {
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "512")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("only a few bytes"))
		w.(http.Flusher).Flush()

		panic(http.ErrAbortHandler)
	}))
	srv.EnableHTTP2 = true
	srv.StartTLS()
	defer srv.Close()

	if _, err := Do(context.Background(), srv.Client(), http.MethodGet, srv.URL, nil, nil); err == nil {
		t.Fatal("Do() error = nil, want a truncated body to fail")
	}
}
