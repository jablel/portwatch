// Package sampler provides periodic port scan sampling with configurable
// intervals, storing lightweight snapshots for trend analysis.
package sampler

import (
	"context"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Sample holds the result of a single scan at a point in time.
type Sample struct {
	At    time.Time
	Ports []scanner.Port
}

// Sampler collects periodic scan samples up to a fixed capacity,
// evicting the oldest entry when full.
type Sampler struct {
	mu       sync.Mutex
	scanner  *scanner.Scanner
	interval time.Duration
	cap      int
	samples  []Sample
}

// New creates a Sampler that scans with the given scanner every interval,
// retaining at most capacity samples.
func New(s *scanner.Scanner, interval time.Duration, capacity int) *Sampler {
	if capacity <= 0 {
		capacity = 1
	}
	return &Sampler{
		scanner:  s,
		interval: interval,
		cap:      capacity,
	}
}

// Run starts the sampling loop. It blocks until ctx is cancelled.
func (s *Sampler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			ports, err := s.scanner.Scan(ctx)
			if err != nil {
				continue
			}
			s.record(t, ports)
		}
	}
}

// record appends a sample, evicting the oldest when at capacity.
func (s *Sampler) record(at time.Time, ports []scanner.Port) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.samples) >= s.cap {
		s.samples = s.samples[1:]
	}
	s.samples = append(s.samples, Sample{At: at, Ports: ports})
}

// All returns a copy of all retained samples, oldest first.
func (s *Sampler) All() []Sample {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Sample, len(s.samples))
	copy(out, s.samples)
	return out
}

// Latest returns the most recent sample, or nil if none exist.
func (s *Sampler) Latest() *Sample {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.samples) == 0 {
		return nil
	}
	copy := s.samples[len(s.samples)-1]
	return &copy
}
