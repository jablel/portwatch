// Package ratelimit implements per-key cooldown-based rate limiting for
// portwatch alert suppression.
//
// # Overview
//
// When portwatch detects port changes it may emit alerts through one or more
// notifier backends. Without rate limiting, a flapping port (one that opens
// and closes rapidly) can flood log files or external notification services.
//
// The Limiter type tracks the last time an alert was emitted for a given key
// (typically a string like "80/tcp" or "added:443/tcp") and suppresses
// repeated alerts that arrive before the configured cooldown has elapsed.
//
// # Usage
//
//	limiter := ratelimit.New(30 * time.Second)
//
//	if limiter.Allow("added:8080/tcp") {
//	    notifier.Notify(diff)
//	}
//
// A cooldown of zero or any negative value disables rate limiting entirely,
// which is useful during testing or when the operator explicitly opts out.
package ratelimit
