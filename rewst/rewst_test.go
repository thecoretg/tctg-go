package rewst

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func newTestClient(t *testing.T) *Client {
	t.Helper()
	_ = godotenv.Load("../.env")
	client, err := NewClient(Config{
		WebhookSecret:   os.Getenv("REWST_WEBHOOK_SECRET"),
		ListOrgsURL:     os.Getenv("REWST_LIST_ORGS_URL"),
		UpsertVarURL:    os.Getenv("REWST_UPSERT_ORG_VAR_URL"),
		GetOrgVarMapURL: os.Getenv("REWST_GET_ORG_VAR_MAP_URL"),
	})
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

func TestBulkGetOrgVarMap(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	varNames := strings.Split(os.Getenv("REWST_BULK_GET_VAR_NAMES"), ",")
	if len(varNames) == 0 || varNames[0] == "" {
		t.Skip("skipping: REWST_BULK_GET_VAR_NAMES must be set (comma-separated)")
	}

	result, err := c.BulkGetOrgVarMap(ctx, BulkGetOrgVarMapInput{
		VarNames: varNames,
	})
	if err != nil {
		t.Fatalf("BulkGetOrgVarMap: %v", err)
	}
	for varName, m := range result {
		t.Logf("var %q: %d entries", varName, len(m.Map))
	}
}

func TestGetOrgVarMap(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	varName := os.Getenv("REWST_GET_VAR_NAME")
	if varName == "" {
		t.Skip("skipping: REWST_GET_VAR_NAME must be set")
	}

	result, err := c.GetOrgVarMap(ctx, GetOrgVarMapInput{
		VarName: varName,
	})
	if err != nil {
		t.Fatalf("GetOrgVarMap: %v", err)
	}
	t.Logf("got %d entries in map", len(result.Map))
	for k, v := range result.Map {
		t.Logf("  %s: %s", k, v)
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
