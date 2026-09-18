package entra

import (
	"context"
	"encoding/json/v2"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ctxKey int

const (
	userKey ctxKey = iota
	sessionKey
)

// UserFromContext returns the user RequireAuth loaded for this request.
func UserFromContext[U any](ctx context.Context) (U, bool) {
	u, ok := ctx.Value(userKey).(U)
	return u, ok
}

// SessionFromContext returns the session RequireAuth loaded for this request.
func SessionFromContext(ctx context.Context) (Session, bool) {
	s, ok := ctx.Value(sessionKey).(Session)
	return s, ok
}

// currentSession loads and validates the session named by the request's
// cookie, returning it with the raw cookie value. It deletes an expired
// session. It does not consult the Provisioner.
func (a *Auth[U]) currentSession(r *http.Request) (s Session, token string, ok bool) {
	token = readCookie(r, a.cfg.SessionCookie)
	if token == "" {
		return Session{}, "", false
	}
	key := hashToken(token)
	s, ok, err := a.cfg.Sessions.GetSession(r.Context(), key)
	if err != nil {
		a.cfg.Logger.ErrorContext(r.Context(), "entra: get session", "err", err)
		return Session{}, "", false
	}
	if !ok {
		return Session{}, "", false
	}
	if !s.ExpiresAt.After(time.Now()) {
		_ = a.cfg.Sessions.DeleteSession(r.Context(), key)
		return Session{}, "", false
	}
	return s, token, true
}

// RequireAuth admits only requests with a live session whose user the
// Provisioner still recognises. It sets Cache-Control: no-store so protected
// pages never survive in the back/forward cache, slides the session forward
// once an eighth of SessionTTL has passed since the last slide, and places
// the user and session in the request context for UserFromContext and
// SessionFromContext. Rejected requests go to Config.Unauthorized.
func (a *Auth[U]) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		ctx := r.Context()

		s, token, ok := a.currentSession(r)
		if !ok {
			a.cfg.Unauthorized(w, r)
			return
		}
		key := hashToken(token)

		user, err := a.prov.Lookup(ctx, s.UserID)
		if err != nil {
			if errors.Is(err, ErrUserDisabled) || errors.Is(err, ErrUserNotFound) {
				_ = a.cfg.Sessions.DeleteSession(ctx, key)
				http.SetCookie(w, expire(a.sessionCookie(r, "")))
				a.cfg.Unauthorized(w, r)
				return
			}
			a.cfg.Logger.ErrorContext(ctx, "entra: lookup user", "user_id", s.UserID, "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		now := time.Now()
		if now.Sub(s.LastSeenAt) > a.cfg.SessionTTL/8 {
			s.LastSeenAt = now
			s.ExpiresAt = now.Add(a.cfg.SessionTTL)
			if err := a.cfg.Sessions.TouchSession(ctx, key, s.LastSeenAt, s.ExpiresAt); err != nil {
				a.cfg.Logger.ErrorContext(ctx, "entra: touch session", "err", err)
			} else {
				// The cookie's own MaxAge must slide with the row or the
				// browser drops it while the server would still honour it.
				http.SetCookie(w, a.sessionCookie(r, token))
			}
		}

		ctx = context.WithValue(ctx, userKey, user)
		ctx = context.WithValue(ctx, sessionKey, s)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// defaultUnauthorized answers JSON to API-shaped requests and redirects
// browsers to LoginPath with ?next set to the requested path.
func (a *Auth[U]) defaultUnauthorized(w http.ResponseWriter, r *http.Request) {
	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.MarshalWrite(w, map[string]string{"error": "unauthenticated", "login": a.cfg.LoginPath})
		return
	}

	u, err := url.Parse(a.cfg.LoginPath)
	if err != nil {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}
	if next := safeNext(r.URL.RequestURI()); next != "" && r.Method == http.MethodGet {
		q := u.Query()
		q.Set("next", next)
		u.RawQuery = q.Encode()
	}
	http.Redirect(w, r, u.String(), http.StatusFound)
}

// wantsJSON reports whether a request came from script rather than a
// navigation: a Sec-Fetch-Dest other than document, or an Accept header that
// names neither text/html nor */*.
func wantsJSON(r *http.Request) bool {
	if d := r.Header.Get("Sec-Fetch-Dest"); d != "" && d != "document" {
		return true
	}
	accept := r.Header.Get("Accept")
	if accept == "" {
		return false
	}
	return !strings.Contains(accept, "text/html") && !strings.Contains(accept, "*/*")
}

// SameOrigin rejects POST, PUT, PATCH and DELETE requests whose Origin header
// names a different host from the one the request arrived at. A missing
// Origin is allowed, since non-browser clients carry no ambient cookies.
// Together with SameSite=Lax cookies this is the package's CSRF defence; wrap
// the whole mux with it.
func SameOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		default:
			next.ServeHTTP(w, r)
			return
		}

		origin := r.Header.Get("Origin")
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}
		u, err := url.Parse(origin)
		if err != nil || u.Host == "" || !strings.EqualFold(u.Host, r.Host) {
			http.Error(w, "cross-origin request rejected", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
