// Package portfreq tracks how frequently each port appears across scans,
// providing a normalised frequency score in the range [0.0, 1.0].
package portfreq

import (
	"fmt"
	"sync"

	"portwatch/internal/scanner"
)

// Tracker counts port observations across a fixed-size scan window.
type Tracker struct {
	mu      sync.Mutex
	window  int
	scans   int
	counts  map[string]int
}

// New returns a Tracker that considers the last windowSize scans.
// windowSize must be >= 1; values below 1 are clamped to 1.
func New(windowSize int) *Tracker {
	if windowSize < 1 {
		windowSize = 1
	}
	return &Tracker{
		window: windowSize,
		counts: make(map[string]int),
	}
}

func portKey(p scanner.Port) string {
	return fmt.Sprintf("%d/%s", p.Number, p.Protocol)
}

// Record registers one scan's worth of observed ports.
// Counts are capped at the configured window size to avoid unbounded growth.
func (t *Tracker) Record(ports []scanner.Port) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.scans < t.window {
		t.scans++
	}

	seen := make(map[string]struct{}, len(ports))
	for _, p := range ports {
		k := portKey(p)
		seen[k] = struct{}{}
		if t.counts[k] < t.window {
			t.counts[k]++
		}
	}

	// Decay counts for ports not observed this scan.
	for k, c := range t.counts {
		if _, ok := seen[k]; !ok && c > 0 {
			t.counts[k] = c - 1
			if t.counts[k] == 0 {
				delete(t.counts, k)
			}
		}
	}
}

// Frequency returns a value in [0.0, 1.0] representing how often the port
// has appeared across the tracked window. Returns 0 for unknown ports.
func (t *Tracker) Frequency(p scanner.Port) float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.scans == 0 {
		return 0
	}
	return float64(t.counts[portKey(p)]) / float64(t.scans)
}

// All returns a snapshot of frequency scores for every tracked port.
func (t *Tracker) All() map[string]float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	out := make(map[string]float64, len(t.counts))
	if t.scans == 0 {
		return out
	}
	for k, c := range t.counts {
		out[k] = float64(c) / float64(t.scans)
	}
	return out
}
