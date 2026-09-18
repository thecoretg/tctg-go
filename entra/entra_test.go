package entra

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json/v2"
	"errors"
	"io"
	"log/slog"
	"maps"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

const (
	testTenant = "11111111-2222-3333-4444-555555555555"
	testClient = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
	testGroup  = "99999999-8888-7777-6666-555555555555"
	testOID    = "0f0f0f0f-1e1e-2d2d-3c3c-4b4b4b4b4b4b"
)

// fakeIssuer is a minimal OIDC provider: discovery, JWKS and a token endpoint
// that signs whatever claims the test asks for.
type fakeIssuer struct {
	srv *httptest.Server
	key *rsa.PrivateKey

	mu              sync.Mutex
	claims          map[string]any // overrides merged onto the defaults; nil value deletes
	tokenErr        string         // when set, /token answers 400 with this error code
	expectChallenge string         // when set, /token rejects a code_verifier that does not hash to it
	gotRedirect     string
	gotCode         string
	gotGrant        string
}

func newFakeIssuer(t *testing.T) *fakeIssuer {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	f := &fakeIssuer{key: key, claims: map[string]any{}}
	f.srv = httptest.NewServer(f.handler())
	t.Cleanup(f.srv.Close)
	return f
}

func (f *fakeIssuer) set(k string, v any) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.claims[k] = v
}

func (f *fakeIssuer) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /.well-known/openid-configuration", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"issuer":                                f.srv.URL,
			"authorization_endpoint":                f.srv.URL + "/authorize",
			"token_endpoint":                        f.srv.URL + "/token",
			"jwks_uri":                              f.srv.URL + "/keys",
			"id_token_signing_alg_values_supported": []string{"RS256"},
			"response_types_supported":              []string{"code"},
			"subject_types_supported":               []string{"pairwise"},
		})
	})
	mux.HandleFunc("GET /keys", func(w http.ResponseWriter, _ *http.Request) {
		pub := f.key.PublicKey
		writeJSON(w, http.StatusOK, map[string]any{"keys": []map[string]any{{
			"kty": "RSA", "kid": "k1", "use": "sig", "alg": "RS256",
			"n": base64.RawURLEncoding.EncodeToString(pub.N.Bytes()),
			"e": base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pub.E)).Bytes()),
		}}})
	})
	mux.HandleFunc("POST /token", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		f.mu.Lock()
		defer f.mu.Unlock()
		f.gotGrant = r.PostForm.Get("grant_type")
		f.gotRedirect = r.PostForm.Get("redirect_uri")
		f.gotCode = r.PostForm.Get("code")

		if f.tokenErr != "" {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": f.tokenErr, "error_description": "AADSTS0000: " + f.tokenErr})
			return
		}
		if f.gotGrant == "client_credentials" {
			writeJSON(w, http.StatusOK, map[string]any{"access_token": "cc", "token_type": "Bearer", "expires_in": 3600})
			return
		}
		if f.expectChallenge != "" && oauth2.S256ChallengeFromVerifier(r.PostForm.Get("code_verifier")) != f.expectChallenge {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_grant", "error_description": "pkce mismatch"})
			return
		}

		now := time.Now()
		claims := map[string]any{
			"iss":                f.srv.URL,
			"aud":                testClient,
			"sub":                "sub-1",
			"iat":                now.Unix(),
			"exp":                now.Add(time.Hour).Unix(),
			"oid":                testOID,
			"tid":                testTenant,
			"name":               "Ada Lovelace",
			"preferred_username": "ada@example.com",
			"email":              "ada@example.com",
			"groups":             []any{testGroup},
		}
		maps.Copy(claims, f.claims)
		maps.DeleteFunc(claims, func(_ string, v any) bool { return v == nil })

		writeJSON(w, http.StatusOK, map[string]any{
			"access_token": "at", "token_type": "Bearer", "expires_in": 3600,
			"id_token": f.sign(claims),
		})
	})
	return mux
}

func (f *fakeIssuer) sign(claims map[string]any) string {
	b64 := base64.RawURLEncoding.EncodeToString
	payload, err := json.Marshal(claims)
	if err != nil {
		panic(err)
	}
	input := b64([]byte(`{"alg":"RS256","kid":"k1","typ":"JWT"}`)) + "." + b64(payload)
	sum := sha256.Sum256([]byte(input))
	sig, err := rsa.SignPKCS1v15(rand.Reader, f.key, crypto.SHA256, sum[:])
	if err != nil {
		panic(err)
	}
	return input + "." + b64(sig)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.MarshalWrite(w, v)
}

// fakeUsers is an in-memory Provisioner.
type testUser struct {
	ID   string
	Name string
}

type fakeUsers struct {
	mu         sync.Mutex
	users      map[string]*testUser
	disabled   map[string]bool
	provisions int
	lookups    int
	lookupErr  error
}

func newFakeUsers() *fakeUsers {
	return &fakeUsers{users: map[string]*testUser{}, disabled: map[string]bool{}}
}

func (s *fakeUsers) Provision(_ context.Context, c Claims) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.provisions++
	if s.disabled[c.OID] {
		return "", ErrUserDisabled
	}
	u, ok := s.users[c.OID]
	if !ok {
		u = &testUser{ID: c.OID}
		s.users[c.OID] = u
	}
	u.Name = c.Name
	return u.ID, nil
}

func (s *fakeUsers) Lookup(_ context.Context, id string) (*testUser, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lookups++
	if s.lookupErr != nil {
		return nil, s.lookupErr
	}
	u, ok := s.users[id]
	switch {
	case !ok:
		return nil, ErrUserNotFound
	case s.disabled[id]:
		return nil, ErrUserDisabled
	}
	return u, nil
}

func (s *fakeUsers) disable(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.disabled[id] = true
}

// harness is an application wired with the package, served over httptest,
// driven by a cookie-holding client that does not follow redirects.
type harness struct {
	t      *testing.T
	issuer *fakeIssuer
	users  *fakeUsers
	store  *MemoryStore
	auth   *Auth[*testUser]
	app    *httptest.Server
	client *http.Client
}

func newHarness(t *testing.T, mutate func(*Config)) *harness {
	t.Helper()
	h := &harness{t: t, issuer: newFakeIssuer(t), users: newFakeUsers(), store: NewMemoryStore()}

	// The app server must exist before New so BaseURL is known.
	mux := http.NewServeMux()
	h.app = httptest.NewServer(mux)
	t.Cleanup(h.app.Close)

	cfg := Config{
		TenantID:     testTenant,
		ClientID:     testClient,
		ClientSecret: "secret",
		BaseURL:      h.app.URL,
		Sessions:     h.store,
		States:       h.store,
		Logger:       slog.New(slog.DiscardHandler),
		issuer:       h.issuer.srv.URL,
	}
	if mutate != nil {
		mutate(&cfg)
	}
	auth, err := New(t.Context(), cfg, h.users)
	if err != nil {
		t.Fatal(err)
	}
	h.auth = auth

	mux.Handle("GET /auth/sso/start", auth.Start())
	mux.Handle("GET /auth/sso/callback", auth.Callback())
	mux.Handle("/logout", auth.Logout())
	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "login page")
	})
	mux.Handle("/me", auth.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, ok := UserFromContext[*testUser](r.Context())
		if !ok {
			http.Error(w, "no user in context", http.StatusInternalServerError)
			return
		}
		_, _ = io.WriteString(w, u.Name)
	})))

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	h.client = &http.Client{
		Jar:           jar,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	return h
}

// get issues a GET against the app and returns the response with its body
// read into a string.
func (h *harness) get(path string, headers map[string]string) (*http.Response, string) {
	h.t.Helper()
	req, err := http.NewRequestWithContext(h.t.Context(), http.MethodGet, h.app.URL+path, nil)
	if err != nil {
		h.t.Fatal(err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	res, err := h.client.Do(req)
	if err != nil {
		h.t.Fatal(err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	return res, string(body)
}

// start begins a flow and returns the authorization URL's query. It teaches
// the fake issuer the nonce and PKCE challenge so the callback can succeed.
func (h *harness) start(next string) url.Values {
	h.t.Helper()
	path := "/auth/sso/start"
	if next != "" {
		path += "?next=" + url.QueryEscape(next)
	}
	res, _ := h.get(path, nil)
	if res.StatusCode != http.StatusFound {
		h.t.Fatalf("start: status %d", res.StatusCode)
	}
	loc, err := url.Parse(res.Header.Get("Location"))
	if err != nil {
		h.t.Fatal(err)
	}
	if !strings.HasPrefix(loc.String(), h.issuer.srv.URL+"/authorize?") {
		h.t.Fatalf("start redirected to %s", loc)
	}
	q := loc.Query()
	h.issuer.set("nonce", q.Get("nonce"))
	h.issuer.mu.Lock()
	h.issuer.expectChallenge = q.Get("code_challenge")
	h.issuer.mu.Unlock()
	return q
}

func (h *harness) callback(q url.Values) *http.Response {
	h.t.Helper()
	res, _ := h.get("/auth/sso/callback?"+q.Encode(), nil)
	return res
}

// login runs the whole flow and fails the test unless it lands on target.
func (h *harness) login(next, target string) {
	h.t.Helper()
	q := h.start(next)
	res := h.callback(url.Values{"code": {"code-1"}, "state": {q.Get("state")}})
	if res.StatusCode != http.StatusFound || res.Header.Get("Location") != target {
		h.t.Fatalf("callback: status %d location %q, want 302 %q", res.StatusCode, res.Header.Get("Location"), target)
	}
}

func (h *harness) cookie(name string) *http.Cookie {
	u, _ := url.Parse(h.app.URL)
	for _, c := range h.client.Jar.Cookies(u) {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func (h *harness) sessionKey() string {
	h.t.Helper()
	c := h.cookie("entra_session")
	if c == nil {
		h.t.Fatal("no session cookie in jar")
	}
	return hashToken(c.Value)
}

// loginError extracts the ?err message from a redirect to /login.
func loginError(t *testing.T, res *http.Response) string {
	t.Helper()
	if res.StatusCode != http.StatusFound {
		t.Fatalf("status %d, want 302", res.StatusCode)
	}
	loc, err := url.Parse(res.Header.Get("Location"))
	if err != nil {
		t.Fatal(err)
	}
	if loc.Path != "/login" {
		t.Fatalf("redirected to %s, want /login", loc)
	}
	return loc.Query().Get("err")
}

func setCookieNamed(res *http.Response, name string) *http.Cookie {
	for _, c := range res.Cookies() {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func TestLoginFlow(t *testing.T) {
	h := newHarness(t, func(c *Config) { c.Authorize = RequireGroup(testGroup) })

	res, _ := h.get("/auth/sso/start?next=/dashboard?tab=1", nil)
	flow := setCookieNamed(res, "entra_flow")
	if flow == nil {
		t.Fatal("start set no flow cookie")
	}
	if flow.Path != "/auth/sso" || !flow.HttpOnly || flow.SameSite != http.SameSiteLaxMode {
		t.Errorf("flow cookie attributes: %+v", flow)
	}
	loc, _ := url.Parse(res.Header.Get("Location"))
	q := loc.Query()
	for _, k := range []string{"state", "nonce", "code_challenge"} {
		if q.Get(k) == "" {
			t.Errorf("authorize URL missing %s", k)
		}
	}
	if q.Get("code_challenge_method") != "S256" || q.Get("redirect_uri") != h.auth.RedirectURI() || q.Get("scope") != "openid profile email" {
		t.Errorf("authorize params: %v", q)
	}
	h.issuer.set("nonce", q.Get("nonce"))

	res = h.callback(url.Values{"code": {"code-1"}, "state": {q.Get("state")}})
	if res.StatusCode != http.StatusFound || res.Header.Get("Location") != "/dashboard?tab=1" {
		t.Fatalf("callback: %d %s", res.StatusCode, res.Header.Get("Location"))
	}
	sess := setCookieNamed(res, "entra_session")
	if sess == nil || !sess.HttpOnly || sess.Path != "/" || sess.MaxAge != int((8*time.Hour)/time.Second) {
		t.Fatalf("session cookie: %+v", sess)
	}
	if c := setCookieNamed(res, "entra_flow"); c == nil || c.MaxAge != -1 {
		t.Errorf("flow cookie not cleared: %+v", c)
	}
	if h.issuer.gotRedirect != h.auth.RedirectURI() || h.issuer.gotCode != "code-1" {
		t.Errorf("token request: redirect=%q code=%q", h.issuer.gotRedirect, h.issuer.gotCode)
	}
	if h.users.provisions != 1 {
		t.Errorf("provisions = %d", h.users.provisions)
	}

	res, body := h.get("/me", nil)
	if res.StatusCode != http.StatusOK || body != "Ada Lovelace" {
		t.Fatalf("/me: %d %q", res.StatusCode, body)
	}
	if res.Header.Get("Cache-Control") != "no-store" {
		t.Error("missing Cache-Control: no-store")
	}
	if setCookieNamed(res, "entra_session") != nil {
		t.Error("fresh session should not be re-sent")
	}

	// A second Start with a live session short-circuits.
	res, _ = h.get("/auth/sso/start", nil)
	if res.StatusCode != http.StatusFound || res.Header.Get("Location") != "/" {
		t.Errorf("start with session: %d %s", res.StatusCode, res.Header.Get("Location"))
	}
}

func TestCallbackFailures(t *testing.T) {
	tests := []struct {
		name    string
		cfg     func(*Config)
		arrange func(h *harness) url.Values // returns the callback query
		wantMsg string
		wantErr error
	}{
		{
			name:    "no flow cookie",
			arrange: func(h *harness) url.Values { return url.Values{"code": {"c"}, "state": {"s"}} },
			wantMsg: "cookies are blocked",
			wantErr: ErrFlowExpired,
		},
		{
			name: "unknown flow (replay)",
			arrange: func(h *harness) url.Values {
				u, _ := url.Parse(h.app.URL)
				h.client.Jar.SetCookies(u, []*http.Cookie{{Name: "entra_flow", Value: "stale", Path: "/auth/sso"}})
				return url.Values{"code": {"c"}, "state": {"s"}}
			},
			wantMsg: "cookies are blocked",
			wantErr: ErrFlowExpired,
		},
		{
			name: "expired flow",
			cfg:  func(c *Config) { c.FlowTTL = time.Nanosecond },
			arrange: func(h *harness) url.Values {
				q := h.start("")
				return url.Values{"code": {"c"}, "state": {q.Get("state")}}
			},
			wantMsg: "took too long",
			wantErr: ErrFlowExpired,
		},
		{
			name: "state mismatch",
			arrange: func(h *harness) url.Values {
				h.start("")
				return url.Values{"code": {"c"}, "state": {"forged"}}
			},
			wantMsg: "could not be verified",
			wantErr: ErrStateMismatch,
		},
		{
			name: "nonce mismatch",
			arrange: func(h *harness) url.Values {
				q := h.start("")
				h.issuer.set("nonce", "other")
				return url.Values{"code": {"c"}, "state": {q.Get("state")}}
			},
			wantMsg: "could not be verified",
			wantErr: ErrStateMismatch,
		},
		{
			name: "user cancelled",
			arrange: func(h *harness) url.Values {
				h.start("")
				return url.Values{"error": {"access_denied"}, "error_description": {"AADSTS65004"}}
			},
			wantMsg: "cancelled",
			wantErr: ErrAccessDenied,
		},
		{
			name: "expired client secret",
			arrange: func(h *harness) url.Values {
				q := h.start("")
				h.issuer.mu.Lock()
				h.issuer.tokenErr = "invalid_client"
				h.issuer.mu.Unlock()
				return url.Values{"code": {"c"}, "state": {q.Get("state")}}
			},
			wantMsg: "secrets expire",
			wantErr: ErrInvalidClient,
		},
		{
			name: "pkce verifier rejected",
			arrange: func(h *harness) url.Values {
				q := h.start("")
				h.issuer.mu.Lock()
				h.issuer.expectChallenge = "not-the-challenge"
				h.issuer.mu.Unlock()
				return url.Values{"code": {"c"}, "state": {q.Get("state")}}
			},
			wantMsg: "pkce mismatch",
		},
		{
			name: "wrong tenant",
			arrange: func(h *harness) url.Values {
				q := h.start("")
				h.issuer.set("tid", "22222222-2222-3333-4444-555555555555")
				return url.Values{"code": {"c"}, "state": {q.Get("state")}}
			},
			wantMsg: "different Microsoft directory",
			wantErr: ErrWrongTenant,
		},
		{
			name: "no oid",
			arrange: func(h *harness) url.Values {
				q := h.start("")
				h.issuer.set("oid", nil)
				return url.Values{"code": {"c"}, "state": {q.Get("state")}}
			},
			wantMsg: "usable ID token",
			wantErr: ErrNoOID,
		},
		{
			name: "not a group member",
			cfg:  func(c *Config) { c.Authorize = RequireGroup(testGroup) },
			arrange: func(h *harness) url.Values {
				q := h.start("")
				h.issuer.set("groups", []any{"33333333-2222-3333-4444-555555555555"})
				return url.Values{"code": {"c"}, "state": {q.Get("state")}}
			},
			wantMsg: "not a member",
			wantErr: ErrNotAuthorized,
		},
		{
			name: "disabled at provision",
			arrange: func(h *harness) url.Values {
				h.users.disable(testOID)
				q := h.start("")
				return url.Values{"code": {"c"}, "state": {q.Get("state")}}
			},
			wantMsg: "has been disabled",
			wantErr: ErrUserDisabled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got error
			h := newHarness(t, func(c *Config) {
				if tt.cfg != nil {
					tt.cfg(c)
				}
				c.OnError = func(w http.ResponseWriter, r *http.Request, err error) {
					got = err
					http.Redirect(w, r, "/login?err="+url.QueryEscape(Message(err)), http.StatusFound)
				}
			})
			q := tt.arrange(h)
			res := h.callback(q)

			msg := loginError(t, res)
			if !strings.Contains(msg, tt.wantMsg) {
				t.Errorf("message %q does not contain %q", msg, tt.wantMsg)
			}
			if tt.wantErr != nil && !errors.Is(got, tt.wantErr) {
				t.Errorf("error %v does not wrap %v", got, tt.wantErr)
			}
			if setCookieNamed(res, "entra_session") != nil {
				t.Error("session cookie set on failure")
			}
			if c := setCookieNamed(res, "entra_flow"); c == nil || c.MaxAge != -1 {
				t.Error("flow cookie not cleared on failure")
			}
			if h.cookie("entra_session") != nil {
				t.Error("session cookie in jar after failure")
			}
		})
	}
}

func TestRequireAuthDeny(t *testing.T) {
	h := newHarness(t, nil)

	res, _ := h.get("/me?a=1", map[string]string{"Accept": "text/html,*/*"})
	if res.StatusCode != http.StatusFound {
		t.Fatalf("browser: status %d", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); loc != "/login?next=%2Fme%3Fa%3D1" {
		t.Errorf("browser redirect: %s", loc)
	}

	for name, hdr := range map[string]map[string]string{
		"accept json":    {"Accept": "application/json"},
		"fetch dest":     {"Accept": "*/*", "Sec-Fetch-Dest": "empty"},
		"no accept, xhr": {"Sec-Fetch-Dest": "empty"},
	} {
		res, body := h.get("/me", hdr)
		if res.StatusCode != http.StatusUnauthorized {
			t.Errorf("%s: status %d", name, res.StatusCode)
		}
		var got map[string]string
		if err := json.Unmarshal([]byte(body), &got); err != nil || got["error"] != "unauthenticated" || got["login"] != "/login" {
			t.Errorf("%s: body %q", name, body)
		}
	}
}

func TestRequireAuthSession(t *testing.T) {
	t.Run("disabled mid-session", func(t *testing.T) {
		h := newHarness(t, nil)
		h.login("", "/")
		key := h.sessionKey()

		h.users.disable(testOID)
		res, _ := h.get("/me", nil)
		if res.StatusCode != http.StatusFound {
			t.Fatalf("status %d, want redirect", res.StatusCode)
		}
		if c := setCookieNamed(res, "entra_session"); c == nil || c.MaxAge != -1 {
			t.Error("session cookie not cleared")
		}
		if _, ok, _ := h.store.GetSession(t.Context(), key); ok {
			t.Error("session still in store")
		}
	})

	t.Run("sliding expiry", func(t *testing.T) {
		h := newHarness(t, nil)
		h.login("", "/")
		key := h.sessionKey()

		res, _ := h.get("/me", nil)
		if setCookieNamed(res, "entra_session") != nil {
			t.Fatal("cookie re-sent before TTL/8 elapsed")
		}

		now := time.Now()
		if err := h.store.TouchSession(t.Context(), key, now.Add(-2*time.Hour), now.Add(6*time.Hour)); err != nil {
			t.Fatal(err)
		}
		res, _ = h.get("/me", nil)
		if res.StatusCode != http.StatusOK {
			t.Fatalf("status %d", res.StatusCode)
		}
		if c := setCookieNamed(res, "entra_session"); c == nil || c.MaxAge != int((8*time.Hour)/time.Second) {
			t.Fatalf("cookie not re-sent after slide: %+v", c)
		}
		s, _, _ := h.store.GetSession(t.Context(), key)
		if s.ExpiresAt.Before(now.Add(7 * time.Hour)) {
			t.Errorf("session did not slide: expires %v", s.ExpiresAt)
		}
	})

	t.Run("expired session", func(t *testing.T) {
		h := newHarness(t, nil)
		h.login("", "/")
		key := h.sessionKey()

		now := time.Now()
		if err := h.store.TouchSession(t.Context(), key, now, now.Add(-time.Second)); err != nil {
			t.Fatal(err)
		}
		res, _ := h.get("/me", nil)
		if res.StatusCode != http.StatusFound {
			t.Fatalf("status %d, want redirect", res.StatusCode)
		}
		if _, ok, _ := h.store.GetSession(t.Context(), key); ok {
			t.Error("expired session not deleted")
		}
	})

	t.Run("lookup failure is 500", func(t *testing.T) {
		h := newHarness(t, nil)
		h.login("", "/")
		h.users.lookupErr = errors.New("db down")
		res, _ := h.get("/me", nil)
		if res.StatusCode != http.StatusInternalServerError {
			t.Fatalf("status %d", res.StatusCode)
		}
	})
}

func TestLogout(t *testing.T) {
	h := newHarness(t, func(c *Config) { c.LogoutRedirect = "/bye" })
	h.login("", "/")
	key := h.sessionKey()

	res, _ := h.get("/logout", nil)
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("GET logout: %d", res.StatusCode)
	}
	if _, ok, _ := h.store.GetSession(t.Context(), key); !ok {
		t.Fatal("GET logout deleted the session")
	}

	req, _ := http.NewRequestWithContext(t.Context(), http.MethodPost, h.app.URL+"/logout", nil)
	res, err := h.client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusFound || res.Header.Get("Location") != "/bye" {
		t.Fatalf("POST logout: %d %s", res.StatusCode, res.Header.Get("Location"))
	}
	if c := setCookieNamed(res, "entra_session"); c == nil || c.MaxAge != -1 {
		t.Error("session cookie not cleared")
	}
	if _, ok, _ := h.store.GetSession(t.Context(), key); ok {
		t.Error("session still in store")
	}
	res, _ = h.get("/me", nil)
	if res.StatusCode != http.StatusFound {
		t.Errorf("/me after logout: %d", res.StatusCode)
	}
}

func TestSecureCookieBehindProxy(t *testing.T) {
	h := newHarness(t, nil)

	for _, tt := range []struct {
		proto string
		want  bool
	}{
		{"", false},
		{"https", true},
		{"https, http", true},
		{"http, https", false},
	} {
		req := httptest.NewRequest(http.MethodGet, "/auth/sso/start", nil)
		if tt.proto != "" {
			req.Header.Set("X-Forwarded-Proto", tt.proto)
		}
		rec := httptest.NewRecorder()
		h.auth.Start().ServeHTTP(rec, req)
		c := setCookieNamed(rec.Result(), "entra_flow")
		if c == nil {
			t.Fatalf("proto %q: no flow cookie", tt.proto)
		}
		if c.Secure != tt.want {
			t.Errorf("proto %q: Secure = %v, want %v", tt.proto, c.Secure, tt.want)
		}
	}
}

func TestTestConnection(t *testing.T) {
	h := newHarness(t, nil)
	if err := h.auth.TestConnection(t.Context()); err != nil {
		t.Fatalf("TestConnection: %v", err)
	}
	if h.issuer.gotGrant != "client_credentials" {
		t.Errorf("grant_type = %q", h.issuer.gotGrant)
	}

	h.issuer.mu.Lock()
	h.issuer.tokenErr = "invalid_client"
	h.issuer.mu.Unlock()
	err := h.auth.TestConnection(t.Context())
	if !errors.Is(err, ErrInvalidClient) {
		t.Fatalf("error %v does not wrap ErrInvalidClient", err)
	}
	if e, ok := errors.AsType[*Error](err); !ok || e.Code != "invalid_client" || !strings.Contains(e.Description, "AADSTS") {
		t.Errorf("Entra error not carried: %v", err)
	}
}

func TestDiscoveryFailureIsNotCached(t *testing.T) {
	h := newHarness(t, func(c *Config) { c.issuer = "http://127.0.0.1:1/dead" })
	res, _ := h.get("/auth/sso/start", nil)
	if msg := loginError(t, res); !strings.Contains(msg, "Check the server log") {
		t.Errorf("message %q", msg)
	}
	if h.auth.provider != nil {
		t.Error("failed discovery was cached")
	}
}

func TestSafeNext(t *testing.T) {
	for in, want := range map[string]string{
		"/":                "/",
		"/dash?x=1":        "/dash?x=1",
		"":                 "",
		"//evil.com":       "",
		"/\\evil.com":      "",
		"https://evil.com": "",
		"dash":             "",
	} {
		if got := safeNext(in); got != want {
			t.Errorf("safeNext(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSameOrigin(t *testing.T) {
	ok := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) })
	h := SameOrigin(ok)

	tests := []struct {
		name   string
		method string
		origin string
		want   int
	}{
		{"get ignores origin", http.MethodGet, "https://evil.example", http.StatusNoContent},
		{"post no origin", http.MethodPost, "", http.StatusNoContent},
		{"post same origin", http.MethodPost, "http://app.example", http.StatusNoContent},
		{"post same origin different case", http.MethodPost, "http://APP.example", http.StatusNoContent},
		{"post cross origin", http.MethodPost, "https://evil.example", http.StatusForbidden},
		{"delete cross origin", http.MethodDelete, "https://evil.example", http.StatusForbidden},
		{"post null origin", http.MethodPost, "null", http.StatusForbidden},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "http://app.example/x", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != tt.want {
				t.Errorf("status %d, want %d", rec.Code, tt.want)
			}
		})
	}
}

func TestNewValidation(t *testing.T) {
	valid := Config{TenantID: testTenant, ClientID: testClient, ClientSecret: "s", BaseURL: "https://app.example.com"}
	users := newFakeUsers()

	tests := []struct {
		name    string
		mutate  func(*Config)
		wantErr string
	}{
		{"missing fields", func(c *Config) { c.ClientID = ""; c.BaseURL = "" }, "missing entra config fields: ClientID, BaseURL"},
		{"domain tenant", func(c *Config) { c.TenantID = "contoso.onmicrosoft.com" }, "directory GUID"},
		{"common tenant", func(c *Config) { c.TenantID = "common" }, "directory GUID"},
		{"base url with path", func(c *Config) { c.BaseURL = "https://app.example.com/panel" }, "origin only"},
		{"base url plain http", func(c *Config) { c.BaseURL = "http://app.example.com" }, "must use https"},
		{"base url relative", func(c *Config) { c.BaseURL = "app.example.com" }, "absolute URL"},
		{"callback path relative", func(c *Config) { c.CallbackPath = "callback" }, "must start with /"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := valid
			tt.mutate(&cfg)
			_, err := New(t.Context(), cfg, users)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error %v, want containing %q", err, tt.wantErr)
			}
		})
	}

	t.Run("nil provisioner", func(t *testing.T) {
		if _, err := New[*testUser](t.Context(), valid, nil); err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("defaults and normalisation", func(t *testing.T) {
		cfg := valid
		cfg.BaseURL = "http://localhost:8080/"
		cfg.Scopes = []string{"profile", "offline_access", "email"}
		a, err := New(t.Context(), cfg, users)
		if err != nil {
			t.Fatal(err)
		}
		if a.RedirectURI() != "http://localhost:8080/auth/sso/callback" {
			t.Errorf("RedirectURI = %s", a.RedirectURI())
		}
		if got := strings.Join(a.cfg.Scopes, " "); got != "openid profile email" {
			t.Errorf("scopes = %q", got)
		}
		if a.cfg.SessionTTL != 8*time.Hour || a.cfg.FlowTTL != 10*time.Minute || a.cfg.LoginPath != "/login" || a.cfg.LogoutRedirect != "/login" {
			t.Errorf("defaults: %+v", a.cfg)
		}
		if a.cfg.Sessions == nil || a.cfg.States == nil || a.cfg.Sessions != a.cfg.States.(SessionStore) {
			t.Error("stores should default to one shared MemoryStore")
		}
		if a.cfg.issuer != "https://login.microsoftonline.com/"+testTenant+"/v2.0" {
			t.Errorf("issuer = %s", a.cfg.issuer)
		}
	})
}

func TestNewFromEnv(t *testing.T) {
	t.Setenv("ENTRA_TENANT_ID", testTenant)
	t.Setenv("ENTRA_CLIENT_ID", testClient)
	t.Setenv("ENTRA_CLIENT_SECRET", "s")
	t.Setenv("ENTRA_BASE_URL", "https://app.example.com")

	t.Setenv("ENTRA_GROUP_ID", "Panel Users")
	if _, err := NewFromEnv(t.Context(), newFakeUsers()); err == nil || !strings.Contains(err.Error(), "object ID") {
		t.Fatalf("group name accepted: %v", err)
	}

	t.Setenv("ENTRA_GROUP_ID", testGroup)
	a, err := NewFromEnv(t.Context(), newFakeUsers())
	if err != nil {
		t.Fatal(err)
	}
	if a.cfg.Authorize == nil {
		t.Fatal("ENTRA_GROUP_ID did not install an Authorizer")
	}
	if err := a.cfg.Authorize(Claims{Groups: []string{testGroup}}); err != nil {
		t.Errorf("member refused: %v", err)
	}
	if err := a.cfg.Authorize(Claims{Groups: []string{}}); !errors.Is(err, ErrNotAuthorized) {
		t.Errorf("non-member allowed: %v", err)
	}
}

// TestLiveConnection talks to the real tenant named in ../.env and is skipped
// when the ENTRA_* variables are absent.
func TestLiveConnection(t *testing.T) {
	_ = godotenv.Load("../.env")
	a, err := NewFromEnv(t.Context(), newFakeUsers())
	if err != nil {
		t.Skip("skipping integration test:", err)
	}
	if err := a.TestConnection(t.Context()); err != nil {
		t.Fatal(err)
	}
}
