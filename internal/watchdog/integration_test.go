package watchdog_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watchdog"
)

// TestWatchdog_MultipleStallCallbacks verifies the callback fires more than
// once across successive stall windows.
func TestWatchdog_MultipleStallCallbacks(t *testing.T) {
	var mu sync.Mutex
	var calls []time.Time

	w := watchdog.New(30*time.Millisecond, func() {
		mu.Lock()
		calls = append(calls, time.Now())
		mu.Unlock()
	})
	defer w.Stop()

	time.Sleep(150 * time.Millisecond)

	mu.Lock()
	count := len(calls)
	mu.Unlock()

	if count < 2 {
		t.Fatalf("expected at least 2 stall callbacks, got %d", count)
	}
}

// TestWatchdog_KickPreventsAnyCallback verifies that continuous kicking
// keeps the callback from ever firing.
func TestWatchdog_KickPreventsAnyCallback(t *testing.T) {
	var mu sync.Mutex
	var calls int

	w := watchdog.New(50*time.Millisecond, func() {
		mu.Lock()
		calls++
		mu.Unlock()
	})
	defer w.Stop()

	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				w.Kick()
			case <-done:
				return
			}
		}
	}()

	time.Sleep(200 * time.Millisecond)
	close(done)

	mu.Lock()
	defer mu.Unlock()
	if calls != 0 {
		t.Fatalf("expected 0 stall callbacks while kicking, got %d", calls)
	}
}
