// Package aggregator combines scan results from multiple sources into a
// unified port list, deduplicating entries by (port, protocol) key.
package aggregator

import (
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Aggregator merges port slices from multiple scan sources.
type Aggregator struct {
	mu      sync.Mutex
	sources map[string][]scanner.Port
}

// New returns an empty Aggregator.
func New() *Aggregator {
	return &Aggregator{
		sources: make(map[string][]scanner.Port),
	}
}

// Update stores the latest scan result for the named source.
func (a *Aggregator) Update(source string, ports []scanner.Port) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.sources[source] = ports
}

// Merge returns a deduplicated, combined slice of all ports across sources.
func (a *Aggregator) Merge() []scanner.Port {
	a.mu.Lock()
	defer a.mu.Unlock()

	seen := make(map[string]struct{})
	var result []scanner.Port

	for _, ports := range a.sources {
		for _, p := range ports {
			key := p.String()
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			result = append(result, p)
		}
	}

	return result
}

// Sources returns the names of all registered sources.
func (a *Aggregator) Sources() []string {
	a.mu.Lock()
	defer a.mu.Unlock()

	names := make([]string, 0, len(a.sources))
	for k := range a.sources {
		names = append(names, k)
	}
	return names
}

// Remove deletes a source from the aggregator.
func (a *Aggregator) Remove(source string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.sources, source)
}
