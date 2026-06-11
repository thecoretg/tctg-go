package salesforce

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func newTestClient(t *testing.T) *Client {
	t.Helper()
	_ = godotenv.Load("../.env")
	client, err := NewClient(context.Background(), Config{
		ClientID:       os.Getenv("SALESFORCE_CLIENT_ID"),
		ClientSecret:   os.Getenv("SALESFORCE_CLIENT_SECRET"),
		CompanyURLName: os.Getenv("SALESFORCE_COMPANY_URL_NAME"),
	})
	if err != nil {
		t.Skip("skipping integration test:", err)
	}
	return client
}

func TestQuery(t *testing.T) {
	c := newTestClient(t)

	records, err := Query[map[string]any](context.Background(), c,
		"SELECT Id, Name, Phone, Support_Agreement__c FROM Account",
		false,
	)
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	t.Logf("got %d records", len(records))
}

func TestQuerySimplify(t *testing.T) {
	c := newTestClient(t)

	records, err := Query[map[string]any](context.Background(), c,
		"SELECT Id, Name, Phone, Support_Agreement__c FROM Account WHERE Type = 'Customer'",
		true,
	)
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	t.Logf("got %d records", len(records))
}

func TestSimplifyRecord(t *testing.T) {
	input := map[string]any{
		"attributes":             map[string]any{"type": "Account", "url": "/services/data/v64.0/sobjects/Account/SomeID"},
		"Id":                     "SomeID",
		"Name":                   "My Client Name",
		"Endpoint_Protection__c": nil,
	}

	got := simplifyRecord(input)

	if _, ok := got["attributes"]; ok {
		t.Error("attributes should be removed")
	}
	if got["id"] != "SomeID" {
		t.Errorf("id: got %v, want %q", got["id"], "SomeID")
	}
	if got["name"] != "My Client Name" {
		t.Errorf("name: got %v, want %q", got["name"], "My Client Name")
	}
	if _, ok := got["endpoint_protection"]; !ok {
		t.Error("endpoint_protection key missing")
	}
	if _, ok := got["Endpoint_Protection__c"]; ok {
		t.Error("original Endpoint_Protection__c key should not be present")
	}
}
