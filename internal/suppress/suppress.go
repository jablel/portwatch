// Package suppress provides a mechanism to suppress repeated alerts
// for ports that have already been reported within a configurable window.
package suppress

import (
	"sync"
	"time"
)

// Entry records when a key was last suppressed.
type Entry struct {
	LastSeen time.Time
	Count    int
}

// Suppressor tracks seen keys and suppresses duplicates within a window.
type Suppressor struct {
	mu     sync.Mutex
	window time.Duration
	seen   map[string]*Entry
	now    func() time.Time
}

// New returns a Suppressor that silences repeated events within window.
// A zero or negative window means nothing is suppressed.
func New(window time.Duration) *Suppressor {
	return &Suppressor{
		window: window,
		seen:   make(map[string]*Entry),
		now:    time.Now,
	}
}

// IsSuppressed reports whether key should be suppressed.
// If not suppressed, the key is recorded and future calls within the
// window will return true.
func (s *Suppressor) IsSuppressed(key string) bool {
	if s.window <= 0 {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	now := s.now()

	if e, ok := s.seen[key]; ok {
		if now.Sub(e.LastSeen) < s.window {
			e.Count++
			return true
		}
	}

	s.seen[key] = &Entry{LastSeen: now, Count: 1}
	return false
}

// Reset clears the suppression record for key, allowing the next event
// through regardless of the window.
func (s *Suppressor) Reset(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.seen, key)
}

// ResetAll clears all suppression records.
func (s *Suppressor) ResetAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seen = make(map[string]*Entry)
}

// Count returns how many times key has been suppressed since it was first seen.
// Returns 0 if the key is unknown or has expired.
func (s *Suppressor) Count(key string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.seen[key]
	if !ok {
		return 0
	}
	if s.window > 0 && s.now().Sub(e.LastSeen) >= s.window {
		return 0
	}
	return e.Count
}
