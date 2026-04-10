// Package digest computes a stable fingerprint for a set of open ports,
// allowing quick equality checks without a full diff.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"portwatch/internal/scanner"
)

// Digest is a hex-encoded SHA-256 fingerprint of a port set.
type Digest string

// Empty is the digest of an empty port set.
const Empty Digest = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

// Compute returns a deterministic Digest for the given ports.
// Ports are sorted by protocol+number before hashing so that insertion
// order does not affect the result.
func Compute(ports []scanner.Port) Digest {
	if len(ports) == 0 {
		return Empty
	}

	keys := make([]string, len(ports))
	for i, p := range ports {
		keys[i] = p.String()
	}
	sort.Strings(keys)

	h := sha256.New()
	for _, k := range keys {
		fmt.Fprintln(h, k)
	}

	return Digest(hex.EncodeToString(h.Sum(nil)))
}

// Equal reports whether two digests are identical.
func Equal(a, b Digest) bool {
	return a == b
}

// String implements fmt.Stringer.
func (d Digest) String() string {
	return string(d)
}
