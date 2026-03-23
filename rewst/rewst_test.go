package rewst

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func newTestClient(t *testing.T) *Client {
	t.Helper()
	_ = godotenv.Load("../.env")
	client, err := NewClient()
	if err != nil {
		t.Skip("skipping integration test:", err)
	}
	return client
}

func TestUpsertOrgVar(t *testing.T) {
	c := newTestClient(t)

	orgID := os.Getenv("REWST_UPSERT_ORG_ID")
	varName := os.Getenv("REWST_UPSERT_VAR_NAME")
	value := os.Getenv("REWST_UPSERT_VALUE")

	if orgID == "" || varName == "" || value == "" {
		t.Skip("skipping: REWST_UPSERT_ORG_ID, REWST_UPSERT_VAR_NAME, and REWST_UPSERT_VALUE must be set")
	}

	err := c.UpsertOrgVar(context.Background(), UpsertOrgVarInput{
		OrgID:   orgID,
		VarName: varName,
		Value:   value,
	})
	if err != nil {
		t.Fatalf("UpsertOrgVar: %v", err)
	}
}

func TestListOrgs(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	orgs, err := c.ListOrgs(ctx)
	if err != nil {
		t.Fatalf("ListOrgs: %v", err)
	}
	t.Logf("got %d orgs", len(orgs))
	for _, o := range orgs {
		t.Logf("  %s: %s", o.ID, o.Name)
	}
}
