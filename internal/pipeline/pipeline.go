// Package pipeline wires together scanning, filtering, diffing, and alerting
// into a single reusable processing step.
package pipeline

import (
	"context"
	"fmt"
	"time"

	"portwatch/internal/filter"
	"portwatch/internal/scanner"
	"portwatch/internal/state"
	"portwatch/internal/notifier"
)

// Result holds the outcome of a single pipeline run.
type Result struct {
	Ports    []scanner.Port
	Diff     state.Diff
	ScannedAt time.Time
}

// Pipeline executes one scan-filter-diff-notify cycle.
type Pipeline struct {
	scanner  *scanner.Scanner
	filter   *filter.Filter
	notifier *notifier.Notifier
}

// New creates a Pipeline from its component dependencies.
func New(sc *scanner.Scanner, f *filter.Filter, n *notifier.Notifier) *Pipeline {
	return &Pipeline{
		scanner:  sc,
		filter:   f,
		notifier: n,
	}
}

// Run performs a full scan cycle against the provided previous port list.
// It returns a Result containing the current ports and the computed diff.
func (p *Pipeline) Run(ctx context.Context, previous []scanner.Port) (*Result, error) {
	ports, err := p.scanner.Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("pipeline scan: %w", err)
	}

	filtered := p.filter.Apply(ports)

	diff := state.Compare(previous, filtered)

	if err := p.notifier.Notify(diff); err != nil {
		return nil, fmt.Errorf("pipeline notify: %w", err)
	}

	return &Result{
		Ports:     filtered,
		Diff:      diff,
		ScannedAt: time.Now(),
	}, nil
}
