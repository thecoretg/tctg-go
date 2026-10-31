package ratelimit_test

import (
	"fmt"
	"testing"
	"time"

	"golang.org/x/time/rate"

	"github.com/thecoretg/tctg-go/addigy"
	"github.com/thecoretg/tctg-go/huntress"
	"github.com/thecoretg/tctg-go/threatdown"
	"github.com/thecoretg/tctg-go/umbrella"
)

type window struct {
	max int
	per time.Duration
}

// TestDefaultLimitersStayInsideVendorLimits drives each package's
// DefaultLimiter with a sender that tries every 1ms on a simulated clock, as
// fast as the limiter allows, and checks that no window of the vendor's
// published limits ever holds more requests than it allows. It also checks
// the sustained rate stays close to the limit, so a default that is safe only
// because it is needlessly slow fails too.
func TestDefaultLimitersStayInsideVendorLimits(t *testing.T) {
	tests := []struct {
		name    string
		lim     *rate.Limiter
		run     time.Duration
		windows []window
		// minPerWindow is the fewest requests the binding window must see
		// under sustained load, i.e. how close to the limit the default runs.
		minPerWindow int
	}{
		{"addigy", addigy.DefaultLimiter(), time.Minute,
			[]window{{1000, 10 * time.Second}}, 500},
		{"huntress", huntress.DefaultLimiter(), 5 * time.Minute,
			[]window{{60, time.Minute}}, 59},
		{"threatdown", threatdown.DefaultLimiter(), 5 * time.Minute,
			// A leaky bucket: 360/min sustained. The minute window may hold the
			// initial burst of 6 on top, which the bucket's depth absorbs.
			[]window{{366, time.Minute}}, 360},
		{"umbrella", umbrella.DefaultLimiter(), 95 * time.Minute,
			[]window{{5, time.Second}, {14, time.Minute}, {350, 30 * time.Minute}}, 348},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sent := simulate(tt.lim, tt.run, time.Millisecond)
			for _, w := range tt.windows {
				got := maxInWindow(sent, w.per)
				t.Logf("%s: busiest window held %d", label(w), got)
				if got > w.max {
					t.Errorf("%s: %d requests in one window, limit %d", label(w), got, w.max)
				}
			}
			binding := tt.windows[len(tt.windows)-1]
			if got := minInWindow(sent, binding.per, tt.run); got < tt.minPerWindow {
				t.Errorf("%s: only %d requests in a window under sustained load, want >= %d",
					label(binding), got, tt.minPerWindow)
			}
		})
	}
}

func label(w window) string { return fmt.Sprintf("%d per %s", w.max, w.per) }

// simulate returns the send times a greedy sender gets from lim over run.
func simulate(lim *rate.Limiter, run, step time.Duration) []time.Time {
	start := time.Unix(0, 0)
	var sent []time.Time
	for at := start; at.Before(start.Add(run)); at = at.Add(step) {
		for lim.AllowN(at, 1) {
			sent = append(sent, at)
		}
	}
	return sent
}

// maxInWindow is the most sends in any half-open span of length per.
func maxInWindow(sent []time.Time, per time.Duration) int {
	best, lo := 0, 0
	for hi := range sent {
		for !sent[lo].Add(per).After(sent[hi]) {
			lo++
		}
		best = max(best, hi-lo+1)
	}
	return best
}

// minInWindow is the fewest sends in any span of length per that starts after
// the initial burst has drained and ends before run does.
func minInWindow(sent []time.Time, per, run time.Duration) int {
	start := sent[0]
	worst := len(sent)
	for from := start.Add(per); !from.Add(per).After(start.Add(run)); from = from.Add(time.Second) {
		n := 0
		for _, s := range sent {
			if !s.Before(from) && s.Before(from.Add(per)) {
				n++
			}
		}
		worst = min(worst, n)
	}
	return worst
}
