package psa

import (
	"context"
	"fmt"
)

func memberIdEndpoint(memberId int) string {
	return fmt.Sprintf("system/members/%d", memberId)
}

func (c *Client) PostMember(ctx context.Context, member *Member) (*Member, error) {
	return Post[Member](ctx, c, "system/members", member)
}

func (c *Client) ListMembers(ctx context.Context, params map[string]string) ([]Member, error) {
	return GetMany[Member](ctx, c, "system/members", params)
}

func (c *Client) GetMemberByIdentifier(ctx context.Context, identifier string) (*Member, error) {
	p := map[string]string{
		"conditions": fmt.Sprintf("identifier='%s'", identifier),
	}

	members, err := c.ListMembers(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("listing members: %w", err)
	}

	if len(members) == 0 {
		return nil, ErrNotFound
	}

	if len(members) > 1 {
		return nil, fmt.Errorf("expected 1 member for identifier %s, received %d", identifier, len(members))
	}

	return &members[0], nil
}

func (c *Client) GetMember(ctx context.Context, memberID int, params map[string]string) (*Member, error) {
	return GetOne[Member](ctx, c, memberIdEndpoint(memberID), params)
}

func (c *Client) PutMember(ctx context.Context, memberID int, member *Member) (*Member, error) {
	return Put[Member](ctx, c, memberIdEndpoint(memberID), member)
}

func (c *Client) PatchMember(ctx context.Context, memberID int, patchOps []PatchOp) (*Member, error) {
	return Patch[Member](ctx, c, memberIdEndpoint(memberID), patchOps)
}

func (c *Client) DeleteMember(ctx context.Context, memberID int) error {
	return Delete(ctx, c, memberIdEndpoint(memberID))
}
