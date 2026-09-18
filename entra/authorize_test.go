package entra

import (
	"errors"
	"strings"
	"testing"
)

func TestRequireGroup(t *testing.T) {
	const want = "11111111-1111-1111-1111-111111111111"
	tests := []struct {
		name    string
		claims  Claims
		wantMsg string // "" means allowed
	}{
		{"member", Claims{Groups: []string{"x", want}}, ""},
		{"non-member", Claims{Groups: []string{"x"}}, "not a member"},
		{"claim present but empty", Claims{Groups: []string{}}, "not granted this app access to any groups"},
		{"claim absent", Claims{}, "groups claim needs configuring"},
		{"claim absent with overage", Claims{Overage: true}, "too many groups"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RequireGroup(want)(tt.claims)
			if tt.wantMsg == "" {
				if err != nil {
					t.Fatalf("refused: %v", err)
				}
				return
			}
			if !errors.Is(err, ErrNotAuthorized) {
				t.Fatalf("error %v does not match ErrNotAuthorized", err)
			}
			denied, ok := errors.AsType[*DeniedError](err)
			if !ok || !strings.Contains(denied.Message, tt.wantMsg) {
				t.Fatalf("message %q does not contain %q", err, tt.wantMsg)
			}
			if strings.Contains(denied.Message, want) {
				t.Error("user-facing message leaks the group ID")
			}
		})
	}
}

func TestRequireRole(t *testing.T) {
	if err := RequireRole("Admin")(Claims{Roles: []string{"Reader", "Admin"}}); err != nil {
		t.Errorf("holder refused: %v", err)
	}
	if err := RequireRole("Admin")(Claims{Roles: []string{"Reader"}}); !errors.Is(err, ErrNotAuthorized) {
		t.Errorf("non-holder allowed: %v", err)
	}
	if err := RequireRole("Admin")(Claims{}); !errors.Is(err, ErrNotAuthorized) {
		t.Errorf("no roles allowed: %v", err)
	}
}

func TestRequireAllAny(t *testing.T) {
	allow := Authorizer(func(Claims) error { return nil })
	denyA := Authorizer(func(Claims) error { return &DeniedError{Message: "A"} })
	denyB := Authorizer(func(Claims) error { return &DeniedError{Message: "B"} })

	if err := RequireAll(allow, allow)(Claims{}); err != nil {
		t.Errorf("all allow: %v", err)
	}
	if err := RequireAll(allow, denyA, denyB)(Claims{}); err == nil || !strings.Contains(err.Error(), "A") {
		t.Errorf("all should return first refusal, got %v", err)
	}
	if err := RequireAny(denyA, allow)(Claims{}); err != nil {
		t.Errorf("any with one allow: %v", err)
	}
	if err := RequireAny(denyA, denyB)(Claims{}); err == nil || !strings.Contains(err.Error(), "A") {
		t.Errorf("any all refused should return first refusal, got %v", err)
	}
	if err := RequireAny()(Claims{}); !errors.Is(err, ErrNotAuthorized) {
		t.Errorf("any with no rules should refuse, got %v", err)
	}
}
