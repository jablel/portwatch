// Package portlifecycle tracks the lifecycle state of observed ports,
// recording when they first appeared, when they were last seen, and how
// many consecutive scans they have been present or absent.
package portlifecycle

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// State describes the current lifecycle phase of a port.
type State string

const (
	StateNew    State = "new"
	StateActive State = "active"
	StateClosed State = "closed"
)

// Entry holds lifecycle metadata for a single port.
type Entry struct {
	Port       scanner.Port
	State      State
	FirstSeen  time.Time
	LastSeen   time.Time
	SeenCount  int
	MissCount  int
}

// String returns a human-readable summary of the entry.
func (e Entry) String() string {
	return fmt.Sprintf("%s state=%s seen=%d miss=%d", e.Port, e.State, e.SeenCount, e.MissCount)
}

// Tracker maintains lifecycle entries for all observed ports.
type Tracker struct {
	mu      sync.Mutex
	entries map[string]*Entry
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{entries: make(map[string]*Entry)}
}

func portKey(p scanner.Port) string {
	return fmt.Sprintf("%s/%d", p.Protocol, p.Number)
}

// Observe updates lifecycle state for the provided set of currently open ports.
// Ports not present in the slice have their MissCount incremented and are
// transitioned to StateClosed.
func (t *Tracker) Observe(ports []scanner.Port, now time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()

	seen := make(map[string]struct{}, len(ports))
	for _, p := range ports {
		k := portKey(p)
		seen[k] = struct{}{}
		if e, ok := t.entries[k]; ok {
			e.LastSeen = now
			e.SeenCount++
			e.MissCount = 0
			e.State = StateActive
		} else {
			t.entries[k] = &Entry{
				Port:      p,
				State:     StateNew,
				FirstSeen: now,
				LastSeen:  now,
				SeenCount: 1,
			}
		}
	}

	for k, e := range t.entries {
		if _, present := seen[k]; !present && e.State != StateClosed {
			e.MissCount++
			e.State = StateClosed
		}
	}
}

// Get returns the lifecycle entry for the given port, and whether it exists.
func (t *Tracker) Get(p scanner.Port) (Entry, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[portKey(p)]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// All returns a snapshot of every tracked entry.
func (t *Tracker) All() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, *e)
	}
	return out
}
