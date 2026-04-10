// Package dedupe provides event deduplication to avoid re-alerting on
// ports that have already been reported in a recent scan cycle.
package dedupe

import (
	"sync"
	"time"

	"portwatch/internal/scanner"
)

// Entry records when a port event was last seen.
type Entry struct {
	SeenAt time.Time
}

// Deduper tracks recently seen port events and suppresses duplicates
// within a configurable retention window.
type Deduper struct {
	mu      sync.Mutex
	seen    map[string]Entry
	window  time.Duration
	nowFunc func() time.Time
}

// New creates a Deduper that suppresses repeated events within window.
func New(window time.Duration) *Deduper {
	return &Deduper{
		seen:    make(map[string]Entry),
		window:  window,
		nowFunc: time.Now,
	}
}

// IsDuplicate returns true if the port was already seen within the
// retention window. If it is not a duplicate the entry is recorded.
func (d *Deduper) IsDuplicate(p scanner.Port) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.nowFunc()
	key := p.String()

	if d.window <= 0 {
		d.seen[key] = Entry{SeenAt: now}
		return false
	}

	if e, ok := d.seen[key]; ok {
		if now.Sub(e.SeenAt) < d.window {
			return true
		}
	}

	d.seen[key] = Entry{SeenAt: now}
	return false
}

// Evict removes all entries whose last-seen time is older than the window.
func (d *Deduper) Evict() {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.nowFunc()
	for key, e := range d.seen {
		if now.Sub(e.SeenAt) >= d.window {
			delete(d.seen, key)
		}
	}
}

// Reset clears all tracked entries.
func (d *Deduper) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]Entry)
}
