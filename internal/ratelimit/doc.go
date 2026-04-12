// Package ratelimit implements a simple per-key cooldown rate limiter suited
// for suppressing repeated alerts about the same port.
//
// Usage:
//
//	limiter := ratelimit.New(30 * time.Second)
//
//	if limiter.Allow("tcp:8080:added") {
//		// send alert
//	}
//
// Keys are arbitrary strings; callers typically encode the port address and
// event kind into the key so that different event types for the same port are
// tracked independently.
//
// A zero or negative cooldown disables rate limiting — every call returns true.
// ResetAll can be used to flush all state, e.g. after a configuration reload.
package ratelimit
