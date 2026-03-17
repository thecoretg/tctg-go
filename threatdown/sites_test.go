package threatdown

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func newTestClient(t *testing.T) *Client {
	t.Helper()
	_ = godotenv.Load("../.env")
	client, err := NewClient(context.Background())
	if err != nil {
		t.Skip("skipping integration test:", err)
	}
	return client
}

func createTestSite(t *testing.T, c *Client) *Site {
	t.Helper()
	ctx := context.Background()
	name := fmt.Sprintf("_test-tctg-go-libs-%d", time.Now().UnixMilli())
	site, err := c.CreateSite(ctx, SiteInput{
		CompanyName: name,
		FirstName:   "Test",
		LastName:    "User",
		Email:       "test@example.com",
	})
	if err != nil {
		t.Fatalf("createTestSite: %v", err)
	}
	t.Logf("createTestSite: site created with ID %s", site.ID)
	t.Cleanup(func() {
		if err := c.DeleteSite(ctx, site.ID); err != nil {
			t.Errorf("cleanup DeleteSite %s: %v", site.ID, err)
		}
		t.Logf("cleanup complete; site %s deleted", site.ID)
	})
	return site
}

func TestListSites(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	sites, err := c.ListSites(ctx)
	if err != nil {
		t.Fatalf("ListSites: %v", err)
	}
	t.Logf("got %d sites", len(sites))
}

func TestSiteLifecycle(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	// Create
	site := createTestSite(t, c)
	t.Logf("created site %s: %s", site.ID, site.CompanyName)

	// Get
	got, err := c.GetSite(ctx, site.ID)
	if err != nil {
		t.Fatalf("GetSite: %v", err)
	}
	if got.CompanyName != site.CompanyName {
		t.Errorf("GetSite company name: got %q, want %q", got.CompanyName, site.CompanyName)
	}

	// Update
	updated, err := c.UpdateSite(ctx, site.ID, SiteInput{CompanyName: site.CompanyName + "-updated"})
	if err != nil {
		t.Fatalf("UpdateSite: %v", err)
	}
	if updated.CompanyName != site.CompanyName+"-updated" {
		t.Errorf("UpdateSite company name: got %q, want %q", updated.CompanyName, site.CompanyName+"-updated")
	}

	// GetSiteByNebulaAccountID — only testable if the site has a linked nebula account
	if site.NebulaAccountID != "" {
		byNebula, err := c.GetSiteByNebulaAccountID(ctx, site.NebulaAccountID)
		if err != nil {
			t.Fatalf("GetSiteByNebulaAccountID: %v", err)
		}
		if byNebula.ID != site.ID {
			t.Errorf("GetSiteByNebulaAccountID: got site %q, want %q", byNebula.ID, site.ID)
		}
	} else {
		t.Log("skipping GetSiteByNebulaAccountID: site has no nebula account ID")
	}
}
