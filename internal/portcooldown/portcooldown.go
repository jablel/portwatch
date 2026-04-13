// Package portcooldown tracks how long a port has been continuously absent
// since it was last seen open. It is useful for suppressing re-alerts for
// ports that flap closed and reopen within a short grace period.
package portcooldown

import (
	"fmt"
	"sync"
	"time"

	"portwatch/internal/scanner"
)

// Tracker records the time a port was last seen and exposes how long it has
// been absent from subsequent scans.
type Tracker struct {
	mu       sync.Mutex
	lastSeen map[string]time.Time
	clock    func() time.Time
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{
		lastSeen: make(map[string]time.Time),
		clock:    time.Now,
	}
}

func portKey(p scanner.Port) string {
	return fmt.Sprintf("%s:%d", p.Protocol, p.Number)
}

// Observe records that the given ports were present in the latest scan.
// Ports not included are considered absent from this point forward.
func (t *Tracker) Observe(ports []scanner.Port) {
	now := t.clock()
	seen := make(map[string]struct{}, len(ports))
	for _, p := range ports {
		k := portKey(p)
		seen[k] = struct{}{}
		t.mu.Lock()
		t.lastSeen[k] = now
		t.mu.Unlock()
	}

	// Remove entries for ports that are still present so absence timer resets.
	t.mu.Lock()
	for k := range t.lastSeen {
		if _, ok := seen[k]; !ok {
			// Keep the entry — it records when the port was last seen.
			_ = k
		}
	}
	t.mu.Unlock()
}

// AbsentFor returns how long the port has been absent since it was last
// observed. If the port has never been observed, the second return value is
// false.
func (t *Tracker) AbsentFor(p scanner.Port) (time.Duration, bool) {
	t.mu.Lock()
	ts, ok := t.lastSeen[portKey(p)]
	t.mu.Unlock()
	if !ok {
		return 0, false
	}
	return t.clock().Sub(ts), true
}

// InCooldown reports whether the port became absent less than d ago.
func (t *Tracker) InCooldown(p scanner.Port, d time.Duration) bool {
	if d <= 0 {
		return false
	}
	absent, ok := t.AbsentFor(p)
	if !ok {
		return false
	}
	return absent < d
}

// Reset removes all tracking state.
func (t *Tracker) Reset() {
	t.mu.Lock()
	t.lastSeen = make(map[string]time.Time)
	t.mu.Unlock()
}
