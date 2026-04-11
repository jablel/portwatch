// Package watchdog monitors the daemon's internal health and restarts
// stalled scan cycles by emitting a signal when a deadline is exceeded.
package watchdog

import (
	"context"
	"sync"
	"time"
)

// Watchdog tracks periodic "kicks" and fires a callback when a kick is
// not received within the configured deadline.
type Watchdog struct {
	mu       sync.Mutex
	deadline time.Duration
	lastKick time.Time
	onStall  func()
	cancel   context.CancelFunc
}

// New creates a Watchdog that calls onStall if no Kick is received within
// deadline. The watchdog begins monitoring immediately.
func New(deadline time.Duration, onStall func()) *Watchdog {
	ctx, cancel := context.WithCancel(context.Background())
	w := &Watchdog{
		deadline: deadline,
		lastKick: time.Now(),
		onStall:  onStall,
		cancel:   cancel,
	}
	go w.run(ctx)
	return w
}

// Kick resets the watchdog timer, signalling that the daemon is alive.
func (w *Watchdog) Kick() {
	w.mu.Lock()
	w.lastKick = time.Now()
	w.mu.Unlock()
}

// Stop shuts down the watchdog goroutine.
func (w *Watchdog) Stop() {
	w.cancel()
}

// Stalled reports whether the watchdog considers the daemon stalled.
func (w *Watchdog) Stalled() bool {
	if w.deadline <= 0 {
		return false
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	return time.Since(w.lastKick) > w.deadline
}

func (w *Watchdog) run(ctx context.Context) {
	if w.deadline <= 0 {
		return
	}
	ticker := time.NewTicker(w.deadline / 2)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if w.Stalled() {
				wdog/watchdog_test.gotest

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watchdog"
)

func TestStalled_FalseAfterKick(t *testing.T) {
	w := watchdog.New(200*time.Millisecond, func() {})
	defer w.Stop()

	w.Kick()
	if w.Stalled() {
		t.Fatal("expected not stalled immediately after kick")
	}
}

func TestStalled_TrueAfterDeadlineExceeded(t *testing.T) {
	w := watchdog.New(50*time.Millisecond, func() {})
	defer w.Stop()

	time.Sleep(80 * time.Millisecond)
	if !w.Stalled() {
		t.Fatal("expected stalled after deadline exceeded")
	}
}

func TestOnStall_CalledWhenStalled(t *testing.T) {
	var called atomic.Int32
	w := watchdog.New(40*time.Millisecond, func() {
		called.Add(1)
	})
	defer w.Stop()

	time.Sleep(120 * time.Millisecond)
	if called.Load() == 0 {
		t.Fatal("expected onStall to have been called")
	}
}

func TestKick_ResetsStallDetection(t *testing.T) {
	var called atomic.Int32
	w := watchdog.New(60*time.Millisecond, func() {
		called.Add(1)
	})
	defer w.Stop()

	// Keep kicking to prevent stall.
	for i := 0; i < 5; i++ {
		time.Sleep(20 * time.Millisecond)
		w.Kick()
	}
	if called.Load() != 0 {
		t.Fatal("expected onStall NOT to have been called while kicking")
	}
}

func TestStop_StopsMonitoring(t *testing.T) {
	var called atomic.Int32
	w := watchdog.New(30*time.Millisecond, func() {
		called.Add(1)
	})
	w.Stop()

	time.Sleep(100 * time.Millisecond)
	if called.Load() != 0 {
		t.Fatal("expected onStall not called after Stop")
	}
}

func TestZeroDeadline_NeverStalls(t *testing.T) {
	var called atomic.Int32
	w := watchdog.New(0, func() {
		called.Add(1)
	})
	defer w.Stop()

	time.Sleep(50 * time.Millisecond)
	if w.Stalled() {
		t.Fatal("zero deadline should never report stalled")
	}
	if called.Load() != 0 {
		t.Fatal("zero deadline should never call onStall")
	}
}
