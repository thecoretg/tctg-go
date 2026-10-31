package ratelimit

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Window is a sliding-log Limiter: it allows at most limit requests in any
// span of length per, and lets them all go at once when that much quota is
// free. It suits APIs that count a sliding window, such as Huntress, where a
// token bucket would have to send one request every per/limit to stay safe.
//
// It remembers the send time of the last limit requests, so it uses memory in
// proportion to limit.
type Window struct {
	limit int
	per   time.Duration

	mu   sync.Mutex
	sent []time.Time // send times in order, at most limit of them
	next int         // index of the oldest send once sent is full
}

// NewWindow returns a Window allowing limit requests in any span of length
// per. It panics if limit < 1 or per <= 0.
func NewWindow(limit int, per time.Duration) *Window {
	if limit < 1 || per <= 0 {
		panic("ratelimit: NewWindow needs limit >= 1 and per > 0")
	}
	return &Window{limit: limit, per: per, sent: make([]time.Time, 0, limit)}
}

// Wait blocks until a request may be sent and records it. Like
// rate.Limiter.Wait, it returns an error immediately when the wait would
// outlast ctx's deadline, and ctx.Err() if ctx ends first.
func (w *Window) Wait(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		now := time.Now()
		wait := w.reserve(now)
		if wait <= 0 {
			return nil
		}
		if dl, ok := ctx.Deadline(); ok && now.Add(wait).After(dl) {
			return fmt.Errorf("ratelimit: waiting %s for a free slot would exceed the context deadline", wait.Round(time.Millisecond))
		}
		timer := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
			// Another waiter may have taken the slot; check again.
		}
	}
}

// AllowAt reports whether a request may be sent at t and, if so, records it.
// It never blocks. Successive calls must not pass a t earlier than a previous
// one; it exists so a simulated clock can drive the limiter in tests.
func (w *Window) AllowAt(t time.Time) bool {
	return w.reserve(t) <= 0
}

// reserve records a send at now and returns 0 when the window has room, or
// returns how long until the oldest send ages out of it.
func (w *Window) reserve(now time.Time) time.Duration {
	w.mu.Lock()
	defer w.mu.Unlock()

	if len(w.sent) < w.limit {
		w.sent = append(w.sent, now)
		return 0
	}
	if wait := w.sent[w.next].Add(w.per).Sub(now); wait > 0 {
		return wait
	}
	w.sent[w.next] = now
	w.next = (w.next + 1) % w.limit
	return 0
}
