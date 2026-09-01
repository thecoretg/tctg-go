package salesforce

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/joho/godotenv"

	"github.com/thecoretg/tctg-go/internal/httpx"
)

// newStubClient points a Client at srv so the query helpers can be exercised
// without Salesforce credentials.
func newStubClient(srv *httptest.Server) *Client {
	return &Client{httpClient: httpx.NewClient(nil, nil, 0), baseURL: srv.URL}
}

func newTestClient(t *testing.T) *Client {
	t.Helper()
	_ = godotenv.Load("../.env")
	client, err := NewClient(context.Background(), Config{
		ClientID:       os.Getenv("SALESFORCE_CLIENT_ID"),
		ClientSecret:   os.Getenv("SALESFORCE_CLIENT_SECRET"),
		CompanyURLName: os.Getenv("SALESFORCE_COMPANY_URL_NAME"),
	})
	if err != nil {
		t.Skip("skipping integration test:", err)
	}
	return client
}

func TestQuery(t *testing.T) {
	c := newTestClient(t)

	records, err := c.Query[map[string]any](context.Background(),
		"SELECT Id, Name, Phone, Support_Agreement__c FROM Account",
	)
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	t.Logf("got %d records", len(records))
}

// TestQueryRecordsIntegration checks the normalizing wrapper against the live
// API: every returned key must already be lowercased and suffix-free.
func TestQueryRecordsIntegration(t *testing.T) {
	c := newTestClient(t)

	records, err := c.QueryRecords(context.Background(),
		"SELECT Id, Name, Phone, Support_Agreement__c FROM Account WHERE Type = 'Customer'",
	)
	if err != nil {
		t.Fatalf("QueryRecords: %v", err)
	}
	for _, r := range records {
		for k := range r {
			if k == "attributes" {
				t.Errorf("attributes key survived normalization")
			}
			if k != strings.ToLower(k) || strings.HasSuffix(k, "__c") {
				t.Errorf("key %q is not normalized", k)
			}
		}
	}
	t.Logf("got %d records", len(records))
}

func TestSimplifyRecord(t *testing.T) {
	nested := map[string]any{"attributes": map[string]any{"type": "Contact"}, "Name": "Nested"}

	tests := []struct {
		name string
		in   map[string]any
		want map[string]any
	}{
		{
			name: "drops attributes, lowercases keys, strips __c",
			in: map[string]any{
				"attributes":             map[string]any{"type": "Account", "url": "/services/data/v64.0/sobjects/Account/SomeID"},
				"Id":                     "SomeID",
				"Name":                   "My Client Name",
				"Endpoint_Protection__c": nil,
			},
			want: map[string]any{
				"id":                  "SomeID",
				"name":                "My Client Name",
				"endpoint_protection": nil,
			},
		},
		{
			name: "empty record",
			in:   map[string]any{},
			want: map[string]any{},
		},
		{
			name: "only attributes",
			in:   map[string]any{"attributes": map[string]any{"type": "Account"}},
			want: map[string]any{},
		},
		{
			name: "preserves value types",
			in: map[string]any{
				"Amount__c":   1234.5,
				"IsActive__c": true,
				"Tags__c":     []any{"a", "b"},
			},
			want: map[string]any{
				"amount":   1234.5,
				"isactive": true,
				"tags":     []any{"a", "b"},
			},
		},
		{
			// Relationship fields end in __r, not __c: the suffix is kept and
			// nested records are returned untouched, attributes included.
			name: "relationship field is lowercased but not unwrapped",
			in:   map[string]any{"Owner__r": nested},
			want: map[string]any{"owner__r": nested},
		},
		{
			// Child subqueries arrive as a nested query response; normalization
			// does not recurse into it.
			name: "child subquery left intact",
			in: map[string]any{
				"Contacts": map[string]any{"totalSize": 1.0, "done": true, "records": []any{nested}},
			},
			want: map[string]any{
				"contacts": map[string]any{"totalSize": 1.0, "done": true, "records": []any{nested}},
			},
		},
		{
			// Only one __c suffix is removed, so a doubled suffix keeps the rest.
			name: "single suffix strip",
			in:   map[string]any{"Weird__c__c": 1},
			want: map[string]any{"weird__c": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := simplifyRecord(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("simplifyRecord() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

// TestSimplifyRecordDoesNotMutateInput guards the caller's copy: the source map
// must survive normalization unchanged.
func TestSimplifyRecordDoesNotMutateInput(t *testing.T) {
	in := map[string]any{"attributes": map[string]any{"type": "Account"}, "Name__c": "x"}

	simplifyRecord(in)

	if len(in) != 2 || in["Name__c"] != "x" {
		t.Errorf("input map was modified: %#v", in)
	}
}

func TestSimplifyRecords(t *testing.T) {
	records := []map[string]any{
		{"attributes": map[string]any{"type": "Account"}, "Id": "1", "Support_Agreement__c": "gold"},
		{"Id": "2"},
	}

	got := SimplifyRecords(records)

	want := []map[string]any{
		{"id": "1", "support_agreement": "gold"},
		{"id": "2"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("SimplifyRecords() = %#v, want %#v", got, want)
	}
	// The slice is normalized in place, so the caller's slice sees the result.
	if !reflect.DeepEqual(records, want) {
		t.Errorf("input slice = %#v, want it normalized in place", records)
	}
}

func TestSimplifyRecordsEmpty(t *testing.T) {
	if got := SimplifyRecords(nil); len(got) != 0 {
		t.Errorf("SimplifyRecords(nil) = %#v, want empty", got)
	}
	if got := SimplifyRecords([]map[string]any{}); len(got) != 0 {
		t.Errorf("SimplifyRecords(empty) = %#v, want empty", got)
	}
}

// TestQueryPagination drives Query against a stub server that pages via
// nextRecordsUrl, and confirms raw Salesforce keys are left untouched.
func TestQueryPagination(t *testing.T) {
	const nextPath = "/services/data/" + versionTag + "/query/01g000000000001-2000"

	var gotPaths, gotQueries []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPaths = append(gotPaths, r.URL.Path)
		gotQueries = append(gotQueries, r.URL.Query().Get("q"))
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == nextPath {
			_, _ = w.Write([]byte(`{"totalSize":3,"done":true,"records":[{"Id":"3"}]}`))
			return
		}
		_, _ = w.Write([]byte(`{"totalSize":3,"done":false,"nextRecordsUrl":"` + nextPath + `",
			"records":[{"attributes":{"type":"Account"},"Id":"1","Support_Agreement__c":"gold"},{"Id":"2"}]}`))
	}))
	defer srv.Close()

	records, err := newStubClient(srv).Query[map[string]any](context.Background(), "SELECT Id FROM Account")
	if err != nil {
		t.Fatalf("Query: %v", err)
	}

	if len(records) != 3 {
		t.Fatalf("got %d records, want 3", len(records))
	}
	for i, want := range []string{"1", "2", "3"} {
		if records[i]["Id"] != want {
			t.Errorf("record %d Id = %v, want %q", i, records[i]["Id"], want)
		}
	}
	if _, ok := records[0]["attributes"]; !ok {
		t.Error("Query should return raw records, attributes included")
	}
	if _, ok := records[0]["Support_Agreement__c"]; !ok {
		t.Error("Query should not normalize keys")
	}

	wantPaths := []string{"/services/data/" + versionTag + "/query", nextPath}
	if !reflect.DeepEqual(gotPaths, wantPaths) {
		t.Errorf("requested paths = %v, want %v", gotPaths, wantPaths)
	}
	// The q param belongs to the first request only; the follow-up URL carries
	// its own cursor.
	if gotQueries[0] != "SELECT Id FROM Account" {
		t.Errorf("first request q = %q", gotQueries[0])
	}
	if gotQueries[1] != "" {
		t.Errorf("second request resent q = %q", gotQueries[1])
	}
}

// TestQueryRecordsNormalizes checks the wrapper normalizes every page, not just
// the first.
func TestQueryRecordsNormalizes(t *testing.T) {
	const nextPath = "/services/data/" + versionTag + "/query/01g000000000001-2000"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == nextPath {
			_, _ = w.Write([]byte(`{"totalSize":2,"done":true,
				"records":[{"attributes":{"type":"Account"},"Id":"2","Endpoint_Protection__c":"crowdstrike"}]}`))
			return
		}
		_, _ = w.Write([]byte(`{"totalSize":2,"done":false,"nextRecordsUrl":"` + nextPath + `",
			"records":[{"attributes":{"type":"Account"},"Id":"1","Support_Agreement__c":"gold"}]}`))
	}))
	defer srv.Close()

	got, err := newStubClient(srv).QueryRecords(context.Background(), "SELECT Id FROM Account")
	if err != nil {
		t.Fatalf("QueryRecords: %v", err)
	}

	want := []map[string]any{
		{"id": "1", "support_agreement": "gold"},
		{"id": "2", "endpoint_protection": "crowdstrike"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("QueryRecords() = %#v, want %#v", got, want)
	}
}

func TestQueryEmptyResult(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"totalSize":0,"done":true,"records":[]}`))
	}))
	defer srv.Close()

	records, err := newStubClient(srv).QueryRecords(context.Background(), "SELECT Id FROM Account")
	if err != nil {
		t.Fatalf("QueryRecords: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("got %d records, want 0", len(records))
	}
}

func TestQueryErrors(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		body    string
		wantErr error
	}{
		{name: "not found", status: http.StatusNotFound, body: `[]`, wantErr: ErrNotFound},
		{name: "malformed soql", status: http.StatusBadRequest, body: `[{"errorCode":"MALFORMED_QUERY"}]`},
		{name: "unauthorized", status: http.StatusUnauthorized, body: `[{"errorCode":"INVALID_SESSION_ID"}]`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			c := newStubClient(srv)
			if _, err := c.Query[map[string]any](context.Background(), "SELECT Id FROM Account"); err == nil {
				t.Fatal("Query: expected error, got nil")
			} else if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("Query: got %v, want %v", err, tt.wantErr)
			}

			if _, err := c.QueryRecords(context.Background(), "SELECT Id FROM Account"); err == nil {
				t.Fatal("QueryRecords: expected error, got nil")
			}
		})
	}
}

// TestQueryPropagatesCancellation makes sure an aborted context stops paging
// instead of looping against the API.
func TestQueryPropagatesCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"totalSize":1,"done":true,"records":[{"Id":"1"}]}`))
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := newStubClient(srv).Query[map[string]any](ctx, "SELECT Id FROM Account"); !errors.Is(err, context.Canceled) {
		t.Fatalf("Query: got %v, want context.Canceled", err)
	}
}
