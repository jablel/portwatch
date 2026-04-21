// Package portpin provides a mechanism to "pin" a set of expected ports
// and detect deviations at runtime.
//
// A Pinner maintains an explicit allowlist of ports that are expected to
// be open. On each scan cycle the caller passes the observed port list to
// Check, which returns Violations for:
//
//   - Unexpected ports: observed but not in the pinned set.
//   - Missing ports:    pinned but not present in the observed set.
//
// All operations are safe for concurrent use.
package portpin
