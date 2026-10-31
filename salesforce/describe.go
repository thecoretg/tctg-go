package salesforce

import (
	"context"
	"fmt"
	"net/url"
)

// Field is one field of an sObject as Describe reports it. Type is
// Salesforce's field type, such as "id", "string", "reference", "picklist",
// "double", "boolean", "date" or "datetime".
type Field struct {
	Name             string          `json:"name"`
	Label            string          `json:"label"`
	Type             string          `json:"type"`
	SoapType         string          `json:"soapType"`
	Length           int             `json:"length"`
	Precision        int             `json:"precision"`
	Scale            int             `json:"scale"`
	Nillable         bool            `json:"nillable"`
	Custom           bool            `json:"custom"`
	Calculated       bool            `json:"calculated"`
	Createable       bool            `json:"createable"`
	Updateable       bool            `json:"updateable"`
	Filterable       bool            `json:"filterable"`
	Unique           bool            `json:"unique"`
	ExternalID       bool            `json:"externalId"`
	ReferenceTo      []string        `json:"referenceTo"`
	RelationshipName string          `json:"relationshipName"`
	PicklistValues   []PicklistValue `json:"picklistValues"`
}

// PicklistValue is one allowed value of a picklist field.
type PicklistValue struct {
	Value        string `json:"value"`
	Label        string `json:"label"`
	Active       bool   `json:"active"`
	DefaultValue bool   `json:"defaultValue"`
}

// Describe returns the fields of the sObject named object, such as "Account"
// or "Custom_Object__c", in the order Salesforce lists them. An unknown object
// returns ErrNotFound.
func (c *Client) Describe(ctx context.Context, object string) ([]Field, error) {
	result, err := c.Get[struct {
		Fields []Field `json:"fields"`
	}](ctx, c.endpointURL("sobjects/"+url.PathEscape(object)+"/describe"), nil)
	if err != nil {
		return nil, fmt.Errorf("describe %s: %w", object, err)
	}
	return result.Fields, nil
}
