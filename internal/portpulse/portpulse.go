// Package portpulse tracks the frequency at which individual ports are
// observed across successive scans, producing a normalized pulse score
// between 0.0 (never seen) and 1.0 (seen in every scan).
package portpulse

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Tracker maintains per-port observation counts over a fixed window of scans.
type Tracker struct {
	mu       sync.RWMutex
	window   int
	counts   map[string]int
	scans    int
}

// New returns a Tracker that considers the last window scans when computing
// pulse scores. window must be >= 1; values below 1 are clamped to 1.
func New(window int) *Tracker {
	if window < 1 {
		window = 1
	}
	return &Tracker{
		window: window,
		counts: make(map[string]int),
	}
}

func portKey(p scanner.Port) string {
	return fmt.Sprintf("%d/%s", p.Number, p.Protocol)
}

// Record registers a set of ports observed in a single scan tick.
// Counts are clamped to the configured window so stale data cannot
// inflate scores indefinitely.
func (t *Tracker) Record(ports []scanner.Port) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.scans < t.window {
		t.scans++
	}

	seen := make(map[string]bool, len(ports))
	for _, p := range ports {
		k := portKey(p)
		seen[k] = true
	}

	for k, c := range t.counts {
		if seen[k] {
			if c < t.window {
				t.counts[k] = c + 1
			}
		} else {
			if c > 0 {
				t.counts[k] = c - 1
			}
			if t.counts[k] == 0 {
				delete(t.counts, k)
			}
		}
	}

	for k := range seen {
		if _, exists := t.counts[k]; !exists {
			t.counts[k] = 1
		}
	}
}

// Pulse returns a score in [0.0, 1.0] representing how consistently the
// given port has been observed. Returns 0 for unknown ports.
func (t *Tracker) Pulse(p scanner.Port) float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.scans == 0 {
		return 0
	}
	c := t.counts[portKey(p)]
	return float64(c) / float64(t.window)
}

// All returns a map of port keys to their current pulse scores.
func (t *Tracker) All() map[string]float64 {
	t.mu.RLock()
	defer t.mu.RUnlock()

	out := make(map[string]float64, len(t.counts))
	for k, c := range t.counts {
		out[k] = float64(c) / float64(t.window)
	}
	return out
}
