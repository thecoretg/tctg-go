package salesforce

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type (
	Account struct {
		ID     string `json:"Id"`
		Name   string `json:"Name"`
		Fields map[string]any
	}
)

type QueryAccountsOpts struct {
	Fields []string
	Where  string
}

func (c *Client) QueryAccounts(ctx context.Context, opts QueryAccountsOpts) ([]Account, error) {
	allFields := append([]string{"Id", "Name"}, opts.Fields...)
	q := fmt.Sprintf("SELECT %s FROM Account", strings.Join(allFields, ", "))
	if opts.Where != "" {
		q += " WHERE " + opts.Where
	}

	accounts, err := Query[Account](ctx, c, q)
	if err != nil {
		return nil, fmt.Errorf("query accounts: %w", err)
	}

	return accounts, nil
}

func (a *Account) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if id, ok := raw["Id"].(string); ok {
		a.ID = id
	}

	if name, ok := raw["Name"].(string); ok {
		a.Name = name
	}

	delete(raw, "Id")
	delete(raw, "Name")
	a.Fields = raw
	return nil
}
