// Package portdrift detects when the set of open ports has drifted
// significantly from a known-good baseline by computing a drift score
// between 0.0 (identical) and 1.0 (completely different).
package portdrift

import (
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Tracker measures how much the current port set differs from a
// previously recorded baseline snapshot.
type Tracker struct {
	mu       sync.Mutex
	baseline map[string]struct{}
}

// New returns a Tracker with no baseline set.
func New() *Tracker {
	return &Tracker{}
}

// SetBaseline records the current port set as the reference snapshot.
func (t *Tracker) SetBaseline(ports []scanner.Port) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.baseline = toSet(ports)
}

// Score returns a drift value in [0.0, 1.0].
// 0.0 means the current ports exactly match the baseline.
// 1.0 means there is no overlap at all.
// If no baseline has been set, Score returns 0.0.
func (t *Tracker) Score(current []scanner.Port) float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.baseline) == 0 && len(current) == 0 {
		return 0.0
	}
	if len(t.baseline) == 0 {
		return 0.0
	}

	curSet := toSet(current)

	union := make(map[string]struct{}, len(t.baseline)+len(curSet))
	for k := range t.baseline {
		union[k] = struct{}{}
	}
	for k := range curSet {
		union[k] = struct{}{}
	}

	intersect := 0
	for k := range t.baseline {
		if _, ok := curSet[k]; ok {
		intersect++
		}
	}

	if len(union) == 0 {
		return 0.0
	}
	return 1.0 - float64(intersect)/float64(len(union))
}

// HasBaseline reports whether a baseline snapshot has been recorded.
func (t *Tracker) HasBaseline() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.baseline != nil
}

func toSet(ports []scanner.Port) map[string]struct{} {
	s := make(map[string]struct{}, len(ports))
	for _, p := range ports {
		s[p.String()] = struct{}{}
	}
	return s
}
