package huntress

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/thecoretg/tctg-go/internal/httpx"
	"github.com/thecoretg/tctg-go/internal/httpxtest"
)

// TestListAllKeepsRawItems pins that ListAll[json.RawMessage] follows the page
// token, builds the URL from the path, reads the items from the key named by
// the path's last segment, and hands each item back byte for byte.
func TestListAllKeepsRawItems(t *testing.T) {
	var paths, tokens []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		tokens = append(tokens, r.URL.Query().Get("page_token"))
		if r.URL.Query().Get("page_token") == "" {
			_, _ = w.Write([]byte(`{"agents":[{"id":1,"unmodelled":{"a":1}}],"pagination":{"next_page_token":"tok2"}}`))
			return
		}
		_, _ = w.Write([]byte(`{"agents":[{"id":2}],"pagination":{}}`))
	}))
	defer srv.Close()
	c := &Client{httpClient: httpx.NewClient(httpxtest.Redirect(srv), nil, httpx.Policy{})}

	items, err := c.ListAll[json.RawMessage](context.Background(), "accounts/7/agents", nil)
	if err != nil {
		t.Fatalf("ListAll() error = %v", err)
	}
	if got := len(items); got != 2 || string(items[0]) != `{"id":1,"unmodelled":{"a":1}}` {
		t.Errorf("items = %s, want both raw agents", items)
	}
	if strings.Join(paths, ",") != "/v1/accounts/7/agents,/v1/accounts/7/agents" || strings.Join(tokens, ",") != ",tok2" {
		t.Errorf("requests = %v with tokens %v", paths, tokens)
	}

	typed, err := c.ListAllAccountAgents(context.Background(), 7, nil)
	if err != nil || len(typed) != 2 || typed[1].ID != 2 {
		t.Errorf("ListAllAccountAgents() = %+v, %v", typed, err)
	}
}

func TestListAllRejectsMissingKey(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"something_else":[],"pagination":{}}`))
	}))
	defer srv.Close()
	c := &Client{httpClient: httpx.NewClient(httpxtest.Redirect(srv), nil, httpx.Policy{})}

	_, err := c.ListAll[json.RawMessage](context.Background(), "agents", nil)
	if err == nil || !strings.Contains(err.Error(), `no "agents" array`) {
		t.Errorf("ListAll() error = %v, want a missing-key error", err)
	}
}
