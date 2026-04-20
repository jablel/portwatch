// Package portanomaly detects statistically anomalous port activity by
// comparing recent observation frequency against a rolling baseline.
package portanomaly

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Anomaly describes a port whose recent frequency deviates significantly from
// its historical baseline.
type Anomaly struct {
	Port      scanner.Port
	Baseline  float64 // historical average frequency [0,1]
	Recent    float64 // recent frequency [0,1]
	Deviation float64 // Recent - Baseline
}

func (a Anomaly) String() string {
	return fmt.Sprintf("%s baseline=%.2f recent=%.2f deviation=%+.2f",
		a.Port, a.Baseline, a.Recent, a.Deviation)
}

// Detector tracks per-port frequencies across two rolling windows and surfaces
// ports whose recent activity deviates beyond a configurable threshold.
type Detector struct {
	mu        sync.Mutex
	threshold float64
	baseline  map[string][]float64 // older window samples
	recent    map[string][]float64 // newer window samples
	window    int                  // samples per half-window
}

// New returns a Detector. threshold is the minimum absolute deviation required
// to flag a port (e.g. 0.3 means recent must differ from baseline by ≥0.30).
// window is the number of scans in each half of the comparison window.
func New(threshold float64, window int) *Detector {
	if window < 1 {
		window = 1
	}
	return &Detector{
		threshold: threshold,
		window:    window,
		baseline:  make(map[string][]float64),
		recent:    make(map[string][]float64),
	}
}

// Record ingests one scan observation. present is the set of ports seen in the
// scan; all tracked ports not in present receive a frequency sample of 0.
func (d *Detector) Record(present []scanner.Port) {
	d.mu.Lock()
	defer d.mu.Unlock()

	seen := make(map[string]bool, len(present))
	for _, p := range present {
		seen[portKey(p)] = true
	}

	// Ensure every known port gets a sample this round.
	keys := make(map[string]bool)
	for k := range d.recent {
		keys[k] = true
	}
	for k := range seen {
		keys[k] = true
	}

	for k := range keys {
		var v float64
		if seen[k] {
			v = 1
		}
		d.recent[k] = append(d.recent[k], v)
		if len(d.recent[k]) > d.window {
			// Promote oldest recent sample into baseline.
			d.baseline[k] = append(d.baseline[k], d.recent[k][0])
			d.recent[k] = d.recent[k][1:]
			if len(d.baseline[k]) > d.window {
				d.baseline[k] = d.baseline[k][1:]
			}
		}
	}
}

// Anomalies returns all ports whose deviation exceeds the configured threshold.
func (d *Detector) Anomalies(ports []scanner.Port) []Anomaly {
	d.mu.Lock()
	defer d.mu.Unlock()

	var out []Anomaly
	for _, p := range ports {
		k := portKey(p)
		bl := d.baseline[k]
		rc := d.recent[k]
		if len(bl) == 0 || len(rc) == 0 {
			continue
		}
		bAvg := avg(bl)
		rAvg := avg(rc)
		dev := rAvg - bAvg
		if abs(dev) >= d.threshold {
			out = append(out, Anomaly{
				Port:      p,
				Baseline:  bAvg,
				Recent:    rAvg,
				Deviation: dev,
			})
		}
	}
	return out
}

func portKey(p scanner.Port) string {
	return fmt.Sprintf("%d/%s", p.Number, p.Protocol)
}

func avg(s []float64) float64 {
	if len(s) == 0 {
		return 0
	}
	var sum float64
	for _, v := range s {
		sum += v
	}
	return sum / float64(len(s))
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
