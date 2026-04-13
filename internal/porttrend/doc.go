// Package porttrend tracks how frequently individual ports appear across
// successive scan samples and derives a simple directional trend.
//
// Usage:
//
//	tr := porttrend.New(5 * time.Minute)
//
//	// Call Record after every scan.
//	tr.Record(scannedPorts)
//
//	// Query the trend for any port.
//	switch tr.Trend(p) {
//	case porttrend.Rising:
//	    // port is being seen more frequently
//	case porttrend.Falling:
//	    // port is disappearing
//	}
//
// Observations older than the configured window are evicted automatically
// on every call to Record or Trend, keeping memory usage bounded.
package porttrend
