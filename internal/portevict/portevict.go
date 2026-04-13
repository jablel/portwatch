// Package portevict tracks ports that have been evicted (closed after being
// open for a sustained period) and records when and how long they were active.
package portevict

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Record holds eviction metadata for a single port.
type Record struct {
	Port      scanner.Port
	FirstSeen time.Time
	LastSeen  time.Time
	EvictedAt time.Time
	Duration  time.Duration
}

// Tracker records ports that disappear after being observed.
type Tracker struct {
	mu      sync.Mutex
	active  map[string]activeEntry
	evicted []Record
	maxLen  int
}

type activeEntry struct {
	port      scanner.Port
	firstSeen time.Time
	lastSeen  time.Time
}

func portKey(p scanner.Port) string {
	return p.Proto + ":" + p.String()
}

// New creates a Tracker that retains at most maxLen eviction records.
func New(maxLen int) *Tracker {
	if maxLen <= 0 {
		maxLen = 256
	}
	return &Tracker{
		active: make(map[string]activeEntry),
		maxLen: maxLen,
	}
}

// Observe updates internal state from the current scan result.
// Ports present in current are marked active; ports previously active but
// absent from current are evicted.
func (t *Tracker) Observe(current []scanner.Port, now time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()

	seen := make(map[string]scanner.Port, len(current))
	for _, p := range current {
		k := portKey(p)
		seen[k] = p
		if e, ok := t.active[k]; ok {
			e.lastSeen = now
			t.active[k] = e
		} else {
			t.active[k] = activeEntry{port: p, firstSeen: now, lastSeen: now}
		}
	}

	for k, e := range t.active {
		if _, ok := seen[k]; !ok {
			rec := Record{
				Port:      e.port,
				FirstSeen: e.firstSeen,
				LastSeen:  e.lastSeen,
				EvictedAt: now,
				Duration:  e.lastSeen.Sub(e.firstSeen),
			}
			if len(t.evicted) >= t.maxLen {
				t.evicted = t.evicted[1:]
			}
			t.evicted = append(t.evicted, rec)
			delete(t.active, k)
		}
	}
}

// Evicted returns a copy of all eviction records, oldest first.
func (t *Tracker) Evicted() []Record {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Record, len(t.evicted))
	copy(out, t.evicted)
	return out
}

// ActiveCount returns the number of currently tracked active ports.
func (t *Tracker) ActiveCount() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.active)
}
