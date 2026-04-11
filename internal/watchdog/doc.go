// Package watchdog provides a deadline-based liveness monitor for the
// portwatch daemon.
//
// The daemon calls Kick on each completed scan cycle. If no kick arrives
// within the configured deadline the watchdog fires an onStall callback,
// allowing the caller to log an error, emit a metric, or attempt recovery.
//
// Usage:
//
//	wd := watchdog.New(30*time.Second, func() {
//		log.Println("scan cycle stalled — check system resources")
//	})
//	defer wd.Stop()
//
//	// inside the scan loop:
//	wd.Kick()
package watchdog
