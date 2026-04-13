// Package portexpiry tracks the continuous open-duration of each observed port
// and identifies ports that have remained open longer than a configured
// maximum age.
//
// Usage:
//
//	tr := portexpiry.New(24 * time.Hour)
//
//	// Call Observe on each scan tick with the current open ports.
//	tr.Observe(ports)
//
//	// Retrieve any ports that have been open longer than maxAge.
//	for _, e := range tr.Expired() {
//		fmt.Printf("port %s open since %s\n", e.Port, e.FirstSeen)
//	}
package portexpiry
