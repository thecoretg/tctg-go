package rewst

import (
	"context"
	"errors"
	"fmt"
)

type (
	UpsertOrgVarInput struct {
		OrgID   string `json:"org_id"`
		VarName string `json:"name"`
		Value   any    `json:"value"`
	}

	UpsertOrgVarResp struct {
		Error string `json:"error"`
	}
)

func (c *Client) UpsertOrgVar(ctx context.Context, input UpsertOrgVarInput) error {
	result, err := post[UpsertOrgVarResp](ctx, c, c.upsertVarURL, input)
	if err != nil {
		return fmt.Errorf("upsert org var: %w", err)
	}

	if result.Error != "" {
		return fmt.Errorf("upsert org var: %w", errors.New(result.Error))
	}

	return nil
}
