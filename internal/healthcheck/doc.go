// Package healthcheck provides a liveness probe for the portwatch daemon.
//
// A Checker tracks the timestamp of the most recent successful port scan and
// compares it against a configurable staleness threshold. If no scan has
// completed within that window the daemon is considered unhealthy.
//
// Typical usage:
//
//	check := healthcheck.New(30 * time.Second)
//
//	// inside the scan loop:
//	check.RecordScan()
//
//	// on error:
//	check.RecordError()
//
//	// to inspect:
//	fmt.Println(check.Status())
package healthcheck
