// Package porttrend tracks how frequently a port appears across scan
// samples and exposes a simple trend (stable, rising, falling).
package porttrend

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Trend describes the direction of a port's observation frequency.
type Trend int

const (
	Stable  Trend = iota // seen consistently
	Rising               // appearing more often
	Falling              // appearing less often
)

func (t Trend) String() string {
	switch t {
	case Rising:
		return "rising"
	case Falling:
		return "falling"
	default:
		return "stable"
	}
}

type bucket struct {
	count int
	at    time.Time
}

// Tracker accumulates per-port observation counts over a sliding window
// and derives a trend by comparing the two halves of that window.
type Tracker struct {
	mu      sync.Mutex
	window  time.Duration
	buckets map[string][]bucket
}

// New returns a Tracker that uses the given window when computing trends.
func New(window time.Duration) *Tracker {
	return &Tracker{
		window:  window,
		buckets: make(map[string][]bucket),
	}
}

func portKey(p scanner.Port) string {
	return fmt.Sprintf("%s/%d", p.Protocol, p.Number)
}

// Record registers one observation of each port in the slice.
func (t *Tracker) Record(ports []scanner.Port) {
	now := time.Now()
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, p := range ports {
		k := portKey(p)
		t.buckets[k] = append(t.buckets[k], bucket{count: 1, at: now})
	}
	t.evict(now)
}

// Trend returns the current trend for a port.
func (t *Tracker) Trend(p scanner.Port) Trend {
	now := time.Now()
	t.mu.Lock()
	defer t.mu.Unlock()
	t.evict(now)

	bs := t.buckets[portKey(p)]
	if len(bs) == 0 {
		return Stable
	}
	mid := now.Add(-t.window / 2)
	var first, second int
	for _, b := range bs {
		if b.at.Before(mid) {
			first += b.count
		} else {
			second += b.count
		}
	}
	switch {
	case second > first:
		return Rising
	case second < first:
		return Falling
	default:
		return Stable
	}
}

// evict removes observations outside the current window. Must be called
// with t.mu held.
func (t *Tracker) evict(now time.Time) {
	cutoff := now.Add(-t.window)
	for k, bs := range t.buckets {
		i := 0
		for i < len(bs) && bs[i].at.Before(cutoff) {
			i++
		}
		if i == len(bs) {
			delete(t.buckets, k)
		} else {
			t.buckets[k] = bs[i:]
		}
	}
}
