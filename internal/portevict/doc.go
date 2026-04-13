// Package portevict tracks ports that disappear after being observed across
// one or more scan cycles.
//
// A Tracker is updated via Observe on each scan tick. When a port that was
// previously active is absent from the current scan, it is recorded as evicted
// along with its first-seen time, last-seen time, and total active duration.
//
// Eviction records are kept in a bounded ring buffer; the oldest record is
// dropped when capacity is exceeded.
//
// Typical usage:
//
//	tr := portevict.New(256)
//	// on each scan tick:
//	tr.Observe(currentPorts, time.Now())
//	// inspect closed ports:
//	for _, rec := range tr.Evicted() {
//		fmt.Printf("port %s closed after %v\n", rec.Port, rec.Duration)
//	}
package portevict
