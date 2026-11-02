package ppe

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"

	"github.com/thecoretg/tctg-go/internal/httpx"
	"github.com/thecoretg/tctg-go/internal/httpxtest"
	"github.com/thecoretg/tctg-go/ratelimit"
)

// newStubClient builds a client with NewClient, so headers and the base URL
// are the real ones, then redirects its transport to srv.
func newStubClient(t *testing.T, srv *httptest.Server, stack string) *Client {
	t.Helper()
	c, err := NewClient(context.Background(), Config{Username: "admin@example.com", Password: "secret", Stack: stack})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	rt := c.httpClient.Transport.(*httpx.RetryTransport)
	rt.Base = httpxtest.Redirect(srv)
	rt.Retries = 0
	return c
}

func TestNewClientMissingCreds(t *testing.T) {
	if _, err := NewClient(context.Background(), Config{}); err == nil {
		t.Fatal("expected error for empty config, got nil")
	}
	if _, err := NewClient(context.Background(), Config{Username: "u"}); err == nil {
		t.Fatal("expected error for missing Password, got nil")
	}
}

func TestStackSelectsBaseURL(t *testing.T) {
	for stack, want := range map[string]string{
		"":    "https://us3.proofpointessentials.com/api/v1",
		"eu1": "https://eu1.proofpointessentials.com/api/v1",
	} {
		c, err := NewClient(context.Background(), Config{Username: "u", Password: "p", Stack: stack})
		if err != nil {
			t.Fatalf("NewClient: %v", err)
		}
		if c.baseURL != want {
			t.Errorf("stack %q: baseURL = %q, want %q", stack, c.baseURL, want)
		}
	}
}

// TestAuthHeadersAndPath checks the credentials go out as X-User/X-Password
// and that path segments, including an email address, are built and escaped.
func TestAuthHeadersAndPath(t *testing.T) {
	var gotUser, gotPass, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUser, gotPass = r.Header.Get("X-User"), r.Header.Get("X-Password")
		gotPath = r.URL.EscapedPath()
		_, _ = w.Write([]byte(`{"primary_email":"bob+x@acme.com","is_active":true}`))
	}))
	defer srv.Close()

	c := newStubClient(t, srv, "us2")
	u, err := c.GetUser(context.Background(), "acme.com", "bob+x@acme.com")
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if gotUser != "admin@example.com" || gotPass != "secret" {
		t.Errorf("auth headers = %q/%q", gotUser, gotPass)
	}
	if want := "/api/v1/orgs/acme.com/users/bob+x@acme.com"; gotPath != want {
		t.Errorf("path = %q, want %q", gotPath, want)
	}
	if u.IsActive == nil || !*u.IsActive {
		t.Errorf("IsActive = %v, want true", u.IsActive)
	}
}

// TestWriteSendsOnlySetFields checks a PUT body leaves out unset fields and
// that a 204 with no body succeeds.
func TestWriteSendsOnlySetFields(t *testing.T) {
	var gotMethod string
	var gotBody map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &gotBody)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newStubClient(t, srv, "")
	active := false
	if err := c.UpdateUser(context.Background(), "acme.com", "bob@acme.com", User{IsActive: &active}); err != nil {
		t.Fatalf("UpdateUser: %v", err)
	}
	if gotMethod != http.MethodPut {
		t.Errorf("method = %s, want PUT", gotMethod)
	}
	if len(gotBody) != 1 || gotBody["is_active"] != false {
		t.Errorf("body = %v, want only is_active=false", gotBody)
	}
}

func TestErrorMapping(t *testing.T) {
	cases := []struct {
		status int
		body   string
		want   error
	}{
		{http.StatusUnauthorized, "", ErrUnauthorized},
		{http.StatusForbidden, "not your org", ErrForbidden},
		{http.StatusNotFound, "", ErrNotFound},
		{http.StatusConflict, "exists", ErrConflict},
	}
	for _, tc := range cases {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tc.status)
			_, _ = w.Write([]byte(tc.body))
		}))
		c := newStubClient(t, srv, "")
		_, err := c.GetOrg(context.Background(), "acme.com")
		srv.Close()
		if !errors.Is(err, tc.want) {
			t.Errorf("status %d: err = %v, want %v", tc.status, err, tc.want)
		}
		if tc.body != "" && (err == nil || !strings.Contains(err.Error(), tc.body)) {
			t.Errorf("status %d: err %v does not carry body %q", tc.status, err, tc.body)
		}
	}
}

func TestUnprocessableCarriesBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"error_message":"primary_email is invalid"}`))
	}))
	defer srv.Close()

	c := newStubClient(t, srv, "")
	_, err := c.CreateUser(context.Background(), "acme.com", User{PrimaryEmail: "nope"})
	if err == nil || !strings.Contains(err.Error(), "primary_email is invalid") || !strings.Contains(err.Error(), "422") {
		t.Fatalf("err = %v, want 422 with body", err)
	}
}

func TestRateLimitError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	c := newStubClient(t, srv, "")
	_, err := c.Me(context.Background())
	if _, ok := errors.AsType[*ratelimit.Error](err); !ok {
		t.Fatalf("err = %v, want *ratelimit.Error", err)
	}
}

// TestListShapes checks lists decode both as the spec's bare array and the
// guide's keyed envelope, and that an unknown envelope fails loudly.
func TestListShapes(t *testing.T) {
	for name, body := range map[string]string{
		"bare":     `[{"primary_email":"a@acme.com"},{"primary_email":"b@acme.com"}]`,
		"envelope": `{"users":[{"primary_email":"a@acme.com"},{"primary_email":"b@acme.com"}]}`,
	} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(body))
		}))
		c := newStubClient(t, srv, "")
		users, err := c.ListUsers(context.Background(), "acme.com")
		srv.Close()
		if err != nil {
			t.Fatalf("%s: ListUsers: %v", name, err)
		}
		if len(users) != 2 || users[1].PrimaryEmail != "b@acme.com" {
			t.Errorf("%s: users = %+v", name, users)
		}
	}

	if _, err := decodeList[User](json.RawMessage(`{"people":[]}`), "users"); err == nil {
		t.Error("expected error for unknown envelope key, got nil")
	}
}

func TestStacks(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"endpoints":["https://us2.proofpointessentials.com","https://eu1.proofpointessentials.com"]}`))
	}))
	defer srv.Close()

	c := newStubClient(t, srv, "")
	stacks, err := c.Stacks(context.Background(), "acme.com")
	if err != nil {
		t.Fatalf("Stacks: %v", err)
	}
	if gotPath != "/api/v1/endpoints/acme.com" {
		t.Errorf("path = %q", gotPath)
	}
	if strings.Join(stacks, ",") != "us2,eu1" {
		t.Errorf("stacks = %v, want [us2 eu1]", stacks)
	}
}

func newTestClient(t *testing.T) *Client {
	t.Helper()
	_ = godotenv.Load("../../.env")
	c, err := NewClient(context.Background(), Config{
		Username: os.Getenv("PROOFPOINT_USERNAME"),
		Password: os.Getenv("PROOFPOINT_PASSWORD"),
		Stack:    os.Getenv("PROOFPOINT_STACK"),
	})
	if err != nil {
		t.Skip("skipping integration test:", err)
	}
	return c
}

// TestReadIntegration runs read-only calls against the live API, starting
// from the authenticated account's own organization.
func TestReadIntegration(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	me, err := c.Me(ctx)
	if err != nil {
		t.Fatalf("Me: %v", err)
	}
	domain := me.EntityPrimaryDomain
	t.Logf("auth OK as %s (%s) on %s", me.PrimaryEmail, me.EntityType, domain)

	stacks, err := c.Stacks(ctx, domain)
	if err != nil {
		t.Fatalf("Stacks: %v", err)
	}
	t.Logf("%s is on stacks %v", domain, stacks)

	if _, err := c.GetOrg(ctx, domain); err != nil {
		t.Fatalf("GetOrg: %v", err)
	}
	users, err := c.ListUsers(ctx, domain)
	if err != nil {
		t.Fatalf("ListUsers: %v", err)
	}
	t.Logf("%d user(s)", len(users))

	domains, err := c.ListDomains(ctx, domain)
	if err != nil {
		t.Fatalf("ListDomains: %v", err)
	}
	t.Logf("%d domain(s)", len(domains))

	if me.EntityType != EntityOrganization {
		orgs, err := c.ListChildOrgs(ctx, domain)
		if err != nil {
			t.Fatalf("ListChildOrgs: %v", err)
		}
		t.Logf("%d child org(s)", len(orgs))

		billing, err := c.GetBillingWithChildren(ctx, domain)
		if err != nil {
			t.Fatalf("GetBillingWithChildren: %v", err)
		}
		t.Logf("%d billing report(s)", len(billing))
	}
}
