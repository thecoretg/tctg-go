package entra

import "slices"

// Authorizer decides whether a verified, tenant-checked sign-in may proceed to
// provisioning. Return nil to allow. Return a *DeniedError (or any error
// wrapping ErrNotAuthorized) to refuse with a message the user may see; any
// other error is treated as an internal failure.
type Authorizer func(Claims) error

// RequireGroup allows only members of the Entra group with the given object
// ID. It distinguishes the three misconfigurations that also present as
// "not a member" so the message points at the fix:
//
//   - groups claim absent and the user is in overage
//   - groups claim absent (not configured on the app registration)
//   - groups claim present but empty (no groups assigned to the enterprise app)
func RequireGroup(groupID string) Authorizer {
	return func(c Claims) error {
		switch {
		case c.Groups == nil && c.Overage:
			return &DeniedError{Message: "Microsoft could not list your group memberships because too many groups are assigned to this app."}
		case c.Groups == nil:
			return &DeniedError{Message: "Microsoft did not include group information in this sign-in. The app registration's groups claim needs configuring."}
		case len(c.Groups) == 0:
			return &DeniedError{Message: "Your organisation has not granted this app access to any groups yet."}
		case !slices.Contains(c.Groups, groupID):
			return &DeniedError{Message: "You are not a member of the group that grants access to this application."}
		}
		return nil
	}
}

// RequireRole allows only users assigned the named app role.
func RequireRole(role string) Authorizer {
	return func(c Claims) error {
		if !slices.Contains(c.Roles, role) {
			return &DeniedError{Message: "Your account has not been assigned the role that grants access to this application."}
		}
		return nil
	}
}

// RequireAll allows the sign-in only when every Authorizer allows it. The
// first refusal is returned.
func RequireAll(authorizers ...Authorizer) Authorizer {
	return func(c Claims) error {
		for _, a := range authorizers {
			if err := a(c); err != nil {
				return err
			}
		}
		return nil
	}
}

// RequireAny allows the sign-in when at least one Authorizer allows it. When
// all refuse, the first refusal is returned. With no Authorizers it refuses.
func RequireAny(authorizers ...Authorizer) Authorizer {
	return func(c Claims) error {
		var first error
		for _, a := range authorizers {
			err := a(c)
			if err == nil {
				return nil
			}
			if first == nil {
				first = err
			}
		}
		if first == nil {
			first = &DeniedError{Message: "No access rule allows this sign-in."}
		}
		return first
	}
}
