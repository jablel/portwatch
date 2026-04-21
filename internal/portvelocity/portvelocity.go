// Package portvelocity tracks the rate of change for observed ports across
// successive scans, producing a velocity score in the range [0, 1] where 0
// means no change and 1 means every port changed between scans.
package portvelocity

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Tracker measures the churn rate between consecutive port scans.
type Tracker struct {
	mu   sync.Mutex
	prev map[string]struct{}
	last float64
}

// New returns a new Tracker with no prior scan data.
func New() *Tracker {
	return &Tracker{}
}

func portKey(p scanner.Port) string {
	return fmt.Sprintf("%s:%d", p.Protocol, p.Number)
}

// Record accepts the current set of ports, computes the velocity relative to
// the previous scan, stores the current set as the new baseline, and returns
// the velocity score.
//
// Velocity is defined as:
//
//	(added + removed) / max(len(prev), len(current), 1)
//
// The first call always returns 0 because there is no prior scan to compare.
func (t *Tracker) Record(ports []scanner.Port) float64 {
	current := make(map[string]struct{}, len(ports))
	for _, p := range ports {
		current[portKey(p)] = struct{}{}
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.prev == nil {
		t.prev = current
		t.last = 0
		return 0
	}

	added := 0
	for k := range current {
		if _, ok := t.prev[k]; !ok {
			added++
		}
	}

	removed := 0
	for k := range t.prev {
		if _, ok := current[k]; !ok {
			removed++
		}
	}

	denom := len(t.prev)
	if len(current) > denom {
		denom = len(current)
	}
	if denom == 0 {
		denom = 1
	}

	v := float64(added+removed) / float64(denom)
	if v > 1 {
		v = 1
	}

	t.prev = current
	t.last = v
	return v
}

// Last returns the velocity score computed during the most recent call to
// Record. It returns 0 if Record has never been called.
func (t *Tracker) Last() float64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.last
}

// Reset clears all accumulated state, including the previous scan snapshot.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.prev = nil
	t.last = 0
}
