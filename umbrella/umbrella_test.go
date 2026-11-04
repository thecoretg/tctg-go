package umbrella

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/time/rate"

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

// TestListCustomersPaging drives the real Get helper and pagination against a
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

	c := &Client{httpClient: httpx.NewClient(nil, nil, httpx.Policy{})}
	all, err := collectPaged(1, limit, func(page, l int) ([]Customer, error) {
		params := map[string]string{"page": strconv.Itoa(page), "limit": strconv.Itoa(l)}
		result, err := c.Get[[]Customer](context.Background(), srv.URL, params)
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

	c := &Client{httpClient: httpx.NewClient(nil, nil, httpx.Policy{})}
	_, err := c.Get[[]Customer](context.Background(), srv.URL, nil)
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

// TestForOrgScopesRequestsToTheChild checks a ForOrg client uses a token issued for its
// child org (the token request carries X-Umbrella-OrgId), reuses it, shares the parent's
// limiter, and leaves the parent on the provider's token.
func TestForOrgScopesRequestsToTheChild(t *testing.T) {
	var tokenOrgs []string
	var used []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/token" {
			org := r.Header.Get(orgHeader)
			tokenOrgs = append(tokenOrgs, org)
			_, _ = fmt.Fprintf(w, `{"access_token":"tok-%s","token_type":"bearer","expires_in":3600}`, org)
			return
		}
		used = append(used, strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "))
		_, _ = w.Write([]byte(`{"organizationId":"8675309","fingerprint":"abc123","userId":"42","region":"us"}`))
	}))
	defer srv.Close()

	creds := &clientcredentials.Config{ClientID: "id", ClientSecret: "secret", TokenURL: srv.URL + "/token"}
	parent := newClient(context.Background(), creds, httpx.Policy{Limiter: rate.NewLimiter(rate.Inf, 1)})
	child := parent.ForOrg(8675309)

	info, err := child.Get[OrgInfo](context.Background(), srv.URL+"/orgInfo", nil)
	if err != nil {
		t.Fatalf("child get: %v", err)
	}
	if info.OrganizationID != "8675309" || info.Fingerprint != "abc123" || info.UserID != "42" {
		t.Errorf("decoded %+v", *info)
	}
	if again := parent.ForOrg(8675309); again != child {
		t.Error("ForOrg built a second client for the same customer; its token would not be reused")
	}
	if _, err := child.Get[OrgInfo](context.Background(), srv.URL+"/orgInfo", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := parent.Get[OrgInfo](context.Background(), srv.URL+"/orgInfo", nil); err != nil {
		t.Fatalf("parent get: %v", err)
	}
	if want := []string{"8675309", ""}; !slices.Equal(tokenOrgs, want) {
		t.Errorf("token requests for orgs %q, want %q (one per org, the child's first)", tokenOrgs, want)
	}
	if want := []string{"tok-8675309", "tok-8675309", "tok-"}; !slices.Equal(used, want) {
		t.Errorf("tokens used %q, want %q", used, want)
	}
	if child.policy.Limiter != parent.policy.Limiter {
		t.Error("the child has a limiter of its own; it must spend the parent's budget")
	}
}

// TestForOrgOrganizationInfoIntegration checks against the live API that a ForOrg client
// reads each customer's own OrgInfo: each answers with its own organization, whose id is
// its customerId, and neither is the provider's.
func TestForOrgOrganizationInfoIntegration(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()
	customers, err := c.ListCustomers(ctx, map[string]string{"limit": "2"})
	if err != nil {
		t.Fatalf("ListCustomers: %v", err)
	}
	if len(customers) < 2 {
		t.Skip("needs two customers")
	}
	provider, err := c.GetOrganizationInfo(ctx)
	if err != nil {
		t.Fatalf("provider's own OrgInfo: %v", err)
	}
	seen := map[string]int{}
	for _, cu := range customers[:2] {
		info, err := c.ForOrg(cu.CustomerID).GetOrganizationInfo(ctx)
		if err != nil {
			t.Fatalf("GetOrganizationInfo for customer %d: %v", cu.CustomerID, err)
		}
		if info.OrganizationID != strconv.Itoa(cu.CustomerID) {
			t.Errorf("customer %d: organizationId %q", cu.CustomerID, info.OrganizationID)
		}
		if info.Fingerprint == "" || info.UserID == "" {
			t.Errorf("customer %d: an OrgInfo field is empty", cu.CustomerID)
		}
		if info.OrganizationID == provider.OrganizationID {
			t.Errorf("customer %d answered with the provider's own organization", cu.CustomerID)
		}
		if other, dup := seen[info.OrganizationID]; dup {
			t.Errorf("customers %d and %d answered with the same organization", other, cu.CustomerID)
		}
		seen[info.OrganizationID] = cu.CustomerID
	}
	t.Logf("2 customers, 2 distinct organizations, neither the provider's")
}
