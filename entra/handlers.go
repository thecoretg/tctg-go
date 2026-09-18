package entra

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// Start begins a sign-in. Mount it on a GET route. It accepts ?next=<path>
// and carries a validated value through to the end of the flow. A browser
// that already holds a valid session is sent to SuccessPath.
func (a *Auth[U]) Start() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if _, _, ok := a.currentSession(r); ok {
			http.Redirect(w, r, a.cfg.SuccessPath, http.StatusFound)
			return
		}

		oc, _, err := a.oauthConfig(ctx)
		if err != nil {
			a.cfg.Logger.ErrorContext(ctx, "entra: cannot start sign-in", "err", err)
			a.cfg.OnError(w, r, err)
			return
		}

		flowToken, err := newToken()
		if err != nil {
			a.fail(w, r, fmt.Errorf("entra: random: %w", err))
			return
		}
		state, err := newToken()
		if err != nil {
			a.fail(w, r, fmt.Errorf("entra: random: %w", err))
			return
		}
		nonce, err := newToken()
		if err != nil {
			a.fail(w, r, fmt.Errorf("entra: random: %w", err))
			return
		}
		verifier := oauth2.GenerateVerifier()

		fs := FlowState{
			State:     state,
			Nonce:     nonce,
			Verifier:  verifier,
			Next:      safeNext(r.URL.Query().Get("next")),
			ExpiresAt: time.Now().Add(a.cfg.FlowTTL),
		}
		if err := a.cfg.States.PutState(ctx, hashToken(flowToken), fs); err != nil {
			a.fail(w, r, fmt.Errorf("entra: store flow state: %w", err))
			return
		}

		http.SetCookie(w, a.flowCookie(r, flowToken))
		// response_mode stays at the default "query": the flow cookie is
		// SameSite=Lax, which browsers send on a top-level GET navigation but
		// withhold on the cross-site POST that form_post would produce.
		http.Redirect(w, r, oc.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier), oidc.Nonce(nonce)), http.StatusFound)
	})
}

// Callback completes a sign-in. Mount it on a GET route at CallbackPath. On
// success it stores a session, sets the session cookie and redirects to the
// validated ?next from Start or to SuccessPath. On any failure the flow cookie
// is cleared and OnError is called with an error wrapping the relevant
// sentinel.
func (a *Auth[U]) Callback() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := a.cfg.Logger

		// The flow is one-shot whether or not it succeeds.
		flowToken := readCookie(r, a.cfg.FlowCookie)
		http.SetCookie(w, expire(a.flowCookie(r, "")))

		if flowToken == "" {
			a.cfg.OnError(w, r, fmt.Errorf("%w: no flow cookie", ErrFlowExpired))
			return
		}
		fs, ok, err := a.cfg.States.TakeState(ctx, hashToken(flowToken))
		if err != nil {
			a.fail(w, r, fmt.Errorf("entra: take flow state: %w", err))
			return
		}
		if !ok {
			log.WarnContext(ctx, "entra: callback with no flow state; a replay or an expired flow")
			a.cfg.OnError(w, r, fmt.Errorf("%w: unknown flow", ErrFlowExpired))
			return
		}
		if !fs.ExpiresAt.After(time.Now()) {
			a.cfg.OnError(w, r, fmt.Errorf("%w: flow past its TTL", ErrFlowExpired))
			return
		}

		q := r.URL.Query()
		if code := q.Get("error"); code != "" {
			e := &Error{Code: code, Description: q.Get("error_description")}
			if code == "access_denied" {
				e.Err = ErrAccessDenied
			} else {
				log.WarnContext(ctx, "entra: authorization error", "code", code, "description", e.Description)
			}
			a.cfg.OnError(w, r, e)
			return
		}
		if q.Get("state") != fs.State {
			log.WarnContext(ctx, "entra: state mismatch on callback")
			a.cfg.OnError(w, r, fmt.Errorf("%w: state", ErrStateMismatch))
			return
		}

		oc, verifier, err := a.oauthConfig(ctx)
		if err != nil {
			log.ErrorContext(ctx, "entra: cannot complete sign-in", "err", err)
			a.cfg.OnError(w, r, err)
			return
		}

		tok, err := oc.Exchange(a.httpContext(ctx), q.Get("code"), oauth2.VerifierOption(fs.Verifier))
		if err != nil {
			err = mapTokenError(err)
			log.ErrorContext(ctx, "entra: token exchange failed", "err", err)
			a.cfg.OnError(w, r, err)
			return
		}
		rawID, _ := tok.Extra("id_token").(string)
		if rawID == "" {
			a.fail(w, r, errors.New("entra: token response carried no id_token"))
			return
		}
		idt, err := verifier.Verify(ctx, rawID)
		if err != nil {
			log.WarnContext(ctx, "entra: id token rejected", "err", err)
			a.fail(w, r, fmt.Errorf("entra: verify id token: %w", err))
			return
		}
		if idt.Nonce != fs.Nonce {
			log.WarnContext(ctx, "entra: nonce mismatch on callback")
			a.cfg.OnError(w, r, fmt.Errorf("%w: nonce", ErrStateMismatch))
			return
		}

		claims, err := parseClaims(idt)
		if err != nil {
			log.WarnContext(ctx, "entra: unusable id token", "err", err)
			a.cfg.OnError(w, r, err)
			return
		}
		if !strings.EqualFold(claims.TID, a.cfg.TenantID) {
			// The tenant goes to the log, not the browser: the callback is
			// reachable by anyone.
			log.WarnContext(ctx, "entra: refused sign-in from another tenant", "tid", claims.TID, "oid", claims.OID)
			a.cfg.OnError(w, r, ErrWrongTenant)
			return
		}

		if a.cfg.Authorize != nil {
			if err := a.cfg.Authorize(claims); err != nil {
				if !errors.Is(err, ErrNotAuthorized) {
					a.fail(w, r, fmt.Errorf("entra: authorize: %w", err))
					return
				}
				log.WarnContext(ctx, "entra: refused sign-in",
					"user", claims.PreferredUsername, "oid", claims.OID,
					"groups", claims.Groups, "roles", claims.Roles, "overage", claims.Overage, "err", err)
				a.cfg.OnError(w, r, err)
				return
			}
		}

		userID, err := a.prov.Provision(ctx, claims)
		if err != nil {
			if errors.Is(err, ErrUserDisabled) {
				log.WarnContext(ctx, "entra: refused disabled account", "oid", claims.OID)
				a.cfg.OnError(w, r, err)
				return
			}
			a.fail(w, r, fmt.Errorf("entra: provision: %w", err))
			return
		}

		sessionToken, err := newToken()
		if err != nil {
			a.fail(w, r, fmt.Errorf("entra: random: %w", err))
			return
		}
		now := time.Now()
		s := Session{UserID: userID, CreatedAt: now, LastSeenAt: now, ExpiresAt: now.Add(a.cfg.SessionTTL)}
		if err := a.cfg.Sessions.PutSession(ctx, hashToken(sessionToken), s); err != nil {
			a.fail(w, r, fmt.Errorf("entra: store session: %w", err))
			return
		}
		http.SetCookie(w, a.sessionCookie(r, sessionToken))
		log.InfoContext(ctx, "entra: signed in", "oid", claims.OID, "user", claims.PreferredUsername)

		// Validated again on purpose: a tampered store cannot become an open
		// redirect.
		target := safeNext(fs.Next)
		if target == "" {
			target = a.cfg.SuccessPath
		}
		http.Redirect(w, r, target, http.StatusFound)
	})
}

// Logout deletes the session server-side, clears the cookie and redirects to
// LogoutRedirect. It answers 405 to anything but POST: a GET logout can be
// fired by link prefetching or an <img> tag. The Microsoft session is left
// alone, so the next Start signs the same user in again without a prompt.
func (a *Auth[U]) Logout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if tok := readCookie(r, a.cfg.SessionCookie); tok != "" {
			if err := a.cfg.Sessions.DeleteSession(r.Context(), hashToken(tok)); err != nil {
				a.cfg.Logger.ErrorContext(r.Context(), "entra: delete session", "err", err)
			}
		}
		http.SetCookie(w, expire(a.sessionCookie(r, "")))
		http.Redirect(w, r, a.cfg.LogoutRedirect, http.StatusFound)
	})
}

// fail handles an internal error: log it and hand a generic error to OnError.
func (a *Auth[U]) fail(w http.ResponseWriter, r *http.Request, err error) {
	a.cfg.Logger.ErrorContext(r.Context(), "entra: sign-in failed", "err", err)
	a.cfg.OnError(w, r, err)
}

// defaultOnError redirects to LoginPath with a user-safe message in ?err.
func (a *Auth[U]) defaultOnError(w http.ResponseWriter, r *http.Request, err error) {
	u, perr := url.Parse(a.cfg.LoginPath)
	if perr != nil {
		http.Error(w, Message(err), http.StatusBadRequest)
		return
	}
	q := u.Query()
	q.Set("err", Message(err))
	u.RawQuery = q.Encode()
	http.Redirect(w, r, u.String(), http.StatusFound)
}

// Message returns a sentence suitable for showing to the person who tried to
// sign in. It never includes tenant, group or user identifiers. Custom
// OnError handlers can use it to keep the same wording.
func Message(err error) string {
	var denied *DeniedError
	switch {
	case errors.As(err, &denied):
		return denied.Message
	case errors.Is(err, ErrAccessDenied):
		return "Sign-in was cancelled."
	case errors.Is(err, ErrInvalidClient):
		return "Microsoft rejected this app's credentials. Entra client secrets expire; check the secret's expiry in the app registration."
	case errors.Is(err, ErrWrongTenant):
		return "That sign-in came from a different Microsoft directory."
	case errors.Is(err, ErrUserDisabled):
		return "This account has been disabled."
	case errors.Is(err, ErrFlowExpired):
		return "Sign-in took too long, or cookies are blocked. Try signing in again."
	case errors.Is(err, ErrStateMismatch):
		return "Sign-in could not be verified. Try signing in again."
	case errors.Is(err, ErrNoOID):
		return "Microsoft did not return a usable ID token."
	}
	if e, ok := errors.AsType[*Error](err); ok && e.Description != "" {
		return "Microsoft rejected the sign-in: " + e.Description
	}
	return "Could not complete sign-in. Check the server log."
}

// safeNext accepts only a same-site path: it must start with a single slash.
// Anything else, including protocol-relative //host, returns "".
func safeNext(v string) string {
	if !strings.HasPrefix(v, "/") || strings.HasPrefix(v, "//") || strings.HasPrefix(v, "/\\") {
		return ""
	}
	return v
}
