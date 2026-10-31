// Package ratelimit defines how tctg-go clients pace their requests.
//
// Clients for APIs with a published quota (Addigy, Huntress, ThreatDown,
// Umbrella) pace themselves by default: each package's DefaultLimiter is sized
// so a client never exceeds that vendor's limit, and waits only when sending
// immediately would have drawn a 429. Pass a Limiter in the client's Config to
// change that:
//
//   - Share one limiter between clients that use the same API key, since the
//     vendor counts requests per key, not per client.
//   - Pass a *rate.Limiter with a different rate when the vendor has raised
//     your quota, or a limiter shared across processes (e.g. Redis-backed)
//     when several processes use the same key.
//   - Pass None to turn pacing off.
//
// Two kinds of limiter fit: a token bucket (*rate.Limiter) for a steady rate,
// and NewWindow's sliding log for an API that counts a sliding window, which it
// matches exactly instead of approximating with a burst of 1.
//
// A limiter only coordinates the requests that go through it. Another process
// using the same API key spends the same quota without this one knowing.
package ratelimit

import "context"

// Limiter paces outgoing requests. Wait blocks until one request may be sent,
// and returns an error if ctx ends first. *rate.Limiter from
// golang.org/x/time/rate satisfies it, and also returns an error immediately
// when the wait would outlast ctx's deadline.
type Limiter interface {
	Wait(ctx context.Context) error
}

// None is a Limiter that never waits, which turns a client's pacing off.
var None Limiter = none{}

type none struct{}

func (none) Wait(context.Context) error { return nil }
