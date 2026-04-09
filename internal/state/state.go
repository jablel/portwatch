package state

import (
	"encoding/json"
	"os"

	"github.com/user/portwatch/internal/scanner"
)

// Diff holds the result of comparing two port snapshots.
type Diff struct {
	Added   []scanner.Port
	Removed []scanner.Port
}

// Snapshot represents a saved set of open ports.
type Snapshot struct {
	Ports []scanner.Port `json:"ports"`
}

// Compare returns the diff between previous and current port lists.
func Compare(previous, current []scanner.Port) Diff {
	prev := toSet(previous)
	curr := toSet(current)

	var added, removed []scanner.Port
	for _, p := range current {
		if !prev[p.String()] {
			added = append(added, p)
		}
	}
	for _, p := range previous {
		if !curr[p.String()] {
			removed = append(removed, p)
		}
	}
	return Diff{Added: added, Removed: removed}
}

// Save persists a port snapshot to the given file path.
func Save(path string, ports []scanner.Port) error {
	snap := Snapshot{Ports: ports}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Load reads a port snapshot from the given file path.
// Returns an empty snapshot if the file does not exist.
func Load(path string) ([]scanner.Port, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, err
	}
	return snap.Ports, nil
}

func toSet(ports []scanner.Port) map[string]bool {
	s := make(map[string]bool, len(ports))
	for _, p := range ports {
		s[p.String()] = true
	}
	return s
}
