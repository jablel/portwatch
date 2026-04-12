// Package profiler tracks per-port scan timing and latency statistics.
package profiler

import (
	"sync"
	"time"
)

// Sample holds timing data for a single scan of a port key.
type Sample struct {
	Key      string
	Duration time.Duration
	RecordedAt time.Time
}

// Stats holds aggregated latency statistics for a port key.
type Stats struct {
	Key    string
	Count  int64
	Total  time.Duration
	Min    time.Duration
	Max    time.Duration
}

// Mean returns the average scan duration.
func (s Stats) Mean() time.Duration {
	if s.Count == 0 {
		return 0
	}
	return s.Total / time.Duration(s.Count)
}

// Profiler records scan durations and exposes aggregated statistics.
type Profiler struct {
	mu    sync.Mutex
	data  map[string]*Stats
}

// New returns an initialised Profiler.
func New() *Profiler {
	return &Profiler{data: make(map[string]*Stats)}
}

// Record adds a timing sample for the given key.
func (p *Profiler) Record(s Sample) {
	p.mu.Lock()
	defer p.mu.Unlock()

	st, ok := p.data[s.Key]
	if !ok {
		st = &Stats{Key: s.Key, Min: s.Duration, Max: s.Duration}
		p.data[s.Key] = st
	}
	st.Count++
	st.Total += s.Duration
	if s.Duration < st.Min {
		st.Min = s.Duration
	}
	if s.Duration > st.Max {
		st.Max = s.Duration
	}
}

// Get returns a copy of the Stats for key, and whether it exists.
func (p *Profiler) Get(key string) (Stats, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	st, ok := p.data[key]
	if !ok {
		return Stats{}, false
	}
	return *st, true
}

// All returns a snapshot of all recorded statistics.
func (p *Profiler) All() []Stats {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]Stats, 0, len(p.data))
	for _, st := range p.data {
		out = append(out, *st)
	}
	return out
}

// Reset clears all recorded data.
func (p *Profiler) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.data = make(map[string]*Stats)
}
