package addigy

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

// TestListAllAdvancesPageInBody pins that ListAll resends the caller's query
// with "page" advanced, stops at metadata.page_count, and keeps raw items.
func TestListAllAdvancesPageInBody(t *testing.T) {
	var bodies []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v2/o/org-1/devices" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		b, _ := io.ReadAll(r.Body)
		bodies = append(bodies, string(b))
		var q struct{ Page int }
		_ = json.Unmarshal(b, &q)
		if q.Page == 1 {
			_, _ = w.Write([]byte(`{"items":[{"agentid":"a","x":1}],"metadata":{"page":1,"page_count":2}}`))
			return
		}
		_, _ = w.Write([]byte(`{"items":[{"agentid":"b"}],"metadata":{"page":2,"page_count":2}}`))
	}))
	defer srv.Close()
	c := &Client{httpClient: httpx.NewClient(httpxtest.Redirect(srv), nil, httpx.Policy{}), orgID: "org-1"}

	items, err := c.ListAll[json.RawMessage](context.Background(), "o/org-1/devices", DeviceFilter{PerPage: 1})
	if err != nil {
		t.Fatalf("ListAll() error = %v", err)
	}
	if len(items) != 2 || string(items[0]) != `{"agentid":"a","x":1}` {
		t.Errorf("items = %s", items)
	}
	want := []string{`{"page":1,"per_page":1}`, `{"page":2,"per_page":1}`}
	if len(bodies) != 2 || bodies[0] != want[0] || bodies[1] != want[1] {
		t.Errorf("bodies = %q, want %q", bodies, want)
	}

	devices, err := c.SearchAllDevices(context.Background(), DeviceFilter{PerPage: 1})
	if err != nil || len(devices) != 2 || devices[1].AgentID != "b" {
		t.Errorf("SearchAllDevices() = %+v, %v", devices, err)
	}
}
