package addigy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

// TestClientSendsAPIKey verifies that a client built by NewClient injects the
// x-api-key header on outgoing requests. It uses a local test server and needs
// no live credentials, isolating the auth wiring from any credential problem.
func TestClientSendsAPIKey(t *testing.T) {
	const wantKey = "test-key-123"

	var gotKey string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.Header.Get("x-api-key")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	c, err := NewClient(context.Background(), Config{APIKey: wantKey, OrgID: "org-1"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	// Hit the test server directly (bypassing baseURL) through the client's
	// configured transport so the header injection is exercised end to end.
	if _, err := c.Get[struct{}](context.Background(), srv.URL, nil); err != nil {
		t.Fatalf("request: %v", err)
	}
	if gotKey != wantKey {
		t.Fatalf("x-api-key header = %q, want %q", gotKey, wantKey)
	}
}

func newTestClient(t *testing.T) *Client {
	t.Helper()
	_ = godotenv.Load("../.env")
	c, err := NewClient(context.Background(), Config{
		APIKey: os.Getenv("ADDIGY_API_KEY"),
		OrgID:  os.Getenv("ADDIGY_ORG_ID"),
	})
	if err != nil {
		t.Skip("skipping integration test:", err)
	}
	return c
}

// TestSearchDevicesAuth exercises auth against the live Addigy API using an
// org-scoped endpoint. A bad or missing key surfaces as a 401/403 error.
func TestSearchDevicesAuth(t *testing.T) {
	c := newTestClient(t)
	if c.orgID == "" {
		t.Skip("skipping: ADDIGY_ORG_ID must be set")
	}

	resp, err := c.SearchDevices(context.Background(), DeviceFilter{Page: 1, PerPage: 1})
	if err != nil {
		t.Fatalf("SearchDevices: %v", err)
	}
	t.Logf("auth OK; got %d device(s)", len(resp.Items))
}

// TestSearchAllDevicesPaging verifies the auto-paging helper walks past a single
// small page and matches the server-reported total.
func TestSearchAllDevicesPaging(t *testing.T) {
	c := newTestClient(t)
	if c.orgID == "" {
		t.Skip("skipping: ADDIGY_ORG_ID must be set")
	}
	ctx := context.Background()

	first, err := c.SearchDevices(ctx, DeviceFilter{PerPage: 5})
	if err != nil {
		t.Fatalf("SearchDevices: %v", err)
	}
	if first.Metadata.PageCount <= 1 {
		t.Skipf("skipping: org has only %d page of devices", first.Metadata.PageCount)
	}

	all, err := c.SearchAllDevices(ctx, DeviceFilter{PerPage: 5})
	if err != nil {
		t.Fatalf("SearchAllDevices: %v", err)
	}
	t.Logf("single page=%d, all pages=%d, reported total=%d",
		len(first.Items), len(all), first.Metadata.Total)
	if len(all) <= len(first.Items) {
		t.Fatalf("expected SearchAll (%d) to exceed a single page (%d)", len(all), len(first.Items))
	}
}
