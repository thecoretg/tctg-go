package salesforce

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDescribeReturnsFields(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/services/data/" + versionTag + "/sobjects/Account/describe":
			_, _ = w.Write([]byte(`{"name":"Account","fields":[
				{"name":"Id","type":"id","nillable":false},
				{"name":"OwnerId","type":"reference","referenceTo":["User"],"relationshipName":"Owner"},
				{"name":"Tier__c","type":"picklist","custom":true,"picklistValues":[{"value":"Gold","label":"Gold","active":true}]}
			]}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()
	c := newStubClient(srv)

	fields, err := c.Describe(context.Background(), "Account")
	if err != nil {
		t.Fatalf("Describe() error = %v", err)
	}
	if len(fields) != 3 || fields[1].ReferenceTo[0] != "User" || !fields[2].Custom || fields[2].PicklistValues[0].Value != "Gold" {
		t.Errorf("fields = %+v", fields)
	}

	if _, err := c.Describe(context.Background(), "Nope__c"); !errors.Is(err, ErrNotFound) {
		t.Errorf("unknown object: error = %v, want ErrNotFound", err)
	}
}
