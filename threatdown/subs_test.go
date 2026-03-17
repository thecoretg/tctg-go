package threatdown

import (
	"context"
	"testing"
)

func TestSubscriptionLifecycle(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	site := createTestSite(t, c)

	// Create
	subs, err := c.CreateSiteSubscription(ctx, site.ID, SiteSubscription{
		Product:     ProductTypeEDR,
		MachineType: MachineTypeWorkstation,
		TermType:    TermTypePaid,
		TermLength:  "0",
	})
	if err != nil {
		t.Fatalf("CreateSiteSubscription: %v", err)
	}
	t.Logf("created %d subscription(s)", len(subs))

	// Get
	got, err := c.GetSiteSubscriptions(ctx, site.ID)
	if err != nil {
		t.Fatalf("GetSiteSubscriptions: %v", err)
	}
	if len(got) == 0 {
		t.Fatal("GetSiteSubscriptions: expected at least one subscription")
	}

	// Update — bump each allocation by 5
	for i := range got {
		got[i].Allocations.EndpointProtection += 5
	}
	updated, err := c.UpdateSiteSubscriptions(ctx, site.ID, got)
	if err != nil {
		t.Fatalf("UpdateSiteSubscriptions: %v", err)
	}
	t.Logf("updated %d subscription(s)", len(updated))
}
