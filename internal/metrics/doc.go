// Package metrics provides lightweight runtime counters for the portwatch daemon.
//
// It tracks the number of completed scans, total ports observed, alerts emitted,
// and errors encountered since the daemon started. All methods are safe for
// concurrent use.
//
// Usage:
//
//	m := metrics.New()
//	m.RecordScan(len(ports))
//	m.RecordAlert()
//	m.Write(os.Stdout)
//
// Counters are kept in memory only; they reset when the daemon restarts.
package metrics
