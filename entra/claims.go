package entra

import (
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
)

// Claims is the subset of ID token claims the package reads, plus the raw
// claim map for anything else the application needs.
type Claims struct {
	// OID is the user's immutable object ID in the tenant. Link local accounts
	// on this and nothing else: Microsoft documents preferred_username and
	// email as mutable.
	OID string
	// TID is the tenant the user signed in from. Always equals Config.TenantID
	// by the time an Authorizer or Provisioner sees it.
	TID               string
	PreferredUsername string
	Email             string
	Name              string
	// Groups is nil when the token carried no groups claim at all, and empty
	// (non-nil) when the claim was present but listed nothing. The distinction
	// matters: nil means the app registration's groups claim is not configured
	// (or the user is in overage), empty means no groups are assigned to the
	// enterprise application.
	Groups []string
	// Roles holds app roles assigned to the user, when the registration
	// defines any.
	Roles []string
	// Overage is true when Entra omitted the groups list because the user is in
	// too many groups (hasgroups, or a _claim_names pointer to Graph).
	Overage bool
	// Raw is the full decoded claim set.
	Raw map[string]any
}

func parseClaims(idt *oidc.IDToken) (Claims, error) {
	var raw map[string]any
	if err := idt.Claims(&raw); err != nil {
		return Claims{}, fmt.Errorf("decode id token claims: %w", err)
	}

	c := Claims{
		Raw:               raw,
		OID:               claimString(raw, "oid"),
		TID:               claimString(raw, "tid"),
		PreferredUsername: claimString(raw, "preferred_username"),
		Email:             claimString(raw, "email"),
		Name:              claimString(raw, "name"),
	}
	if g, ok := raw["groups"].([]any); ok {
		c.Groups = claimStrings(g)
	}
	if r, ok := raw["roles"].([]any); ok {
		c.Roles = claimStrings(r)
	}
	_, hasClaimNames := raw["_claim_names"]
	c.Overage = raw["hasgroups"] == true || hasClaimNames

	if c.OID == "" {
		return c, ErrNoOID
	}
	return c, nil
}

func claimString(raw map[string]any, key string) string {
	s, _ := raw[key].(string)
	return s
}

// claimStrings always returns a non-nil slice so callers can tell "claim
// present but empty" from "claim absent".
func claimStrings(vals []any) []string {
	out := make([]string, 0, len(vals))
	for _, v := range vals {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}
