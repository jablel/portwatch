package eventlog

import (
	"time"

	"github.com/user/portwatch/internal/state"
)

// ForPort returns all events matching the given port number and protocol.
func (l *EventLog) ForPort(port state.Port) []Event {
	l.mu.RLock()
	defer l.mu.RUnlock()
	var out []Event
	for _, e := range l.events {
		if e.Port.Number == port.Number && e.Port.Protocol == port.Protocol {
			out = append(out, e)
		}
	}
	return out
}

// CountByKind returns the number of events matching the given kind ("added" or "removed").
func (l *EventLog) CountByKind(kind string) int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	n := 0
	for _, e := range l.events {
		if e.Kind == kind {
			n++
		}
	}
	return n
}

// Between returns events whose timestamps fall within [from, to] inclusive.
func (l *EventLog) Between(from, to time.Time) []Event {
	l.mu.RLock()
	defer l.mu.RUnlock()
	var out []Event
	for _, e := range l.events {
		if !e.Timestamp.Before(from) && !e.Timestamp.After(to) {
			out = append(out, e)
		}
	}
	return out
}

// Latest returns the most recent event, or nil if the log is empty.
func (l *EventLog) Latest() *Event {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if len(l.events) == 0 {
		return nil
	}
	e := l.events[len(l.events)-1]
	return &e
}
