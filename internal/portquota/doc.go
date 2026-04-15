// Package portquota provides per-protocol quota enforcement for open ports.
//
// A Quota holds a configurable limit for each network protocol (e.g. "tcp",
// "udp"). Calling Check against a slice of observed ports returns a
// []Violation describing every protocol whose open-port count exceeds its
// configured limit.
//
// Typical usage:
//
//	q := portquota.New()
//	q.Set("tcp", 50)
//	q.Set("udp", 20)
//
//	violations := q.Check(currentPorts)
//	for _, v := range violations {
//		log.Printf("quota breach: %s", v)
//	}
//
// A limit of zero is treated as unlimited and will never produce a violation.
// All methods are safe for concurrent use.
package portquota
