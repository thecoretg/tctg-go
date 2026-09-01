package automate

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

// dataCall is one non-login request the fake server received.
type dataCall struct {
	method string
	path   string
	auth   string
	body   string
}

// fakeAutomate stands in for an Automate server. It issues incrementing bearer
// tokens from the apitoken endpoint and hands every other request to respond,
// which decides the status and body. Setting revoked makes the server reject a
// token that it previously issued, which is how token expiry is simulated.
type fakeAutomate struct {
	mu      sync.Mutex
	logins  int
	issued  string
	revoked string
	calls   []dataCall

	loginStatus int
	respond     func(f *fakeAutomate, call dataCall, n int) (int, string)
}

func (f *fakeAutomate) handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/"+apiTokenEndpoint || r.URL.Path == "/"+apiTokenEndpoint+"/refresh" {
			f.logins++
			if f.loginStatus != 0 && f.loginStatus != http.StatusOK {
				w.WriteHeader(f.loginStatus)
				_, _ = w.Write([]byte(`{}`))
				return
			}
			f.issued = fmt.Sprintf("token-%d", f.logins)
			_, _ = fmt.Fprintf(w, `{"AccessToken":%q,"TokenType":"Bearer"}`, f.issued)
			return
		}

		body, _ := io.ReadAll(r.Body)
		call := dataCall{method: r.Method, path: r.URL.Path, auth: r.Header.Get("Authorization"), body: string(body)}
		f.calls = append(f.calls, call)

		status, payload := f.respond(f, call, len(f.calls))
		w.WriteHeader(status)
		_, _ = w.Write([]byte(payload))
	})
}

func (f *fakeAutomate) loginCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.logins
}

func (f *fakeAutomate) dataCalls() []dataCall {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]dataCall(nil), f.calls...)
}

// unauthorizedUntil rejects every request until the nth, so the first attempt
// always 401s and the replay succeeds.
func unauthorizedUntil(n int, payload string) func(*fakeAutomate, dataCall, int) (int, string) {
	return func(_ *fakeAutomate, _ dataCall, call int) (int, string) {
		if call < n {
			return http.StatusUnauthorized, `{"Message":"expired"}`
		}
		return http.StatusOK, payload
	}
}

func startFake(t *testing.T, f *fakeAutomate) *httptest.Server {
	t.Helper()

	srv := httptest.NewServer(f.handler())
	t.Cleanup(srv.Close)
	return srv
}

func testConfig(serverURL string) Config {
	return Config{ServerURL: serverURL, ClientID: "integrator-guid", Username: "user", Password: "pass"}
}

// TestUnauthorizedTriggersReloginAndReplay pins the fix for a client whose
// bearer token expires while the process keeps running: the 401 must produce a
// fresh login and a replay of the original request, not an error.
func TestUnauthorizedTriggersReloginAndReplay(t *testing.T) {
	f := &fakeAutomate{respond: unauthorizedUntil(2, `[{"Id":"1","Name":"Acme"}]`)}
	srv := startFake(t, f)

	c, err := NewClient(context.Background(), testConfig(srv.URL))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	companies, err := c.ListCompanies(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListCompanies: %v", err)
	}
	if len(companies) != 1 || companies[0].Name != "Acme" {
		t.Fatalf("companies = %+v, want one named Acme", companies)
	}

	if got := f.loginCount(); got != 2 {
		t.Errorf("logins = %d, want 2 (initial + one relogin)", got)
	}

	calls := f.dataCalls()
	if len(calls) != 2 {
		t.Fatalf("data calls = %d, want 2", len(calls))
	}
	if calls[0].auth != "Bearer token-1" {
		t.Errorf("first attempt auth = %q, want the original token", calls[0].auth)
	}
	if calls[1].auth != "Bearer token-2" {
		t.Errorf("replay auth = %q, want the reissued token", calls[1].auth)
	}
}

// TestUnauthorizedRelogsInOnlyOnce guards against a 401 that survives the
// relogin turning into an infinite login loop.
func TestUnauthorizedRelogsInOnlyOnce(t *testing.T) {
	f := &fakeAutomate{respond: func(*fakeAutomate, dataCall, int) (int, string) {
		return http.StatusUnauthorized, `{"Message":"nope"}`
	}}
	srv := startFake(t, f)

	c, err := NewClient(context.Background(), testConfig(srv.URL))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	if _, err := c.ListCompanies(context.Background(), nil); err == nil {
		t.Fatal("ListCompanies error = nil, want the 401 to surface")
	}

	if got := f.loginCount(); got != 2 {
		t.Errorf("logins = %d, want 2 (initial + exactly one relogin)", got)
	}
	if got := len(f.dataCalls()); got != 2 {
		t.Errorf("data calls = %d, want 2", got)
	}
}

// TestUnauthorizedWithoutCredentialsDoesNotRelogin covers a client built from a
// pre-issued token: there is nothing to log in with, so the 401 must surface
// untouched.
func TestUnauthorizedWithoutCredentialsDoesNotRelogin(t *testing.T) {
	f := &fakeAutomate{respond: func(*fakeAutomate, dataCall, int) (int, string) {
		return http.StatusUnauthorized, `{"Message":"nope"}`
	}}
	srv := startFake(t, f)

	c, err := NewClient(context.Background(), Config{
		ServerURL: srv.URL,
		ClientID:  "integrator-guid",
		Token:     "pre-issued",
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	if _, err := c.ListCompanies(context.Background(), nil); err == nil {
		t.Fatal("ListCompanies error = nil, want the 401 to surface")
	}

	if got := f.loginCount(); got != 0 {
		t.Errorf("logins = %d, want 0", got)
	}
	if got := len(f.dataCalls()); got != 1 {
		t.Errorf("data calls = %d, want 1", got)
	}
}

// TestUnauthorizedLoginDoesNotRecurse pins the exclusion of the apitoken
// endpoint itself: a rejected credential must fail once, not retry forever.
func TestUnauthorizedLoginDoesNotRecurse(t *testing.T) {
	f := &fakeAutomate{
		loginStatus: http.StatusUnauthorized,
		respond: func(*fakeAutomate, dataCall, int) (int, string) {
			return http.StatusOK, `[]`
		},
	}
	srv := startFake(t, f)

	if _, err := NewClient(context.Background(), testConfig(srv.URL)); err == nil {
		t.Fatal("NewClient error = nil, want the rejected login to surface")
	}

	if got := f.loginCount(); got != 1 {
		t.Errorf("logins = %d, want 1", got)
	}
}

// TestConcurrentUnauthorizedLogsInOnce pins the single-flight guard: a burst of
// calls that all discover the same expired token must share one login rather
// than issuing one apiece.
func TestConcurrentUnauthorizedLogsInOnce(t *testing.T) {
	f := &fakeAutomate{respond: func(f *fakeAutomate, call dataCall, _ int) (int, string) {
		if f.revoked != "" && call.auth == "Bearer "+f.revoked {
			return http.StatusUnauthorized, `{"Message":"expired"}`
		}
		return http.StatusOK, `[{"Id":"1","Name":"Acme"}]`
	}}
	srv := startFake(t, f)

	c, err := NewClient(context.Background(), testConfig(srv.URL))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	f.mu.Lock()
	f.revoked = f.issued
	f.mu.Unlock()

	const callers = 8
	errs := make(chan error, callers)
	var wg sync.WaitGroup
	for range callers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := c.ListCompanies(context.Background(), nil)
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			t.Fatalf("ListCompanies: %v", err)
		}
	}

	if got := f.loginCount(); got != 2 {
		t.Errorf("logins = %d, want 2 (initial + one shared relogin)", got)
	}
}

// TestReplayResendsRequestBody covers a mutating call caught by an expired
// token: the replay must carry the same body as the first attempt.
func TestReplayResendsRequestBody(t *testing.T) {
	f := &fakeAutomate{respond: unauthorizedUntil(2, `{"Id":"5","Name":"Renamed"}`)}
	srv := startFake(t, f)

	c, err := NewClient(context.Background(), testConfig(srv.URL))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ops := []PatchOp{{Op: "replace", Path: "/Name", Value: "Renamed"}}
	if _, err := c.PatchCompany(context.Background(), 5, ops); err != nil {
		t.Fatalf("PatchCompany: %v", err)
	}

	calls := f.dataCalls()
	if len(calls) != 2 {
		t.Fatalf("data calls = %d, want 2", len(calls))
	}
	if calls[0].body == "" {
		t.Fatal("first attempt sent an empty body")
	}
	if calls[0].body != calls[1].body {
		t.Errorf("replay body = %q, want the original %q", calls[1].body, calls[0].body)
	}
}

// TestNotFoundIsErrNotFoundOnEveryMethod pins the 404 sentinel across the
// helpers: Post and Patch previously returned an opaque error, so callers could
// not tell a missing record from any other failure.
func TestNotFoundIsErrNotFoundOnEveryMethod(t *testing.T) {
	f := &fakeAutomate{respond: func(*fakeAutomate, dataCall, int) (int, string) {
		return http.StatusNotFound, `{"Message":"not found"}`
	}}
	srv := startFake(t, f)

	c, err := NewClient(context.Background(), testConfig(srv.URL))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	ctx := context.Background()
	cases := []struct {
		name string
		call func() error
	}{
		{"GetCompany", func() error { _, err := c.GetCompany(ctx, 5, nil); return err }},
		{"PostCompany", func() error { _, err := c.PostCompany(ctx, &Company{Name: "Acme"}); return err }},
		{"PutCompany", func() error { _, err := c.PutCompany(ctx, 5, &Company{Name: "Acme"}); return err }},
		{"PatchCompany", func() error {
			_, err := c.PatchCompany(ctx, 5, []PatchOp{{Op: "replace", Path: "/Name", Value: "x"}})
			return err
		}},
		{"DeleteCompany", func() error { return c.DeleteCompany(ctx, 5) }},
		{"GetLocation", func() error { _, err := c.GetLocation(ctx, 5, nil); return err }},
	}

	for _, tc := range cases {
		if err := tc.call(); !errors.Is(err, ErrNotFound) {
			t.Errorf("%s error = %v, want ErrNotFound", tc.name, err)
		}
	}
}
