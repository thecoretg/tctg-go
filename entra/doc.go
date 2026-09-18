// Package entra signs users into a net/http application with Microsoft Entra
// ID and provisions them just in time.
//
// It runs the OpenID Connect authorization-code flow with PKCE against a
// single tenant, verifies the ID token, checks the tenant and any Authorizer
// you supply, hands the claims to your Provisioner, and then maintains its
// own opaque-cookie session. Microsoft is contacted only at sign-in: no
// refresh tokens are requested, no Graph calls are made, and the access token
// is discarded.
//
// # Wiring
//
// The application supplies a Provisioner that owns the user model and its
// storage, then mounts three handlers and one middleware:
//
//	type user struct {
//		ID       string
//		Name     string
//		Disabled bool
//	}
//
//	type users struct {
//		mu    sync.Mutex
//		byOID map[string]*user
//	}
//
//	func (s *users) Provision(_ context.Context, c entra.Claims) (string, error) {
//		s.mu.Lock()
//		defer s.mu.Unlock()
//		u, ok := s.byOID[c.OID]
//		if !ok {
//			u = &user{ID: c.OID}
//			s.byOID[c.OID] = u
//		}
//		if u.Disabled {
//			return "", entra.ErrUserDisabled
//		}
//		u.Name = c.Name
//		return u.ID, nil
//	}
//
//	func (s *users) Lookup(_ context.Context, id string) (*user, error) {
//		s.mu.Lock()
//		defer s.mu.Unlock()
//		u, ok := s.byOID[id]
//		switch {
//		case !ok:
//			return nil, entra.ErrUserNotFound
//		case u.Disabled:
//			return nil, entra.ErrUserDisabled
//		}
//		return u, nil
//	}
//
//	func main() {
//		ctx := context.Background()
//		store := &users{byOID: map[string]*user{}}
//		auth, err := entra.NewFromEnv[*user](ctx, store)
//		if err != nil {
//			log.Fatal(err)
//		}
//		log.Printf("register redirect URI %s", auth.RedirectURI())
//
//		mux := http.NewServeMux()
//		mux.Handle("GET /login", loginPage)
//		mux.Handle("GET /auth/sso/start", auth.Start())
//		mux.Handle("GET /auth/sso/callback", auth.Callback())
//		mux.Handle("POST /logout", auth.Logout())
//		mux.Handle("/", auth.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			u, _ := entra.UserFromContext[*user](r.Context())
//			fmt.Fprintf(w, "hello %s", u.Name)
//		})))
//
//		srv := &http.Server{Addr: ":8080", Handler: entra.SameOrigin(mux), ReadHeaderTimeout: 5 * time.Second}
//		log.Fatal(srv.ListenAndServe())
//	}
//
// # Provisioning rules
//
// Link local records on Claims.OID and nothing else. Microsoft documents
// preferred_username and email as mutable, and a username collision between
// two directories must never merge two people.
//
// Do not resurrect a record the application has disabled. Return
// ErrUserDisabled from Provision so the sign-in is refused; a plain upsert
// silently re-enables whoever was just off-boarded. Lookup runs on every
// authenticated request, so returning ErrUserDisabled there cuts a live
// session off within one request. Call Sessions().DeleteUserSessions when
// disabling an account to remove its sessions immediately.
//
// # Authorization
//
// The tenant claim is always checked. Config.Authorize adds group or role
// gating: RequireGroup, RequireRole, RequireAll and RequireAny cover the
// common cases, and any func(Claims) error works. Group removal in Entra
// only blocks the next sign-in; a running session lasts until SessionTTL
// (default 8h) unless the application ends it.
//
// # App registration checklist
//
//  1. App registrations > New registration, single tenant.
//  2. Authentication > Add a platform > Web (not SPA: a SPA registration
//     refuses the client-secret exchange with AADSTS9002327). Redirect URI
//     is exactly Auth.RedirectURI(). http://localhost is allowed for local
//     runs.
//  3. Overview: copy the Application (client) ID and the Directory (tenant)
//     ID. The tenant must be the GUID, not the domain.
//  4. Certificates & secrets > New client secret. Copy the Value. Note the
//     expiry: when it lapses every sign-in fails with ErrInvalidClient.
//  5. For RequireGroup: Token configuration > Add groups claim > "Groups
//     assigned to the application", ticked for the ID token, value "Group
//     ID". Then Enterprise applications > this app > Users and groups >
//     assign the group, or the claim is empty for everyone.
//  6. For RequireRole: App roles > Create app role, then assign it under
//     Enterprise applications > Users and groups.
//  7. API permissions: nothing to add. openid, profile and email are default
//     and need no admin consent.
//
// # Deployment notes
//
// BaseURL must be the public origin browsers use; the redirect URI is built
// from it and never from the incoming request. Behind a TLS-terminating
// proxy, cookies get the Secure flag when X-Forwarded-Proto says https.
//
// MemoryStore is the default for both sessions and flow state and is
// single-instance only. Deployments with more than one replica must supply a
// shared SessionStore and StateStore.
//
// ID token expiry is checked with no clock skew allowance; keep the host
// clock synchronised.
package entra
