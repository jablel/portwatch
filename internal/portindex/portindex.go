// Package portindex maintains an in-memory index of observed ports,
// allowing fast lookup by port number, protocol, or tag.
package portindex

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Index holds a thread-safe, multi-key index over a set of ports.
type Index struct {
	mu      sync.RWMutex
	byNum   map[int][]scanner.Port
	byProto map[string][]scanner.Port
	byTag   map[string][]scanner.Port
}

// New returns an empty Index.
func New() *Index {
	return &Index{
		byNum:   make(map[int][]scanner.Port),
		byProto: make(map[string][]scanner.Port),
		byTag:   make(map[string][]scanner.Port),
	}
}

// Build replaces the entire index with the provided ports.
func (idx *Index) Build(ports []scanner.Port) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	idx.byNum = make(map[int][]scanner.Port, len(ports))
	idx.byProto = make(map[string][]scanner.Port)
	idx.byTag = make(map[string][]scanner.Port)

	for _, p := range ports {
		idx.byNum[p.Number] = append(idx.byNum[p.Number], p)
		idx.byProto[p.Protocol] = append(idx.byProto[p.Protocol], p)
		for _, t := range p.Tags {
			idx.byTag[t] = append(idx.byTag[t], p)
		}
	}
}

// ByNumber returns all ports matching the given port number.
func (idx *Index) ByNumber(n int) []scanner.Port {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return cloneSlice(idx.byNum[n])
}

// ByProtocol returns all ports matching the given protocol string.
func (idx *Index) ByProtocol(proto string) []scanner.Port {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return cloneSlice(idx.byProto[proto])
}

// ByTag returns all ports that carry the given tag.
func (idx *Index) ByTag(tag string) []scanner.Port {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return cloneSlice(idx.byTag[tag])
}

// Size returns the total number of unique (number, protocol) pairs indexed.
func (idx *Index) Size() int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	total := 0
	for _, v := range idx.byNum {
		total += len(v)
	}
	return total
}

// Key returns a canonical string key for a port.
func Key(p scanner.Port) string {
	return fmt.Sprintf("%s:%d", p.Protocol, p.Number)
}

func cloneSlice(src []scanner.Port) []scanner.Port {
	if len(src) == 0 {
		return nil
	}
	out := make([]scanner.Port, len(src))
	copy(out, src)
	return out
}
