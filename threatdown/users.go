package threatdown

import (
	"context"
	"fmt"
)

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

type usersResp struct {
	Users      []User `json:"users"`
	NextCursor string `json:"next_cursor"`
}

func (c *Client) ListUsers(ctx context.Context) ([]User, error) {
	users, err := getAll(ctx, c, endpointURLV1("users"), map[string]string{"page_size": "1000"}, func(r usersResp) ([]User, string) {
		return r.Users, r.NextCursor
	})
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return users, nil
}
