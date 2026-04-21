// Package portsteady tracks how long a port has remained continuously open
// and exposes a stability score in the range [0, 1].
package portsteady

import (
	"fmt"
	"sync"
	"time"

	"portwatch/internal/scanner"
)

const defaultWindow = 10

type entry struct {
	firstSeen time.Time
	lastSeen  time.Time
	scans     int
}

// Tracker measures port stability across successive scans.
type Tracker struct {
	mu     sync.Mutex
	ports  map[string]*entry
	window int // number of recent scans considered
}

// New returns a Tracker that uses the given scan-window size.
// window <= 0 falls back to the default of 10.
func New(window int) *Tracker {
	if window <= 0 {
		window = defaultWindow
	}
	return &Tracker{
		ports:  make(map[string]*entry),
		window: window,
	}
}

func portKey(p scanner.Port) string {
	return fmt.Sprintf("%s:%d", p.Protocol, p.Number)
}

// Observe records one scan snapshot. Ports present in the snapshot are
// credited; ports absent are removed from the tracker.
func (t *Tracker) Observe(ports []scanner.Port, now time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()

	seen := make(map[string]struct{}, len(ports))
	for _, p := range ports {
		k := portKey(p)
		seen[k] = struct{}{}
		if e, ok := t.ports[k]; ok {
			e.lastSeen = now
			if e.scans < t.window {
				e.scans++
			}
		} else {
			t.ports[k] = &entry{firstSeen: now, lastSeen: now, scans: 1}
		}
	}

	for k := range t.ports {
		if _, ok := seen[k]; !ok {
			delete(t.ports, k)
		}
	}
}

// Stability returns a score in [0, 1] representing how consistently the port
// has appeared over the observation window. Returns 0 for unknown ports.
func (t *Tracker) Stability(p scanner.Port) float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	e, ok := t.ports[portKey(p)]
	if !ok {
		return 0
	}
	return float64(e.scans) / float64(t.window)
}

// Uptime returns how long the port has been continuously observed.
// Returns zero duration for unknown ports.
func (t *Tracker) Uptime(p scanner.Port, now time.Time) time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()

	e, ok := t.ports[portKey(p)]
	if !ok {
		return 0
	}
	return now.Sub(e.firstSeen)
}
