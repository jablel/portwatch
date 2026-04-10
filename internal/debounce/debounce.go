// Package debounce provides a mechanism to delay and coalesce rapid
// successive events for the same key, emitting only the final event
// after a quiet period has elapsed.
package debounce

import (
	"sync"
	"time"
)

// Func is a callback invoked with the key after the debounce period.
type Func func(key string)

// Debouncer delays calls for a given key until no new calls have been
// made for the configured wait duration.
type Debouncer struct {
	wait  time.Duration
	mu    sync.Mutex
	timers map[string]*time.Timer
}

// New creates a Debouncer that waits for the given duration of inactivity
// before firing the callback.
func New(wait time.Duration) *Debouncer {
	return &Debouncer{
		wait:  wait,
		timers: make(map[string]*time.Timer),
	}
}

// Trigger schedules fn to be called with key after the wait period.
// If Trigger is called again for the same key before the timer fires,
// the timer is reset and fn will be called only once.
func (d *Debouncer) Trigger(key string, fn Func) {
	if d.wait <= 0 {
		fn(key)
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[key]; ok {
		t.Stop()
	}

	d.timers[key] = time.AfterFunc(d.wait, func() {
		d.mu.Lock()
		delete(d.timers, key)
		d.mu.Unlock()
		fn(key)
	})
}

// Cancel stops any pending timer for the given key without invoking the callback.
func (d *Debouncer) Cancel(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[key]; ok {
		t.Stop()
		delete(d.timers, key)
	}
}

// Pending returns the number of keys with active pending timers.
func (d *Debouncer) Pending() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.timers)
}
