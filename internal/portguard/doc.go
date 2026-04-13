// Package portguard provides an allowlist-based guard for open ports.
//
// A Guard holds a set of approved (port, protocol) pairs. After each
// scan cycle the caller passes the observed ports to Check; any port
// absent from the allowlist is returned as a Violation that the caller
// can forward to the notifier or event log.
//
// The allowlist is safe for concurrent use and can be modified at
// runtime via Allow and Revoke to support dynamic policy updates
// without restarting the daemon.
package portguard
