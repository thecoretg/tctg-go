package psa

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

// pagedServer serves pages of `perPage` companies up to `total`, advertising
// the next page through the Link header the way ConnectWise does.
func pagedServer(t *testing.T, total, perPage int, requests *atomic.Int64) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	mux.HandleFunc("/company/companies", func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)

		page := 1
		if v := r.URL.Query().Get("page"); v != "" {
			if _, err := fmt.Sscanf(v, "%d", &page); err != nil {
				t.Errorf("bad page param %q: %v", v, err)
			}
		}

		size := perPage
		if v := r.URL.Query().Get("pageSize"); v != "" {
			if _, err := fmt.Sscanf(v, "%d", &size); err != nil {
				t.Errorf("bad pageSize param %q: %v", v, err)
			}
		}

		start := (page - 1) * size
		end := min(start+size, total)
		if start > end {
			start = end
		}

		body := "["
		for i := start; i < end; i++ {
			if i > start {
				body += ","
			}
			body += fmt.Sprintf(`{"id":%d,"name":"Company %d"}`, i+1, i+1)
		}
		body += "]"

		if end < total {
			next := fmt.Sprintf("%s/company/companies?page=%d&pageSize=%d", srv.URL, page+1, size)
			w.Header().Set("Link", fmt.Sprintf(`<%s>; rel="next"`, next))
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(body)); err != nil {
			t.Errorf("writing body: %v", err)
		}
	})

	return srv
}

func testClient(t *testing.T, srv *httptest.Server) *Client {
	t.Helper()

	return &Client{httpClient: srv.Client(), baseURL: srv.URL}
}

func TestListCompaniesFollowsEveryPageByDefault(t *testing.T) {
	var requests atomic.Int64
	srv := pagedServer(t, 250, 25, &requests)

	companies, err := testClient(t, srv).ListCompanies(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListCompanies: %v", err)
	}

	if len(companies) != 250 {
		t.Errorf("companies = %d, want 250", len(companies))
	}
	if got := requests.Load(); got != 10 {
		t.Errorf("requests = %d, want 10", got)
	}
}

func TestWithLimitCapsResultAndStopsEarly(t *testing.T) {
	var requests atomic.Int64
	srv := pagedServer(t, 5000, 25, &requests)

	companies, err := testClient(t, srv).ListCompanies(context.Background(), nil, WithLimit(10))
	if err != nil {
		t.Fatalf("ListCompanies: %v", err)
	}

	if len(companies) != 10 {
		t.Fatalf("companies = %d, want 10", len(companies))
	}
	if companies[0].Name != "Company 1" || companies[9].Name != "Company 10" {
		t.Errorf("got %q..%q, want Company 1..Company 10", companies[0].Name, companies[9].Name)
	}
	if got := requests.Load(); got != 1 {
		t.Errorf("requests = %d, want 1 - the limit should size the page", got)
	}
}

func TestWithLimitStopsAtMaxPageSize(t *testing.T) {
	var requests atomic.Int64
	srv := pagedServer(t, 5000, 25, &requests)

	companies, err := testClient(t, srv).ListCompanies(context.Background(), nil, WithLimit(2500))
	if err != nil {
		t.Fatalf("ListCompanies: %v", err)
	}

	if len(companies) != 2500 {
		t.Errorf("companies = %d, want 2500", len(companies))
	}
	if got := requests.Load(); got != 3 {
		t.Errorf("requests = %d, want 3 pages of %d", got, maxPageSize)
	}
}

func TestWithLimitLargerThanTotalReturnsEverything(t *testing.T) {
	var requests atomic.Int64
	srv := pagedServer(t, 7, 25, &requests)

	companies, err := testClient(t, srv).ListCompanies(context.Background(), nil, WithLimit(100))
	if err != nil {
		t.Fatalf("ListCompanies: %v", err)
	}

	if len(companies) != 7 {
		t.Errorf("companies = %d, want 7", len(companies))
	}
}

func TestWithLimitHonorsCallerPageSize(t *testing.T) {
	var requests atomic.Int64
	srv := pagedServer(t, 5000, 25, &requests)

	params := map[string]string{"pageSize": "5"}
	companies, err := testClient(t, srv).ListCompanies(context.Background(), params, WithLimit(12))
	if err != nil {
		t.Fatalf("ListCompanies: %v", err)
	}

	if len(companies) != 12 {
		t.Errorf("companies = %d, want 12", len(companies))
	}
	if params["pageSize"] != "5" {
		t.Errorf("caller params mutated: pageSize = %q, want 5", params["pageSize"])
	}
	if got := requests.Load(); got != 3 {
		t.Errorf("requests = %d, want 3 pages of 5", got)
	}
}

func TestWithLimitZeroMeansNoCap(t *testing.T) {
	var requests atomic.Int64
	srv := pagedServer(t, 60, 25, &requests)

	companies, err := testClient(t, srv).ListCompanies(context.Background(), nil, WithLimit(0))
	if err != nil {
		t.Fatalf("ListCompanies: %v", err)
	}

	if len(companies) != 60 {
		t.Errorf("companies = %d, want 60", len(companies))
	}
}

func TestListCompaniesPassesConditionsThrough(t *testing.T) {
	got := make(chan string, 1)

	mux := http.NewServeMux()
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	mux.HandleFunc("/company/companies", func(w http.ResponseWriter, r *http.Request) {
		got <- r.URL.Query().Get("conditions")
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte("[]")); err != nil {
			t.Errorf("writing body: %v", err)
		}
	})

	conditions := `status/name="Active" and deletedFlag=false`
	if _, err := testClient(t, srv).ListCompanies(context.Background(), map[string]string{"conditions": conditions}, WithLimit(5)); err != nil {
		t.Fatalf("ListCompanies: %v", err)
	}

	if sent := <-got; sent != conditions {
		t.Errorf("conditions = %q, want %q", sent, conditions)
	}
}
