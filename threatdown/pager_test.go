package threatdown

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thecoretg/tctg-go/internal/httpx"
	"github.com/thecoretg/tctg-go/internal/httpxtest"
)

func TestListAllFollowsCursorParam(t *testing.T) {
	var cursors []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/oneview/v1/users" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		cursors = append(cursors, r.URL.Query().Get("cursor"))
		if r.URL.Query().Get("cursor") == "" {
			_, _ = w.Write([]byte(`{"users":[{"id":"u1","x":1}],"next_cursor":"c2"}`))
			return
		}
		_, _ = w.Write([]byte(`{"users":[{"id":"u2"}],"next_cursor":""}`))
	}))
	defer srv.Close()
	c := &Client{httpClient: httpx.NewClient(httpxtest.Redirect(srv), nil, httpx.Policy{})}

	items, err := c.ListAll[json.RawMessage](context.Background(), "v1/users", nil)
	if err != nil {
		t.Fatalf("ListAll() error = %v", err)
	}
	if len(items) != 2 || string(items[0]) != `{"id":"u1","x":1}` {
		t.Errorf("items = %s", items)
	}
	if len(cursors) != 2 || cursors[1] != "c2" {
		t.Errorf("cursors = %v, want [ c2]", cursors)
	}
}

func TestSearchAllSendsCursorInBody(t *testing.T) {
	var bodies []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/oneview/v1/accounts/acc-1/endpoints" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		b, _ := io.ReadAll(r.Body)
		bodies = append(bodies, string(b))
		if len(bodies) == 1 {
			_, _ = w.Write([]byte(`{"endpoints":[{"display_name":"a"}],"next_cursor":"c2"}`))
			return
		}
		_, _ = w.Write([]byte(`{"endpoints":[{"display_name":"b"}]}`))
	}))
	defer srv.Close()
	c := &Client{httpClient: httpx.NewClient(httpxtest.Redirect(srv), nil, httpx.Policy{})}

	eps, err := c.SearchEndpoints(context.Background(), "acc-1", map[string]string{"page_size": "100"})
	if err != nil {
		t.Fatalf("SearchEndpoints() error = %v", err)
	}
	if len(eps) != 2 || eps[1].DisplayName != "b" {
		t.Errorf("endpoints = %+v", eps)
	}
	want := []string{`{"page_size":"100"}`, `{"next_cursor":"c2","page_size":"100"}`}
	if len(bodies) != 2 || bodies[0] != want[0] || bodies[1] != want[1] {
		t.Errorf("bodies = %q, want %q", bodies, want)
	}
}
