package entra

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"path"
	"strings"
	"time"
)

// newToken returns 32 bytes of randomness, base64url encoded, for use as a
// flow or session cookie value.
func newToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// hashToken is the store key for a cookie value. Stores hold only the digest,
// so a leaked store cannot be replayed as cookies.
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// isSecure reports whether the browser reached us over TLS, either directly or
// through a proxy that set X-Forwarded-Proto. Only the first (client-most)
// value of a comma-separated header counts.
func isSecure(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	first, _, _ := strings.Cut(r.Header.Get("X-Forwarded-Proto"), ",")
	return strings.EqualFold(strings.TrimSpace(first), "https")
}

func readCookie(r *http.Request, name string) string {
	c, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return c.Value
}

func (a *Auth[U]) sessionCookie(r *http.Request, value string) *http.Cookie {
	return &http.Cookie{
		Name:     a.cfg.SessionCookie,
		Value:    value,
		Path:     "/",
		MaxAge:   int(a.cfg.SessionTTL / time.Second),
		HttpOnly: true,
		Secure:   isSecure(r),
		SameSite: http.SameSiteLaxMode,
	}
}

// flowCookiePath scopes the flow cookie to the directory holding the callback
// so it is never sent with ordinary page requests.
func (a *Auth[U]) flowCookiePath() string {
	dir := path.Dir(a.cfg.CallbackPath)
	if dir == "." || dir == "" {
		return "/"
	}
	return dir
}

func (a *Auth[U]) flowCookie(r *http.Request, value string) *http.Cookie {
	return &http.Cookie{
		Name:     a.cfg.FlowCookie,
		Value:    value,
		Path:     a.flowCookiePath(),
		MaxAge:   int(a.cfg.FlowTTL / time.Second),
		HttpOnly: true,
		Secure:   isSecure(r),
		SameSite: http.SameSiteLaxMode,
	}
}

// expire turns a cookie into a deletion. Path must match the original or the
// browser keeps the old cookie.
func expire(c *http.Cookie) *http.Cookie {
	c.Value = ""
	c.MaxAge = -1
	return c
}
