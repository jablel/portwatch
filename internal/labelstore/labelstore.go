// Package labelstore maintains a persistent map of port keys to
// user-defined string labels so that alerts and reports can show
// human-readable names alongside raw port numbers.
package labelstore

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Store maps port keys to labels.
type Store struct {
	mu     sync.RWMutex
	labels map[string]string
	path   string
}

// New returns an empty Store backed by the given file path.
func New(path string) *Store {
	return &Store{
		labels: make(map[string]string),
		path:   path,
	}
}

// key returns a canonical string for a Port.
func key(p scanner.Port) string {
	return fmt.Sprintf("%s:%d", p.Proto, p.Number)
}

// Set assigns a label to a port, overwriting any previous value.
func (s *Store) Set(p scanner.Port, label string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.labels[key(p)] = label
}

// Get returns the label for a port and whether it was found.
func (s *Store) Get(p scanner.Port) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.labels[key(p)]
	return v, ok
}

// Delete removes the label for a port.
func (s *Store) Delete(p scanner.Port) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.labels, key(p))
}

// Save persists the label map to disk as JSON.
func (s *Store) Save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, err := json.MarshalIndent(s.labels, "", "  ")
	if err != nil {
		return fmt.Errorf("labelstore: marshal: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0o644); err != nil {
		return fmt.Errorf("labelstore: write %s: %w", s.path, err)
	}
	return nil
}

// Load reads a previously saved label map from disk.
// If the file does not exist the store remains empty.
func (s *Store) Load() error {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("labelstore: read %s: %w", s.path, err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := json.Unmarshal(data, &s.labels); err != nil {
		return fmt.Errorf("labelstore: unmarshal: %w", err)
	}
	return nil
}
