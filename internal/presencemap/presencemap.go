// Package presencemap tracks which ports have been continuously present
// across consecutive scans, enabling stable-port detection.
package presencemap

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Entry records how many consecutive scans a port has been observed.
type Entry struct {
	Port   scanner.Port
	Streak int
}

// Map maintains per-port presence streaks.
type Map struct {
	mu      sync.RWMutex
	entries map[string]*Entry
}

// New returns an initialised Map.
func New() *Map {
	return &Map{entries: make(map[string]*Entry)}
}

func portKey(p scanner.Port) string {
	return fmt.Sprintf("%s/%d", p.Protocol, p.Number)
}

// Observe updates streaks given the current set of open ports.
// Ports not present in current are removed from the map.
func (m *Map) Observe(current []scanner.Port) {
	m.mu.Lock()
	defer m.mu.Unlock()

	seen := make(map[string]struct{}, len(current))
	for _, p := range current {
		k := portKey(p)
		seen[k] = struct{}{}
		if e, ok := m.entries[k]; ok {
			e.Streak++
		} else {
			m.entries[k] = &Entry{Port: p, Streak: 1}
		}
	}

	for k := range m.entries {
		if _, ok := seen[k]; !ok {
			delete(m.entries, k)
		}
	}
}

// Streak returns the current consecutive-scan count for a port.
// Returns 0 if the port is not tracked.
func (m *Map) Streak(p scanner.Port) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if e, ok := m.entries[portKey(p)]; ok {
		return e.Streak
	}
	return 0
}

// Stable returns all ports whose streak is at least minStreak.
func (m *Map) Stable(minStreak int) []scanner.Port {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []scanner.Port
	for _, e := range m.entries {
		if e.Streak >= minStreak {
			out = append(out, e.Port)
		}
	}
	return out
}

// Reset clears all tracked entries.
func (m *Map) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = make(map[string]*Entry)
}
