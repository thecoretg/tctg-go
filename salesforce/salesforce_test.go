package salesforce

import (
	"context"
	"encoding/json"
	"testing"

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

func TestQueryAccounts(t *testing.T) {
	c := newTestClient(t)

	accounts, err := c.QueryAccounts(context.Background(), QueryAccountsOpts{
		Fields: []string{"Phone", "Support_Agreement__c"},
	})
	if err != nil {
		t.Fatalf("QueryAccounts: %v", err)
	}
	t.Logf("got %d accounts", len(accounts))
}

func TestQueryAccountsWhere(t *testing.T) {
	c := newTestClient(t)

	accounts, err := c.QueryAccounts(context.Background(), QueryAccountsOpts{
		Fields: []string{"Phone", "Support_Agreement__c"},
		Where:  "Type = 'Customer'",
	})
	if err != nil {
		t.Fatalf("QueryAccounts: %v", err)
	}
	t.Logf("got %d accounts", len(accounts))
}

func TestAccountUnmarshal(t *testing.T) {
	data := []byte(`{
		"attributes": {"type": "Account"},
		"Id": "001Dp0000056v5LIAQ",
		"Name": "Acme Corp",
		"Phone": "(713) 393-6500",
		"Support_Agreement__c": "Co-Managed IT"
	}`)

	var a Account
	if err := json.Unmarshal(data, &a); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if a.ID != "001Dp0000056v5LIAQ" {
		t.Errorf("ID: got %q, want %q", a.ID, "001Dp0000056v5LIAQ")
	}
	if a.Name != "Acme Corp" {
		t.Errorf("Name: got %q, want %q", a.Name, "Acme Corp")
	}
	if a.Fields["Phone"] != "(713) 393-6500" {
		t.Errorf("Fields[Phone]: got %v, want %q", a.Fields["Phone"], "(713) 393-6500")
	}
	if a.Fields["Support_Agreement__c"] != "Co-Managed IT" {
		t.Errorf("Fields[Support_Agreement__c]: got %v, want %q", a.Fields["Support_Agreement__c"], "Co-Managed IT")
	}
	if _, ok := a.Fields["Id"]; ok {
		t.Error("Fields should not contain Id")
	}
	if _, ok := a.Fields["Name"]; ok {
		t.Error("Fields should not contain Name")
	}
}
