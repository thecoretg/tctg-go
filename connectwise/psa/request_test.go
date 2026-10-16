package psa

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const securityBody = `{"code":"Security","message":"You do not have security permission to perform this action."}`

// TestCheckStatusEveryMethod confirms each request helper maps 404 and 403 to their sentinels,
// and that a 403 keeps ConnectWise's body in the message.
func TestCheckStatusEveryMethod(t *testing.T) {
	ctx := context.Background()
	calls := map[string]func(c *Client) error{
		"Get": func(c *Client) error {
			_, err := c.Get[Ticket](ctx, "service/tickets/1", nil)
			return err
		},
		"GetMany": func(c *Client) error {
			_, err := c.GetMany[Ticket](ctx, "service/tickets", nil)
			return err
		},
		"Post": func(c *Client) error {
			_, err := c.Post[Ticket](ctx, "service/tickets", &Ticket{})
			return err
		},
		"Put": func(c *Client) error {
			_, err := c.Put[Ticket](ctx, "service/tickets/1", &Ticket{})
			return err
		},
		"Patch": func(c *Client) error {
			_, err := c.Patch[Ticket](ctx, "service/tickets/1", []PatchOp{})
			return err
		},
		"Delete": func(c *Client) error {
			return c.Delete(ctx, "service/tickets/1")
		},
	}

	statuses := []struct {
		code    int
		body    string
		want    error
		notWant error
		inMsg   string
	}{
		{code: http.StatusNotFound, body: `{"code":"NotFound"}`, want: ErrNotFound, notWant: ErrForbidden},
		{code: http.StatusForbidden, body: securityBody, want: ErrForbidden, notWant: ErrNotFound, inMsg: "You do not have security permission"},
		{code: http.StatusBadRequest, body: `{"code":"InvalidObject"}`, inMsg: "InvalidObject"},
	}

	for name, call := range calls {
		for _, st := range statuses {
			t.Run(name+"/"+http.StatusText(st.code), func(t *testing.T) {
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(st.code)
					_, _ = w.Write([]byte(st.body))
				}))
				t.Cleanup(srv.Close)
				c := &Client{httpClient: srv.Client(), baseURL: srv.URL}

				err := call(c)
				if err == nil {
					t.Fatal("err = nil, want an error")
				}
				if st.want != nil && !errors.Is(err, st.want) {
					t.Errorf("err = %v, want errors.Is %v", err, st.want)
				}
				if st.notWant != nil && errors.Is(err, st.notWant) {
					t.Errorf("err = %v, should not match %v", err, st.notWant)
				}
				if st.want == nil && (errors.Is(err, ErrNotFound) || errors.Is(err, ErrForbidden)) {
					t.Errorf("err = %v, want neither sentinel", err)
				}
				if st.inMsg != "" && !strings.Contains(err.Error(), st.inMsg) {
					t.Errorf("err = %q, want it to contain %q", err, st.inMsg)
				}
			})
		}
	}
}
