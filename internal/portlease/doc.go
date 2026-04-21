// Package portlease provides TTL-based port lease tracking for portwatch.
//
// A lease declares that a specific port is expected to be open for a limited
// duration. The Tracker monitors active ports against granted leases and
// surfaces violations when a port remains open past its lease expiry.
//
// Typical usage:
//
//	tr := portlease.New()
//	tr.Grant(port, 30*time.Minute)
//	// ... later, after scanning ...
//	violations := tr.Check(currentPorts)
//	for _, v := range violations {
//		fmt.Println(v)
//	}
package portlease
