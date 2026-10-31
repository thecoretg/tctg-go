package threatdown

import "context"

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

// ListUsers returns every OneView user.
func (c *Client) ListUsers(ctx context.Context) ([]User, error) {
	return c.ListAll[User](ctx, "v1/users", map[string]string{"page_size": "1000"})
}
