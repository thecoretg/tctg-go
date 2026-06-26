package automate

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

// TestTokenResultParsesZonelessTime reproduces the login failure: Automate
// returns expiration timestamps with no timezone offset, which the standard
// time.Time rejects. It needs no credentials.
func TestTokenResultParsesZonelessTime(t *testing.T) {
	body := `{"AccessToken":"abc","ExpirationDate":"2026-06-26T15:52:53","AbsoluteExpirationDate":"2026-06-27T15:52:53"}`

	var tr TokenResult
	if err := json.Unmarshal([]byte(body), &tr); err != nil {
		t.Fatalf("unmarshal TokenResult: %v", err)
	}
	want := time.Date(2026, 6, 26, 15, 52, 53, 0, time.UTC)
	if !tr.ExpirationDate.Equal(want) {
		t.Fatalf("ExpirationDate = %v, want %v", tr.ExpirationDate.Time, want)
	}
}

func newTestClient(t *testing.T) *Client {
	t.Helper()
	_ = godotenv.Load("../../.env")
	c, err := NewClient(context.Background(), Config{
		ServerURL:         os.Getenv("CWA_SERVER_URL"),
		ClientID:          os.Getenv("CWA_CLIENT_ID"),
		Username:          os.Getenv("CWA_USERNAME"),
		Password:          os.Getenv("CWA_PASSWORD"),
		TwoFactorPasscode: os.Getenv("CWA_2FA"),
		Token:             os.Getenv("CWA_TOKEN"),
	})
	if err != nil {
		t.Skip("skipping integration test:", err)
	}
	return c
}

// TestLoginAuth exercises the live login flow. NewClient performs the apitoken
// exchange, so a non-empty token confirms auth (and timestamp parsing) works.
func TestLoginAuth(t *testing.T) {
	c := newTestClient(t)
	if c.Token() == "" {
		t.Fatal("expected a bearer token after login")
	}
	t.Logf("auth OK; token acquired (len %d)", len(c.Token()))
}

// TestListCompaniesAuth performs an authenticated GET against the live API.
func TestListCompaniesAuth(t *testing.T) {
	c := newTestClient(t)
	companies, err := c.ListCompanies(context.Background(), map[string]string{
		"pageSize": "5",
	})
	if err != nil {
		t.Fatalf("ListCompanies: %v", err)
	}
	t.Logf("auth OK; got %d company/companies", len(companies))
	for _, co := range companies {
		t.Logf("  - %s (Id %s)", co.Name, co.ID)
	}
}

// TestListAllCompaniesPaging verifies the pagination helper walks every page:
// a single page (pagesize=5) must return fewer rows than ListAll, which should
// equal the Total-Count the server reports.
func TestListAllCompaniesPaging(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	onePage, err := c.ListCompanies(ctx, map[string]string{"pagesize": "5"})
	if err != nil {
		t.Fatalf("ListCompanies: %v", err)
	}
	all, err := c.ListAllCompanies(ctx, map[string]string{"pagesize": "5"})
	if err != nil {
		t.Fatalf("ListAllCompanies: %v", err)
	}
	t.Logf("single page=%d, all pages=%d", len(onePage), len(all))
	if len(all) <= len(onePage) {
		t.Fatalf("expected ListAll (%d) to exceed a single page (%d)", len(all), len(onePage))
	}
}
