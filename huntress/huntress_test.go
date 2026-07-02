package huntress

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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
	if _, err := NewClient(context.Background(), Config{APIKey: "key"}); err == nil {
		t.Fatal("expected error for missing APISecret, got nil")
	}
}

// TestCollectPagedStopsOnEmptyToken checks the cursor pagination helper walks
// pages until NextPageToken is empty.
func TestCollectPagedStopsOnEmptyToken(t *testing.T) {
	pages := map[string]struct {
		items []int
		next  string
	}{
		"":     {[]int{1, 2, 3}, "tok2"},
		"tok2": {[]int{4, 5, 6}, "tok3"},
		"tok3": {[]int{7, 8}, ""},
	}

	var calls int
	all, err := collectPaged(func(token string) ([]int, Pagination, error) {
		calls++
		p := pages[token]
		return p.items, Pagination{NextPageToken: p.next}, nil
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

// TestListPagingAndAuth drives the real get helper, cursor pagination, and
// HTTP Basic auth against a stub server: a first page with a next-page token
// followed by a final page with an empty token.
func TestListPagingAndAuth(t *testing.T) {
	var gotTokens []string
	var gotAuth string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotTokens = append(gotTokens, r.URL.Query().Get("page_token"))
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Query().Get("page_token") {
		case "":
			_, _ = w.Write([]byte(`{"organizations":[{"id":1},{"id":2}],"pagination":{"next_page_token":"tok2"}}`))
		default:
			_, _ = w.Write([]byte(`{"organizations":[{"id":3}],"pagination":{}}`))
		}
	}))
	defer srv.Close()

	c, err := NewClient(context.Background(), Config{APIKey: "key", APISecret: "secret"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	all, err := collectPaged(func(token string) ([]Organization, Pagination, error) {
		resp, err := get[OrganizationsResponse](context.Background(), c, srv.URL, pageParams(nil, token))
		if err != nil {
			return nil, Pagination{}, err
		}
		return resp.Organizations, resp.Pagination, nil
	})
	if err != nil {
		t.Fatalf("collectPaged: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("collected %d organizations, want 3", len(all))
	}
	if fmt.Sprint(gotTokens) != "[ tok2]" {
		t.Fatalf("requested page tokens = %v, want [ tok2]", gotTokens)
	}
	wantAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("key:secret"))
	if gotAuth != wantAuth {
		t.Fatalf("Authorization header = %q, want %q", gotAuth, wantAuth)
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
	_, err := get[AccountsResponse](context.Background(), c, srv.URL, nil)
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
		APIKey:    os.Getenv("HUNTRESS_API_KEY"),
		APISecret: os.Getenv("HUNTRESS_API_SECRET"),
	})
	if err != nil {
		t.Skip("skipping integration test:", err)
	}
	return c
}

// TestGetAccountIntegration exercises auth against the live Huntress API. A bad
// or missing credential surfaces as a 401/403 error.
func TestGetAccountIntegration(t *testing.T) {
	c := newTestClient(t)
	account, err := c.GetAccount(context.Background())
	if err != nil {
		t.Fatalf("GetAccount: %v", err)
	}
	t.Logf("auth OK; account %d (%s)", account.ID, account.Name)
}

// TestListAllOrganizationsIntegration verifies the auto-paging helper against
// the live API.
func TestListAllOrganizationsIntegration(t *testing.T) {
	c := newTestClient(t)
	all, err := c.ListAllOrganizations(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListAllOrganizations: %v", err)
	}
	t.Logf("fetched %d organization(s) across all pages", len(all))
}
