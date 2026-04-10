// Package digest provides a lightweight fingerprinting mechanism for sets of
// open ports.
//
// A Digest is a hex-encoded SHA-256 hash computed from the sorted string
// representations of all ports in a snapshot. It is designed for fast
// equality checks: if two consecutive scans produce the same Digest, no diff
// computation or alerting is necessary.
//
// Usage:
//
//	ports := scanner.Scan(cfg)
//	d := digest.Compute(ports)
//	if digest.Equal(d, previous) {
//	    // nothing changed — skip expensive diff
//	    return
//	}
//
// The Empty constant holds the digest of a nil/empty port slice and can be
// used as a safe initial value before the first scan completes.
package digest
