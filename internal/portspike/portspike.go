// Package portspike detects sudden spikes in the number of open ports
// observed across consecutive scans. A spike is defined as an increase
// in port count that exceeds a configurable threshold ratio within a
// single scan interval.
package portspike

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Spike describes a detected port-count spike.
type Spike struct {
	Previous int
	Current  int
	Delta    int
	Ratio    float64
}

func (s Spike) String() string {
	return fmt.Sprintf("spike: %d -> %d (+%d, %.2fx)", s.Previous, s.Current, s.Delta, s.Ratio)
}

// Detector tracks port counts across scans and reports spikes.
type Detector struct {
	mu        sync.Mutex
	threshold float64 // minimum ratio increase to qualify as a spike
	prevCount int
	ready     bool
}

// New returns a Detector that fires when the port count grows by more
// than threshold as a ratio (e.g. 0.5 means a 50 % increase).
// A threshold <= 0 disables spike detection (Record always returns nil).
func New(threshold float64) *Detector {
	return &Detector{threshold: threshold}
}

// Record accepts the latest scan result and returns a non-nil Spike if
// the port count has grown beyond the configured threshold since the
// previous call. The first call always returns nil (no baseline yet).
func (d *Detector) Record(ports []scanner.Port) *Spike {
	d.mu.Lock()
	defer d.mu.Unlock()

	current := len(ports)

	if !d.ready {
		d.prevCount = current
		d.ready = true
		return nil
	}

	prev := d.prevCount
	d.prevCount = current

	if d.threshold <= 0 || prev == 0 {
		return nil
	}

	delta := current - prev
	if delta <= 0 {
		return nil
	}

	ratio := float64(delta) / float64(prev)
	if ratio < d.threshold {
		return nil
	}

	return &Spike{
		Previous: prev,
		Current:  current,
		Delta:    delta,
		Ratio:    ratio,
	}
}

// Reset clears the stored baseline so the next Record call is treated
// as the first observation.
func (d *Detector) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.ready = false
	d.prevCount = 0
}
