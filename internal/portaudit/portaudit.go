// Package portaudit records and queries a tamper-evident audit trail of
// port-change events, associating each entry with a wall-clock timestamp
// and an optional actor/source label.
package portaudit

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Kind classifies the type of audit event.
type Kind string

const (
	KindAdded   Kind = "added"
	KindRemoved Kind = "removed"
)

// Entry is a single audit record.
type Entry struct {
	Timestamp time.Time    `json:"timestamp"`
	Kind      Kind         `json:"kind"`
	Port      scanner.Port `json:"port"`
	Actor     string       `json:"actor,omitempty"`
}

// Log holds an in-memory audit trail with optional persistence.
type Log struct {
	mu      sync.RWMutex
	entries []Entry
	maxSize int
}

// New creates a Log that retains at most maxSize entries (0 = unlimited).
func New(maxSize int) *Log {
	return &Log{maxSize: maxSize}
}

// Record appends a new entry to the log.
func (l *Log) Record(kind Kind, port scanner.Port, actor string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	e := Entry{Timestamp: time.Now(), Kind: kind, Port: port, Actor: actor}
	l.entries = append(l.entries, e)
	if l.maxSize > 0 && len(l.entries) > l.maxSize {
		l.entries = l.entries[len(l.entries)-l.maxSize:]
	}
}

// All returns a copy of all entries in chronological order.
func (l *Log) All() []Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Entry, len(l.entries))
	copy(out, l.entries)
	return out
}

// Since returns entries recorded at or after cutoff.
func (l *Log) Since(cutoff time.Time) []Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	var out []Entry
	for _, e := range l.entries {
		if !e.Timestamp.Before(cutoff) {
			out = append(out, e)
		}
	}
	return out
}

// Save writes the current log to path as JSON.
func (l *Log) Save(path string) error {
	l.mu.RLock()
	defer l.mu.RUnlock()
	data, err := json.MarshalIndent(l.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// Load replaces the in-memory log with entries read from path.
func (l *Log) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = entries
	return nil
}
