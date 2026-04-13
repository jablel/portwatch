// Package portlifecycle tracks the lifecycle of observed ports across scans.
//
// A Tracker is updated on every scan tick via Observe. Each port transitions
// through the following states:
//
//	"new"    – first observation in the current run
//	"active" – present in at least two consecutive scans
//	"closed" – was active but absent from the most recent scan
//
// Lifecycle data can be used by alerting and reporting components to suppress
// noise from ephemeral ports or to highlight long-lived unexpected listeners.
package portlifecycle
