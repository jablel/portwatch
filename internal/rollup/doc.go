// Package rollup provides a time-windowed accumulator for port-change diffs.
//
// When the monitored system experiences a burst of port changes — for example
// during a service restart — individual diffs arrive in rapid succession.
// Forwarding each diff as a separate alert would overwhelm operators.
//
// Rollup solves this by collecting all diffs that arrive within a configurable
// window and merging them into a single summary diff before invoking the
// downstream handler. The window is reset each time a new diff arrives, so
// the handler is only called once the activity has settled.
//
// Usage:
//
//	r := rollup.New(500*time.Millisecond, func(d rollup.Diff) {
//	    notifier.Notify(d)
//	})
//	// Feed diffs from the watcher:
//	r.Add(diff)
package rollup
