// Package portcache provides a lightweight, thread-safe, TTL-bounded cache
// for the results of a port scan.
//
// Typical usage:
//
//	c := portcache.New(30 * time.Second)
//	c.Set(ports)
//	if entry, ok := c.Get(); ok {
//		// use entry.Ports
//	}
//
// A zero TTL disables expiry so that entries remain valid until explicitly
// invalidated with Invalidate.
package portcache
