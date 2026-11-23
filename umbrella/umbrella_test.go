package umbrella

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/time/rate"

	"github.com/thecoretg/tctg-go/internal/httpx"
	"github.com/thecoretg/tctg-go/ratelimit"
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

// countingLimiter counts the waits a client makes, one per attempt.
type countingLimiter struct {
	ratelimit.Limiter
	waits atomic.Int32
}

func (l *countingLimiter) Wait(ctx context.Context) error {
	l.waits.Add(1)
	return l.Limiter.Wait(ctx)
}

// TestRefusedTokenIsNotRetried checks a token request answered with a 4xx other than 429
// costs one attempt, one limiter wait and one token request, while a 429 or 5xx token
// answer is still retried.
func TestRefusedTokenIsNotRetried(t *testing.T) {
	for _, tc := range []struct {
		status       int
		wantRequests int32
	}{
		{http.StatusBadRequest, 1},
		{http.StatusUnauthorized, 1},
		{http.StatusForbidden, 1},
		{http.StatusTooManyRequests, 2},
		{http.StatusServiceUnavailable, 2},
	} {
		t.Run(strconv.Itoa(tc.status), func(t *testing.T) {
			var tokenRequests, apiRequests atomic.Int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/token" {
					apiRequests.Add(1)
					return
				}
				tokenRequests.Add(1)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(`{"error":"invalid_client"}`))
			}))
			defer srv.Close()

			lim := &countingLimiter{Limiter: rate.NewLimiter(rate.Inf, 1)}
			creds := &clientcredentials.Config{ClientID: "id", ClientSecret: "secret", TokenURL: srv.URL + "/token", AuthStyle: oauth2.AuthStyleInHeader}
			// One retry keeps the retried cases fast; the policy is NewClient's otherwise.
			parent := newClient(context.Background(), creds, httpx.Policy{Retries: 1, Limiter: lim, CheckRetry: checkRetry})
			_, err := parent.ForOrg(8675309).Get[OrgInfo](context.Background(), srv.URL+"/orgInfo", nil)

			var re *oauth2.RetrieveError
			if !errors.As(err, &re) || re.Response.StatusCode != tc.status {
				t.Fatalf("err = %v, want a %d *oauth2.RetrieveError", err, tc.status)
			}
			if got := tokenRequests.Load(); got != tc.wantRequests {
				t.Errorf("token requests = %d, want %d", got, tc.wantRequests)
			}
			if got := lim.waits.Load(); got != tc.wantRequests {
				t.Errorf("limiter waits = %d, want %d", got, tc.wantRequests)
			}
			if apiRequests.Load() != 0 {
				t.Error("the API was called without a token")
			}
		})
	}
}

// TestNewClientDoesNotRetryRefusedToken checks NewClient's policy carries checkRetry.
func TestNewClientDoesNotRetryRefusedToken(t *testing.T) {
	c, err := NewClient(context.Background(), Config{ClientID: "id", ClientSecret: "secret", Limiter: ratelimit.None})
	if err != nil {
		t.Fatal(err)
	}
	if c.creds.AuthStyle != oauth2.AuthStyleInHeader {
		t.Error("NewClient leaves the auth style to auto-detect, which sends a refused token request twice")
	}
	refused := &oauth2.RetrieveError{Response: &http.Response{StatusCode: http.StatusUnauthorized}}
	if c.policy.CheckRetry == nil || c.policy.CheckRetry(nil, fmt.Errorf("wrapped: %w", refused)) {
		t.Error("NewClient's policy retries a refused token")
	}
	if c.ForOrg(1).policy.CheckRetry == nil {
		t.Error("ForOrg's child lost the check")
	}
}

// TestForOrgRefusedTokenIntegration checks against the live API that a child org the key
// may not act for is refused on the token request with a 4xx, after one attempt.
func TestForOrgRefusedTokenIntegration(t *testing.T) {
	c := newTestClient(t)
	lim := &countingLimiter{Limiter: c.policy.Limiter}
	c.policy.Limiter = lim
	_, err := c.ForOrg(1).GetOrganizationInfo(context.Background())
	var re *oauth2.RetrieveError
	if !errors.As(err, &re) {
		t.Fatalf("err = %v, want a *oauth2.RetrieveError", err)
	}
	if code := re.Response.StatusCode; code < 400 || code >= 500 || code == http.StatusTooManyRequests {
		t.Fatalf("token answer %d, want a 4xx refusal other than 429", code)
	}
	if got := lim.waits.Load(); got != 1 {
		t.Errorf("limiter waits = %d, want 1", got)
	}
	t.Logf("org 1 refused with %d after one attempt: %v", re.Response.StatusCode, err)
}
