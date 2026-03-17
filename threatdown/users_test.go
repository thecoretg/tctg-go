
package threatdown

import (
	"context"
	"testing"
)

func TestListUsers(t *testing.T) {
	c := newTestClient(t)
	ctx := context.Background()

	users, err := c.ListUsers(ctx)
	if err != nil {
		t.Fatalf("ListUsers: %v", err)
	}
	t.Logf("got %d users", len(users))
}
