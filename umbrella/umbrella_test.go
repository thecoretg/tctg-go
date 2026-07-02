package umbrella

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"

	"github.com/thecoretg/tctg-go/internal/httpx"
)

// TestNewClientMissingCreds verifies config validation reports the missing
// fields rather than building an unusable client.
func TestNewClientMissingCreds(t *testing.T) {
	if _, err := NewClient(context.Background(), Config{}); err == nil {
		t.Fatal("expected error for empty config, got nil")
	}
	if _, err := NewClient(context.Background(), Config{ClientID: "id"}); err == nil {
		t.Fatal("expected error for missing ClientSecret, got nil")
	}
}

// TestCollectPagedStopsOnShortPage checks the pagination helper walks full
// pages and stops as soon as a page returns fewer than limit items.
func TestCollectPagedStopsOnShortPage(t *testing.T) {
	const limit = 3
	pages := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8}} // last page is short

	var calls int
	all, err := collectPaged(1, limit, func(page, l int) ([]int, error) {
		calls++
		if l != limit {
			t.Fatalf("limit passed to fetch = %d, want %d", l, limit)
		}
		return pages[page-1], nil
	})
	if err != nil {
		t.Fatalf("collectPaged: %v", err)
	}
	if calls != 3 {
		t.Fatalf("fetch called %d times, want 3", calls)
	}
	if len(all) != 8 {
		t.Fatalf("collected %d items, want 8", len(all))
	}
}

// TestListCustomersPaging drives the real get helper and pagination against a
// stub server: a full first page followed by a short second page.
func TestListCustomersPaging(t *testing.T) {
	const limit = 2
	var gotPages []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPages = append(gotPages, r.URL.Query().Get("page"))
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Query().Get("page") {
		case "1":
			_, _ = w.Write([]byte(`[{"customerId":1},{"customerId":2}]`))
		default:
			_, _ = w.Write([]byte(`[{"customerId":3}]`))
		}
	}))
	defer srv.Close()

	c := &Client{httpClient: httpx.NewClient(nil, nil, 0)}
	all, err := collectPaged(1, limit, func(page, l int) ([]Customer, error) {
		params := map[string]string{"page": strconv.Itoa(page), "limit": strconv.Itoa(l)}
		result, err := get[[]Customer](context.Background(), c, srv.URL, params)
		if err != nil {
			return nil, err
		}
		return *result, nil
	})
	if err != nil {
		t.Fatalf("collectPaged: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("collected %d customers, want 3", len(all))
	}
	if fmt.Sprint(gotPages) != "[1 2]" {
		t.Fatalf("requested pages = %v, want [1 2]", gotPages)
	}
}

// TestAPIErrorDecodesMessage verifies error responses surface the API message.
func TestAPIErrorDecodesMessage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message":"Invalid request"}`))
	}))
	defer srv.Close()

	c := &Client{httpClient: httpx.NewClient(nil, nil, 0)}
	_, err := get[[]Customer](context.Background(), c, srv.URL, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if want := "Invalid request"; !contains(err.Error(), want) {
		t.Fatalf("error %q does not contain %q", err.Error(), want)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func newTestClient(t *testing.T) *Client {
	t.Helper()
	_ = godotenv.Load("../.env")
	c, err := NewClient(context.Background(), Config{
		ClientID:     os.Getenv("UMBRELLA_CLIENT_ID"),
		ClientSecret: os.Getenv("UMBRELLA_CLIENT_SECRET"),
	})
	if err != nil {
		t.Skip("skipping integration test:", err)
	}
	return c
}

// TestListCustomersIntegration exercises auth against the live Umbrella API. A
// bad or missing credential surfaces as a 401/403 error.
func TestListCustomersIntegration(t *testing.T) {
	c := newTestClient(t)
	customers, err := c.ListCustomers(context.Background(), map[string]string{"limit": "1"})
	if err != nil {
		t.Fatalf("ListCustomers: %v", err)
	}
	t.Logf("auth OK; got %d customer(s)", len(customers))
}

// TestListAllCustomersIntegration verifies the auto-paging helper against the
// live API.
func TestListAllCustomersIntegration(t *testing.T) {
	c := newTestClient(t)
	all, err := c.ListAllCustomers(context.Background())
	if err != nil {
		t.Fatalf("ListAllCustomers: %v", err)
	}
	t.Logf("fetched %d customer(s) across all pages", len(all))
}
