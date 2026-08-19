package psa

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func newLiveClient(t *testing.T) *Client {
	t.Helper()
	_ = godotenv.Load("../../.env")
	c, err := NewClient(context.Background(), Config{
		PublicKey:  os.Getenv("CW_PUB_KEY"),
		PrivateKey: os.Getenv("CW_PRIV_KEY"),
		ClientID:   os.Getenv("CW_CLIENT_ID"),
		CompanyID:  os.Getenv("CW_COMPANY_ID"),
	})
	if err != nil {
		t.Skip("skipping integration test:", err)
	}
	return c
}

// TestListCompaniesLimitAuth caps a live call. It must stay capped: an
// uncapped ListCompanies walks the whole tenant.
func TestListCompaniesLimitAuth(t *testing.T) {
	c := newLiveClient(t)

	start := time.Now()
	companies, err := c.ListCompanies(context.Background(), nil, WithLimit(5))
	if err != nil {
		t.Fatalf("ListCompanies: %v", err)
	}

	if len(companies) > 5 {
		t.Fatalf("companies = %d, want at most 5", len(companies))
	}
	t.Logf("auth OK; got %d companies in %s", len(companies), time.Since(start).Round(time.Millisecond))
	for _, co := range companies {
		t.Logf("  - %s (id %d, identifier %s)", co.Name, co.ID, co.Identifier)
	}
}

// TestListCompaniesConditionsAuth confirms ConnectWise applies a conditions
// filter passed through params.
func TestListCompaniesConditionsAuth(t *testing.T) {
	c := newLiveClient(t)

	params := map[string]string{"conditions": `status/name="Active" and deletedFlag=false`}
	companies, err := c.ListCompanies(context.Background(), params, WithLimit(5))
	if err != nil {
		t.Fatalf("ListCompanies: %v", err)
	}

	t.Logf("got %d active companies", len(companies))
	for _, co := range companies {
		if co.DeletedFlag {
			t.Errorf("%s came back deleted despite the filter", co.Name)
		}
		if co.Status.Name != "Active" {
			t.Errorf("%s status = %q, want Active", co.Name, co.Status.Name)
		}
	}
}
