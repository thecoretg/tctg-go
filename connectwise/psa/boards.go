package psa

import (
	"context"
	"fmt"
)

func boardIDEndpoint(boardID int) string {
	return fmt.Sprintf("service/boards/%d", boardID)
}

func boardIDStatusEndpoint(boardID int) string {
	return fmt.Sprintf("%s/statuses", boardIDEndpoint(boardID))
}

func boardIDStatusIDEndpoint(boardID, statusID int) string {
	return fmt.Sprintf("%s/%d", boardIDStatusEndpoint(boardID), statusID)
}

func boardIDTypeEndpoint(boardID int) string {
	return fmt.Sprintf("%s/types", boardIDEndpoint(boardID))
}

func boardIDTypeIDEndpoint(boardID, typeID int) string {
	return fmt.Sprintf("%s/%d", boardIDTypeEndpoint(boardID), typeID)
}

func boardIDSubTypeEndpoint(boardID int) string {
	return fmt.Sprintf("%s/subtypes", boardIDEndpoint(boardID))
}

func boardIDSubTypeIDEndpoint(boardID, subTypeID int) string {
	return fmt.Sprintf("%s/%d", boardIDSubTypeEndpoint(boardID), subTypeID)
}

func boardIDItemEndpoint(boardID int) string {
	return fmt.Sprintf("%s/items", boardIDEndpoint(boardID))
}

func boardIDItemIDEndpoint(boardID, itemID int) string {
	return fmt.Sprintf("%s/%d", boardIDItemEndpoint(boardID), itemID)
}

func (c *Client) PostBoard(ctx context.Context, board *Board) (*Board, error) {
	return c.Post[Board](ctx, "service/boards", board)
}

func (c *Client) ListBoards(ctx context.Context, params map[string]string, opts ...ListOption) ([]Board, error) {
	return c.GetMany[Board](ctx, "service/boards", params, opts...)
}

func (c *Client) GetBoard(ctx context.Context, boardID int, params map[string]string) (*Board, error) {
	return c.Get[Board](ctx, boardIDEndpoint(boardID), params)
}

func (c *Client) PutBoard(ctx context.Context, boardID int, board *Board) (*Board, error) {
	return c.Put[Board](ctx, boardIDEndpoint(boardID), board)
}

func (c *Client) PatchBoard(ctx context.Context, boardID int, patchOps []PatchOp) (*Board, error) {
	return c.Patch[Board](ctx, boardIDEndpoint(boardID), patchOps)
}

func (c *Client) DeleteBoard(ctx context.Context, boardID int) error {
	return c.Delete(ctx, boardIDEndpoint(boardID))
}

func (c *Client) PostBoardStatus(ctx context.Context, boardStatus *BoardStatus, boardID int) (*BoardStatus, error) {
	return c.Post[BoardStatus](ctx, boardIDStatusEndpoint(boardID), boardStatus)
}

func (c *Client) ListBoardStatuses(ctx context.Context, params map[string]string, boardID int, opts ...ListOption) ([]BoardStatus, error) {
	return c.GetMany[BoardStatus](ctx, boardIDStatusEndpoint(boardID), params, opts...)
}

func (c *Client) GetBoardStatus(ctx context.Context, statusID int, params map[string]string, boardID int) (*BoardStatus, error) {
	return c.Get[BoardStatus](ctx, boardIDStatusIDEndpoint(boardID, statusID), params)
}

func (c *Client) PutBoardStatus(ctx context.Context, statusID int, boardStatus *BoardStatus, boardID int) (*BoardStatus, error) {
	return c.Put[BoardStatus](ctx, boardIDStatusIDEndpoint(boardID, statusID), boardStatus)
}

func (c *Client) PatchBoardStatus(ctx context.Context, statusID int, patchOps []PatchOp, boardID int) (*BoardStatus, error) {
	return c.Patch[BoardStatus](ctx, boardIDStatusIDEndpoint(boardID, statusID), patchOps)
}

func (c *Client) DeleteBoardStatus(ctx context.Context, statusID int, boardID int) error {
	return c.Delete(ctx, boardIDStatusIDEndpoint(boardID, statusID))
}

func (c *Client) ListBoardTypes(ctx context.Context, params map[string]string, boardID int, opts ...ListOption) ([]BoardType, error) {
	return c.GetMany[BoardType](ctx, boardIDTypeEndpoint(boardID), params, opts...)
}

func (c *Client) GetBoardType(ctx context.Context, typeID int, params map[string]string, boardID int) (*BoardType, error) {
	return c.Get[BoardType](ctx, boardIDTypeIDEndpoint(boardID, typeID), params)
}

func (c *Client) ListBoardSubTypes(ctx context.Context, params map[string]string, boardID int, opts ...ListOption) ([]BoardSubType, error) {
	return c.GetMany[BoardSubType](ctx, boardIDSubTypeEndpoint(boardID), params, opts...)
}

func (c *Client) GetBoardSubType(ctx context.Context, subTypeID int, params map[string]string, boardID int) (*BoardSubType, error) {
	return c.Get[BoardSubType](ctx, boardIDSubTypeIDEndpoint(boardID, subTypeID), params)
}

func (c *Client) ListBoardItems(ctx context.Context, params map[string]string, boardID int, opts ...ListOption) ([]BoardItem, error) {
	return c.GetMany[BoardItem](ctx, boardIDItemEndpoint(boardID), params, opts...)
}

func (c *Client) GetBoardItem(ctx context.Context, itemID int, params map[string]string, boardID int) (*BoardItem, error) {
	return c.Get[BoardItem](ctx, boardIDItemIDEndpoint(boardID, itemID), params)
}
