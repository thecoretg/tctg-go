package automate

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListAllKeepsRawItems(t *testing.T) {
	var pages []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cwa/api/v1/Computers" || r.URL.Query().Get("pagesize") != "2" {
			t.Errorf("request = %s?%s", r.URL.Path, r.URL.RawQuery)
		}
		pages = append(pages, r.URL.Query().Get("page"))
		w.Header().Set("Total-Count", "3")
		if r.URL.Query().Get("page") == "1" {
			_, _ = w.Write([]byte(`[{"Id":1,"Extra":"x"},{"Id":2}]`))
			return
		}
		_, _ = w.Write([]byte(`[{"Id":3}]`))
	}))
	defer srv.Close()
	c, err := NewClient(context.Background(), Config{ServerURL: srv.URL, ClientID: "id", Token: "tok"})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	items, err := c.ListAll[json.RawMessage](context.Background(), "cwa/api/v1/Computers", map[string]string{"pagesize": "2"})
	if err != nil {
		t.Fatalf("ListAll() error = %v", err)
	}
	if len(items) != 3 || string(items[0]) != `{"Id":1,"Extra":"x"}` {
		t.Errorf("items = %s", items)
	}
	if len(pages) != 2 {
		t.Errorf("pages = %v, want 2 requests", pages)
	}
}
