package history

import (
	"time"

	"github.com/user/portwatch/internal/state"
)

// Since returns all entries recorded at or after the given time.
func (h *History) Since(t time.Time) []Entry {
	var result []Entry
	for _, e := range h.Entries {
		if !e.Timestamp.Before(t) {
			result = append(result, e)
		}
	}
	return result
}

// PortSeen reports whether a port with the given number and protocol appeared
// in any historical entry.
func (h *History) PortSeen(number int, protocol string) bool {
	for _, e := range h.Entries {
		for _, p := range e.Ports {
			if p.Number == number && p.Protocol == protocol {
				return true
			}
		}
	}
	return false
}

// UniquePortsInRange returns deduplicated ports seen across all entries whose
// port numbers fall within [low, high] inclusive.
func (h *History) UniquePortsInRange(low, high int) []state.Port {
	seen := make(map[string]state.Port)
	for _, e := range h.Entries {
		for _, p := range e.Ports {
			if p.Number >= low && p.Number <= high {
				key := p.String()
				if _, ok := seen[key]; !ok {
					seen[key] = p
				}
			}
		}
	}
	result := make([]state.Port, 0, len(seen))
	for _, p := range seen {
		result = append(result, p)
	}
	return result
}
