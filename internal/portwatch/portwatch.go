// Package portwatch ties together scanning, diffing, and alerting into a
// single high-level coordinator used by the daemon tick loop.
package portwatch

import (
	"context"
	"fmt"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
	"github.com/user/portwatch/internal/notifier"
)

// Result holds the outcome of a single watch cycle.
type Result struct {
	Ports    []scanner.Port
	Diff     state.Diff
	ScannedAt time.Time
}

// Watcher coordinates one scan-diff-notify cycle.
type Watcher struct {
	scanner  *scanner.Scanner
	notifier *notifier.Notifier
	prev     []scanner.Port
}

// New creates a Watcher with the given scanner and notifier.
func New(s *scanner.Scanner, n *notifier.Notifier) *Watcher {
	return &Watcher{scanner: s, notifier: n}
}

// Tick runs one scan cycle: scan ports, compute diff against previous state,
// send notifications for any changes, and update internal state.
// It returns a Result describing what was observed.
func (w *Watcher) Tick(ctx context.Context) (Result, error) {
	ports, err := w.scanner.Scan(ctx)
	if err != nil {
		return Result{}, fmt.Errorf("portwatch: scan failed: %w", err)
	}

	diff := state.Compare(w.prev, ports)

	if err := w.notifier.Notify(diff); err != nil {
		return Result{}, fmt.Errorf("portwatch: notify failed: %w", err)
	}

	w.prev = ports

	return Result{
		Ports:     ports,
		Diff:      diff,
		ScannedAt: time.Now(),
	}, nil
}

// Reset clears the previously observed port list so the next Tick treats
// all discovered ports as newly added.
func (w *Watcher) Reset() {
	w.prev = nil
}
