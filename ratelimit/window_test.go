package ratelimit

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestWindowBurstsThenWaitsForOldestToAgeOut(t *testing.T) {
	const per = 150 * time.Millisecond
	w := NewWindow(3, per)
	ctx := context.Background()

	start := time.Now()
	for range 3 {
		if err := w.Wait(ctx); err != nil {
			t.Fatalf("Wait() error = %v", err)
		}
	}
	if d := time.Since(start); d > 50*time.Millisecond {
		t.Errorf("first 3 took %s, want them at once", d)
	}

	if err := w.Wait(ctx); err != nil {
		t.Fatalf("Wait() error = %v", err)
	}
	if d := time.Since(start); d < per-10*time.Millisecond {
		t.Errorf("4th request went after %s, want it held until the first aged out (%s)", d, per)
	}
}

func TestWindowFailsFastPastDeadline(t *testing.T) {
	w := NewWindow(1, time.Hour)
	if err := w.Wait(context.Background()); err != nil {
		t.Fatalf("Wait() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	start := time.Now()
	err := w.Wait(ctx)
	if err == nil || errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Wait() error = %v, want an immediate refusal", err)
	}
	if d := time.Since(start); d > 50*time.Millisecond {
		t.Errorf("refusal took %s, want it immediate", d)
	}
}

func TestWindowHonoursCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := NewWindow(1, time.Second).Wait(ctx); !errors.Is(err, context.Canceled) {
		t.Errorf("Wait() error = %v, want context.Canceled", err)
	}
}

// TestWindowConcurrentWaiters sends 20 requests from 20 goroutines through a
// 5-per-100ms window: four full windows, so the last must wait about 300ms.
func TestWindowConcurrentWaiters(t *testing.T) {
	w := NewWindow(5, 100*time.Millisecond)
	start := time.Now()
	var wg sync.WaitGroup
	for range 20 {
		wg.Go(func() {
			if err := w.Wait(context.Background()); err != nil {
				t.Errorf("Wait() error = %v", err)
			}
		})
	}
	wg.Wait()
	if d := time.Since(start); d < 290*time.Millisecond || d > time.Second {
		t.Errorf("20 requests took %s, want about 300ms", d)
	}
}
