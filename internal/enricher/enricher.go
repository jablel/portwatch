// Package enricher attaches human-readable metadata to scanned ports,
// combining tagger labels with well-known service names.
package enricher

import (
	"fmt"
	"strconv"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/tagger"
)

// EnrichedPort wraps a scanner.Port with additional metadata.
type EnrichedPort struct {
	scanner.Port
	ServiceName string
	Tags        []string
}

// String returns a human-readable representation of the enriched port.
func (e EnrichedPort) String() string {
	if e.ServiceName != "" {
		return fmt.Sprintf("%s (%s/%s)", e.ServiceName, strconv.Itoa(int(e.Port.Number)), e.Port.Protocol)
	}
	return e.Port.String()
}

// Enricher annotates ports with tags and service names.
type Enricher struct {
	tagger *tagger.Tagger
}

// New creates an Enricher backed by the given Tagger.
func New(t *tagger.Tagger) *Enricher {
	return &Enricher{tagger: t}
}

// Enrich converts a slice of scanner.Port into EnrichedPort values.
func (e *Enricher) Enrich(ports []scanner.Port) []EnrichedPort {
	out := make([]EnrichedPort, 0, len(ports))
	for _, p := range ports {
		tagged := e.tagger.Tag(p)
		out = append(out, EnrichedPort{
			Port:        p,
			ServiceName: tagged.Label,
			Tags:        tagged.Tags,
		})
	}
	return out
}

// EnrichOne annotates a single port.
func (e *Enricher) EnrichOne(p scanner.Port) EnrichedPort {
	tagged := e.tagger.Tag(p)
	return EnrichedPort{
		Port:        p,
		ServiceName: tagged.Label,
		Tags:        tagged.Tags,
	}
}
