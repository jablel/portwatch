// Package rollup groups rapid port-change events into a single summary
// notification, reducing alert noise during bursts of activity.
package rollup

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Diff mirrors the state.Diff type used across the project.
type Diff = state.Diff

// Handler is called with the accumulated diff once the window closes.
type Handler func(d Diff)

// Rollup collects diffs within a time window and merges them before
// forwarding to the registered handler.
type Rollup struct {
	mu      sync.Mutex
	window  time.Duration
	handler Handler
	pending *Diff
	timer   *time.Timer
	clock   func() time.Time
}

// New creates a Rollup that flushes accumulated diffs after window elapses
// with no new activity.
func New(window time.Duration, h Handler) *Rollup {
	return &Rollup{
		window:  window,
		handler: h,
		clock:   time.Now,
	}
}

// Add merges d into the pending diff and resets the flush timer.
func (r *Rollup) Add(d Diff) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.pending == nil {
		copy := Diff{
			Added:   append([]state.Port(nil), d.Added...),
			Removed: append([]state.Port(nil), d.Removed...),
		}
		r.pending = &copy
	} else {
		r.pending.Added = append(r.pending.Added, d.Added...)
		r.pending.Removed = append(r.pending.Removed, d.Removed...)
	}

	if r.window <= 0 {
		r.flush()
		return
	}

	if r.timer != nil {
		r.timer.Reset(r.window)
	} else {
		r.timer = time.AfterFunc(r.window, r.flushLocked)
	}
}

// Flush forces an immediate flush of any pending diff.
func (r *Rollup) Flush() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.flush()
}

// flush must be called with r.mu held.
func (r *Rollup) flush() {
	if r.timer != nil {
		r.timer.Stop()
		r.timer = nil
	}
	if r.pending == nil {
		return
	}
	d := *r.pending
	r.pending = nil
	go r.handler(d)
}

// flushLocked acquires the lock then flushes; used by time.AfterFunc.
func (r *Rollup) flushLocked() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.flush()
}
