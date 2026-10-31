// Package httpxtest helps tests drive a real client against an httptest
// server, so the URLs a client builds from its fixed base URL are exercised
// rather than bypassed.
package httpxtest

import (
	"net/http"
	"net/http/httptest"
	"net/url"
)

// Redirect returns a RoundTripper that sends every request to srv, keeping
// its path and query, so a test can assert on the exact path a client built.
func Redirect(srv *httptest.Server) http.RoundTripper {
	target, err := url.Parse(srv.URL)
	if err != nil {
		panic(err)
	}
	return redirect{target: target, base: srv.Client().Transport}
}

type redirect struct {
	target *url.URL
	base   http.RoundTripper
}

func (r redirect) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.URL.Scheme = r.target.Scheme
	req.URL.Host = r.target.Host
	req.Host = r.target.Host
	return r.base.RoundTrip(req)
}
