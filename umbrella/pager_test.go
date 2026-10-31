package umbrella

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thecoretg/tctg-go/internal/httpx"
	"github.com/thecoretg/tctg-go/internal/httpxtest"
)

// TestListAllKeepsRawItems pins that ListAll[json.RawMessage] builds the URL
// from a host-relative path, sends the caller's params with every page, and
// stops at the first page shorter than the page size.
func TestListAllKeepsRawItems(t *testing.T) {
	var pages []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if r.URL.Path != "/deployments/v2/roamingcomputers" || q.Get("status") != "Off" || q.Get("limit") != "2" {
			t.Errorf("request = %s?%s", r.URL.Path, r.URL.RawQuery)
		}
		pages = append(pages, q.Get("page"))
		if q.Get("page") == "1" {
			_, _ = w.Write([]byte(`[{"deviceId":"a","extra":true},{"deviceId":"b"}]`))
			return
		}
		_, _ = w.Write([]byte(`[{"deviceId":"c"}]`))
	}))
	defer srv.Close()
	c := &Client{httpClient: httpx.NewClient(httpxtest.Redirect(srv), nil, httpx.Policy{})}

	items, err := c.ListAll[json.RawMessage](context.Background(), "deployments/v2/roamingcomputers",
		map[string]string{"status": "Off", "limit": "2"})
	if err != nil {
		t.Fatalf("ListAll() error = %v", err)
	}
	if len(items) != 3 || string(items[0]) != `{"deviceId":"a","extra":true}` {
		t.Errorf("items = %s, want 3 raw computers", items)
	}
	if len(pages) != 2 || pages[0] != "1" || pages[1] != "2" {
		t.Errorf("pages = %v, want [1 2]", pages)
	}
}
