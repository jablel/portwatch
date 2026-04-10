// Package baseline manages the trusted set of ports that define
// the expected state of the system. Deviations from the baseline
// trigger alerts in the daemon loop.
package baseline

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// ErrNoBaseline is returned when no baseline file exists on disk.
var ErrNoBaseline = errors.New("baseline: no baseline file found")

// Baseline holds a trusted snapshot of open ports.
type Baseline struct {
	Ports     []scanner.Port `json:"ports"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// New creates an empty Baseline.
func New() *Baseline {
	now := time.Now().UTC()
	return &Baseline{
		Ports:     []scanner.Port{},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Set replaces the trusted port list and updates the timestamp.
func (b *Baseline) Set(ports []scanner.Port) {
	b.Ports = ports
	b.UpdatedAt = time.Now().UTC()
}

// Save writes the baseline to the given file path as JSON.
func (b *Baseline) Save(path string) error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// Load reads a baseline from the given file path.
// Returns ErrNoBaseline if the file does not exist.
func Load(path string) (*Baseline, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNoBaseline
	}
	if err != nil {
		return nil, err
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// Contains reports whether the given port is part of the baseline.
func (b *Baseline) Contains(p scanner.Port) bool {
	for _, bp := range b.Ports {
		if bp.Number == p.Number && bp.Protocol == p.Protocol {
			return true
		}
	}
	return false
}
