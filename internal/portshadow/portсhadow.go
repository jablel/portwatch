// Package portсhadow tracks ports that appear briefly and disappear
// within a single observation window — so-called "shadow" ports that
// may indicate transient or suspicious activity.
package portсhadow

import (
	"sync"
	"time"

	"portwatch/internal/scanner"
)

// Entry records when a shadow port was first and last seen.
type Entry struct {
	Port      scanner.Port
	FirstSeen time.Time
	LastSeen  time.Time
	Count     int
}

// Tracker detects ports that appear and disappear within a short window.
type Tracker struct {
	mu      sync.Mutex
	window  time.Duration
	active  map[string]*Entry // ports seen in current window
	shadows []Entry           // confirmed shadow ports
	now     func() time.Time
}

// New returns a Tracker that considers a port a shadow if it vanishes
// within window after first being seen.
func New(window time.Duration) *Tracker {
	return &Tracker{
		window: window,
		active: make(map[string]*Entry),
		now:    time.Now,
	}
}

func portKey(p scanner.Port) string {
	return p.Protocol + ":" + p.String()
}

// Observe updates the tracker with the current set of open ports.
// Ports that were active but are no longer present are evaluated;
// those that disappeared within the window are recorded as shadows.
func (t *Tracker) Observe(ports []scanner.Port) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	current := make(map[string]struct{}, len(ports))
	for _, p := range ports {
		key := portKey(p)
		current[key] = struct{}{}
		if e, ok := t.active[key]; ok {
			e.LastSeen = now
			e.Count++
		} else {
			t.active[key] = &Entry{
				Port:      p,
				FirstSeen: now,
				LastSeen:  now,
				Count:     1,
			}
		}
	}

	for key, e := range t.active {
		if _, seen := current[key]; !seen {
			if t.window > 0 && e.LastSeen.Sub(e.FirstSeen) <= t.window {
				t.shadows = append(t.shadows, *e)
			}
			delete(t.active, key)
		}
	}
}

// Shadows returns all shadow port entries recorded so far and resets
// the internal list.
func (t *Tracker) Shadows() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Entry, len(t.shadows))
	copy(out, t.shadows)
	t.shadows = t.shadows[:0]
	return out
}

// Active returns the ports currently being tracked (seen but not yet gone).
func (t *Tracker) Active() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Entry, 0, len(t.active))
	for _, e := range t.active {
		out = append(out, *e)
	}
	return out
}
