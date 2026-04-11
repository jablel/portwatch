// Package eventlog provides a structured, append-only log of port change events
// with support for querying by time range and port.
package eventlog

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Event represents a single port change event.
type Event struct {
	Timestamp time.Time  `json:"timestamp"`
	Kind      string     `json:"kind"` // "added" or "removed"
	Port      state.Port `json:"port"`
}

// EventLog is a thread-safe, bounded append-only log of events.
type EventLog struct {
	mu      sync.RWMutex
	events  []Event
	maxSize int
}

// New creates a new EventLog with the given maximum size.
func New(maxSize int) *EventLog {
	if maxSize <= 0 {
		maxSize = 1000
	}
	return &EventLog{maxSize: maxSize}
}

// Append adds one or more events to the log, evicting the oldest when full.
func (l *EventLog) Append(events ...Event) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, e := range events {
		if len(l.events) >= l.maxSize {
			l.events = l.events[1:]
		}
		l.events = append(l.events, e)
	}
}

// All returns a copy of all events in the log.
func (l *EventLog) All() []Event {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Event, len(l.events))
	copy(out, l.events)
	return out
}

// Since returns events with a timestamp at or after cutoff.
func (l *EventLog) Since(cutoff time.Time) []Event {
	l.mu.RLock()
	defer l.mu.RUnlock()
	var out []Event
	for _, e := range l.events {
		if !e.Timestamp.Before(cutoff) {
			out = append(out, e)
		}
	}
	return out
}

// Save persists the log to the given file path as newline-delimited JSON.
func (l *EventLog) Save(path string) error {
	l.mu.RLock()
	defer l.mu.RUnlock()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, e := range l.events {
		if err := enc.Encode(e); err != nil {
			return err
		}
	}
	return nil
}

// Load reads a newline-delimited JSON file into the log, appending to existing entries.
func Load(path string, maxSize int) (*EventLog, error) {
	l := New(maxSize)
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return l, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Event
		if err := dec.Decode(&e); err != nil {
			return nil, err
		}
		l.events = append(l.events, e)
	}
	return l, nil
}
