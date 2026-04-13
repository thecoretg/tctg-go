package rewst

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

type (
	BulkGetOrgVarMapInput struct {
		VarNames       []string
		ExcludedOrgIDs []string
	}

	GetOrgVarMapInput struct {
		VarName        string   `json:"var_name"`
		ExcludedOrgIDs []string `json:"excluded_org_ids"`
	}

	OrgVarMapResp struct {
		*OrgVarMap
		Error string `json:"error"`
	}

	OrgVarMap struct {
		Map        map[string]string `json:"map"`
		ReverseMap map[string]string `json:"reverse_map"`
	}

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
	result, err := Post[UpsertOrgVarResp](ctx, c.wc, c.upsertVarURL, input)
	if err != nil {
		return fmt.Errorf("upsert org var: %w", err)
	}

	if result.Error != "" {
		return fmt.Errorf("upsert org var: %w", errors.New(result.Error))
	}

	return nil
}

func (c *Client) BulkGetOrgVarMap(ctx context.Context, input BulkGetOrgVarMapInput) (map[string]OrgVarMap, error) {
	data := make(map[string]OrgVarMap, len(input.VarNames))
	g, ctx := errgroup.WithContext(ctx)

	for _, n := range input.VarNames {
		g.Go(func() error {
			in := GetOrgVarMapInput{
				VarName:        n,
				ExcludedOrgIDs: input.ExcludedOrgIDs,
			}
			m, err := c.GetOrgVarMap(ctx, in)
			if err != nil {
				return fmt.Errorf("getting org var map: %w", err)
			}
			data[n] = *m
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Client) GetOrgVarMap(ctx context.Context, input GetOrgVarMapInput) (*OrgVarMap, error) {
	result, err := Post[OrgVarMapResp](ctx, c.wc, c.getOrgVarMapURL, input)
	if err != nil {
		return nil, fmt.Errorf("get org var map: %w", err)
	}

	if result.Error != "" {
		return nil, fmt.Errorf("get org var map: %w", errors.New(result.Error))
	}

	return result.OrgVarMap, nil
}
