// Package portecho tracks how frequently each port appears across
// successive scans, producing a normalised echo score in [0.0, 1.0].
// A score of 1.0 means the port was present in every observed scan;
// 0.0 means it has never been seen.
package portecho

import (
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Entry holds the observation counts for a single port.
type Entry struct {
	Seen  int
	Total int
}

// Score returns the ratio of scans in which the port was present.
func (e Entry) Score() float64 {
	if e.Total == 0 {
		return 0
	}
	return float64(e.Seen) / float64(e.Total)
}

// Tracker records per-port presence across scans.
type Tracker struct {
	mu      sync.Mutex
	entries map[string]*Entry
	scans   int
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{entries: make(map[string]*Entry)}
}

func portKey(p scanner.Port) string {
	return p.Proto + ":" + p.String()
}

// Record advances the scan counter and marks which ports were present.
func (t *Tracker) Record(ports []scanner.Port) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.scans++

	seen := make(map[string]struct{}, len(ports))
	for _, p := range ports {
		k := portKey(p)
		seen[k] = struct{}{}
		if t.entries[k] == nil {
			t.entries[k] = &Entry{}
		}
		t.entries[k].Seen++
	}

	// Bump Total for every tracked port on each scan.
	for k, e := range t.entries {
		e.Total++
		_ = k
	}
}

// Score returns the echo score for the given port, or 0 if unknown.
func (t *Tracker) Score(p scanner.Port) float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	e, ok := t.entries[portKey(p)]
	if !ok {
		return 0
	}
	return e.Score()
}

// All returns a snapshot of every tracked entry keyed by "proto:port".
func (t *Tracker) All() map[string]Entry {
	t.mu.Lock()
	defer t.mu.Unlock()

	out := make(map[string]Entry, len(t.entries))
	for k, e := range t.entries {
		out[k] = *e
	}
	return out
}
