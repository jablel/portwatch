// Package ratelimit provides per-key cooldown-based rate limiting for
// port change alerts. It prevents alert storms when a port flaps rapidly
// by suppressing repeated notifications within a configurable cooldown window.
//
// Usage:
//
//	limiter := ratelimit.New(5 * time.Second)
//	if limiter.Allow("tcp:8080") {
//		// send alert
//	}
package ratelimit
