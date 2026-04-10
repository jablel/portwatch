// Package snapshot provides point-in-time captures of open ports,
// allowing portwatch to compare current state against a named reference.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot holds a named capture of open ports at a specific time.
type Snapshot struct {
	Name      string         `json:"name"`
	CapturedAt time.Time     `json:"captured_at"`
	Ports     []scanner.Port `json:"ports"`
}

// Store manages named snapshots persisted to a directory.
type Store struct {
	dir string
}

// New returns a Store rooted at dir, creating it if necessary.
func New(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("snapshot: create dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Save writes a snapshot under the given name, overwriting any previous one.
func (s *Store) Save(name string, ports []scanner.Port) error {
	snap := Snapshot{
		Name:       name,
		CapturedAt: time.Now().UTC(),
		Ports:      ports,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	return os.WriteFile(s.path(name), data, 0o644)
}

// Load retrieves a previously saved snapshot by name.
// Returns os.ErrNotExist if no snapshot with that name exists.
func (s *Store) Load(name string) (*Snapshot, error) {
	data, err := os.ReadFile(s.path(name))
	if err != nil {
		return nil, fmt.Errorf("snapshot: load %q: %w", name, err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return &snap, nil
}

// Delete removes a named snapshot. Returns nil if it did not exist.
func (s *Store) Delete(name string) error {
	err := os.Remove(s.path(name))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// List returns the names of all stored snapshots.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, fmt.Errorf("snapshot: list: %w", err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".json" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}

func (s *Store) path(name string) string {
	return filepath.Join(s.dir, name+".json")
}
