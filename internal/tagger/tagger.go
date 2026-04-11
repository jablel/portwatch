// Package tagger assigns human-readable labels to ports based on
// well-known service mappings and user-defined rules.
package tagger

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// wellKnown maps port numbers to common service names.
var wellKnown = map[uint16]string{
	21:   "ftp",
	22:   "ssh",
	25:   "smtp",
	53:   "dns",
	80:   "http",
	110:  "pop3",
	143:  "imap",
	443:  "https",
	3306: "mysql",
	5432: "postgres",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	27017: "mongodb",
}

// Tagger assigns labels to ports.
type Tagger struct {
	mu      sync.RWMutex
	custom  map[uint16]string
}

// New returns a Tagger with no custom rules.
func New() *Tagger {
	return &Tagger{
		custom: make(map[uint16]string),
	}
}

// Define registers a custom label for a port number, overriding well-known
// mappings.
func (t *Tagger) Define(port uint16, label string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.custom[port] = label
}

// Tag returns the label for the given port. Custom rules take precedence over
// well-known names. If no mapping exists the label is "unknown:PORT".
func (t *Tagger) Tag(p scanner.Port) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if label, ok := t.custom[p.Number]; ok {
		return label
	}
	if label, ok := wellKnown[p.Number]; ok {
		return label
	}
	return fmt.Sprintf("unknown:%d", p.Number)
}

// TagAll annotates a slice of ports, returning a map from port to label.
func (t *Tagger) TagAll(ports []scanner.Port) map[scanner.Port]string {
	out := make(map[scanner.Port]string, len(ports))
	for _, p := range ports {
		out[p] = t.Tag(p)
	}
	return out
}
