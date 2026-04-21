// Package portchurn tracks the rate at which ports open and close over a
// sliding observation window, producing a churn score between 0.0 (perfectly
// stable) and 1.0 (completely replaced every scan).
package portchurn

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// entry records a single scan observation.
type entry struct {
	ports map[string]struct{}
	at    time.Time
}

// Tracker accumulates scan snapshots and computes a churn score.
type Tracker struct {
	mu      sync.Mutex
	window  time.Duration
	entries []entry
}

// New returns a Tracker that retains observations within window.
// A zero or negative window keeps all observations.
func New(window time.Duration) *Tracker {
	return &Tracker{window: window}
}

// Record registers the current set of ports at the given timestamp.
// It evicts observations older than the configured window before recording.
func (t *Tracker) Record(ports []scanner.Port, now time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.window > 0 {
		cutoff := now.Add(-t.window)
		i := 0
		for i < len(t.entries) && t.entries[i].at.Before(cutoff) {
			i++
		}
		t.entries = t.entries[i:]
	}

	set := make(map[string]struct{}, len(ports))
	for _, p := range ports {
		set[p.String()] = struct{}{}
	}
	t.entries = append(t.entries, entry{ports: set, at: now})
}

// Churn returns the fraction of unique ports that changed (appeared or
// disappeared) across all retained observations.  Returns 0 when fewer than
// two observations are available.
func (t *Tracker) Churn() float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.entries) < 2 {
		return 0
	}

	all := make(map[string]struct{})
	for _, e := range t.entries {
		for k := range e.ports {
			all[k] = struct{}{}
		}
	}

	// A port is "churned" if it is absent from at least one observation.
	churned := 0
	for k := range all {
		for _, e := range t.entries {
			if _, ok := e.ports[k]; !ok {
				churned++
				break
			}
		}
	}

	if len(all) == 0 {
		return 0
	}
	return float64(churned) / float64(len(all))
}

// Reset clears all retained observations.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = nil
}
