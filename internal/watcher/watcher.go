package watcher

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

// ChangeEvent holds the diff result and the timestamp when it was detected.
type ChangeEvent struct {
	Diff      state.Diff
	DetectedAt time.Time
}

// Watcher polls the port scanner at a fixed interval and emits change events.
type Watcher struct {
	scanner  *scanner.Scanner
	interval time.Duration
	events   chan ChangeEvent
	stop     chan struct{}
}

// New creates a Watcher that scans using s every interval duration.
func New(s *scanner.Scanner, interval time.Duration) *Watcher {
	return &Watcher{
		scanner:  s,
		interval: interval,
		events:   make(chan ChangeEvent, 8),
		stop:     make(chan struct{}),
	}
}

// Events returns the read-only channel of detected change events.
func (w *Watcher) Events() <-chan ChangeEvent {
	return w.events
}

// Start begins polling in a background goroutine.
// The previous snapshot is provided so the first tick can detect changes
// relative to a known baseline.
func (w *Watcher) Start(baseline []scanner.Port) {
	go func() {
		prev := baseline
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-w.stop:
				close(w.events)
				return
			case <-ticker.C:
				current, err := w.scanner.Scan()
				if err != nil {
					continue
				}
				diff := state.Compare(prev, current)
				if len(diff.Added) > 0 || len(diff.Removed) > 0 {
					w.events <- ChangeEvent{
						Diff:       diff,
						DetectedAt: time.Now(),
					}
				}
				prev = current
			}
		}
	}()
}

// Stop signals the background goroutine to exit.
func (w *Watcher) Stop() {
	close(w.stop)
}
