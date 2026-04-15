// Package portmap provides a bidirectional mapping between port numbers and
// their canonical service names, enabling fast lookups in both directions.
package portmap

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// PortMap holds bidirectional mappings between port numbers and service names.
type PortMap struct {
	mu          sync.RWMutex
	byNumber    map[string]string // "proto:number" -> service name
	byService   map[string][]scanner.Port // service name -> ports
}

// New returns an empty PortMap.
func New() *PortMap {
	return &PortMap{
		byNumber:  make(map[string]string),
		byService: make(map[string][]scanner.Port),
	}
}

func portKey(p scanner.Port) string {
	return fmt.Sprintf("%s:%d", p.Proto, p.Number)
}

// Register associates a port with a service name.
// Registering the same port twice overwrites the previous entry.
func (m *PortMap) Register(p scanner.Port, service string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	k := portKey(p)
	// Remove from old service list if previously registered.
	if old, ok := m.byNumber[k]; ok && old != service {
		filtered := m.byService[old][:0]
		for _, existing := range m.byService[old] {
			if portKey(existing) != k {
				filtered = append(filtered, existing)
			}
		}
		m.byService[old] = filtered
	}

	m.byNumber[k] = service
	m.byService[service] = append(m.byService[service], p)
}

// LookupService returns the service name for the given port, and whether it was found.
func (m *PortMap) LookupService(p scanner.Port) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	svc, ok := m.byNumber[portKey(p)]
	return svc, ok
}

// LookupPorts returns all ports registered under the given service name.
func (m *PortMap) LookupPorts(service string) []scanner.Port {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]scanner.Port, len(m.byService[service]))
	copy(result, m.byService[service])
	return result
}

// Len returns the total number of registered port entries.
func (m *PortMap) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.byNumber)
}
