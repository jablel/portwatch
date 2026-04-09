package state

import (
	"encoding/json"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot represents a saved state of open ports at a point in time.
type Snapshot struct {
	Timestamp time.Time      `json:"timestamp"`
	Ports     []scanner.Port `json:"ports"`
}

// Diff holds the ports that were added or removed between two snapshots.
type Diff struct {
	Added   []scanner.Port
	Removed []scanner.Port
}

// HasChanges returns true if any ports were added or removed.
func (d Diff) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0
}

// Compare returns a Diff between a previous and current set of ports.
func Compare(previous, current []scanner.Port) Diff {
	prevSet := toSet(previous)
	currSet := toSet(current)

	var added, removed []scanner.Port

	for _, p := range current {
		if !prevSet[p.String()] {
			added = append(added, p)
		}
	}
	for _, p := range previous {
		if !currSet[p.String()] {
			removed = append(removed, p)
		}
	}

	return Diff{Added: added, Removed: removed}
}

// Save writes a snapshot to the given file path as JSON.
func Save(path string, ports []scanner.Port) error {
	snap := Snapshot{
		Timestamp: time.Now(),
		Ports:     ports,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, err
	}
	return &snap, nil
}

func toSet(ports []scanner.Port) map[string]bool {
	set := make(map[string]bool, len(ports))
	for _, p := range ports {
		set[p.String()] = true
	}
	return set
}
