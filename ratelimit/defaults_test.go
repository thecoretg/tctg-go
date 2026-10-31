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
// published limits ever holds more requests than it allows. It also checks the
// opening burst and the average rate, so a default that is safe only because
// it is needlessly slow fails too.
func TestDefaultLimitersStayInsideVendorLimits(t *testing.T) {
	tests := []struct {
		name    string
		allow   func(time.Time) bool
		run     time.Duration
		windows []window
		// wantBurst is how many requests go out at once from a cold start.
		wantBurst int
		// minAvg is the fewest requests the last window must see on average
		// under sustained load, i.e. how close to the limit the default runs.
		minAvg int
	}{
		{"addigy", bucket(addigy.DefaultLimiter()), time.Minute,
			[]window{{1000, 10 * time.Second}}, 50, 500},
		{"huntress", huntress.DefaultLimiter().AllowAt, 5 * time.Minute,
			[]window{{60, time.Minute}}, 60, 59},
		{"threatdown", bucket(threatdown.DefaultLimiter()), 5 * time.Minute,
			// A leaky bucket: 360/min sustained. The minute window may hold the
			// initial burst of 6 on top, which the bucket's depth absorbs.
			[]window{{366, time.Minute}}, 6, 360},
		{"umbrella", bucket(umbrella.DefaultLimiter()), 95 * time.Minute,
			[]window{{5, time.Second}, {14, time.Minute}, {350, 30 * time.Minute}}, 2, 348},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sent := simulate(tt.allow, tt.run, time.Millisecond)
			for _, w := range tt.windows {
				got := maxInWindow(sent, w.per)
				t.Logf("%s: busiest window held %d", label(w), got)
				if got > w.max {
					t.Errorf("%s: %d requests in one window, limit %d", label(w), got, w.max)
				}
			}
			if got := burst(sent); got != tt.wantBurst {
				t.Errorf("opening burst = %d, want %d", got, tt.wantBurst)
			}
			last := tt.windows[len(tt.windows)-1]
			if got := len(sent) * int(last.per) / int(tt.run); got < tt.minAvg {
				t.Errorf("%s: averaged %d requests per window under sustained load, want >= %d",
					label(last), got, tt.minAvg)
			}
		})
	}
}

func label(w window) string { return fmt.Sprintf("%d per %s", w.max, w.per) }

// bucket adapts a token bucket to the simulated clock.
func bucket(l *rate.Limiter) func(time.Time) bool {
	return func(t time.Time) bool { return l.AllowN(t, 1) }
}

// simulate returns the send times a greedy sender gets from allow over run.
func simulate(allow func(time.Time) bool, run, step time.Duration) []time.Time {
	start := time.Unix(0, 0)
	var sent []time.Time
	for at := start; at.Before(start.Add(run)); at = at.Add(step) {
		for allow(at) {
			sent = append(sent, at)
		}
	}
	return sent
}

// burst is how many sends share the first send's time.
func burst(sent []time.Time) int {
	n := 0
	for n < len(sent) && sent[n].Equal(sent[0]) {
		n++
	}
	return n
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
