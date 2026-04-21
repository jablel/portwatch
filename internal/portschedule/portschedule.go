// Package portschedule tracks when ports are expected to be active based on
// observed time-of-day patterns, flagging ports seen outside their normal schedule.
package portschedule

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// hourBucket is a bitmask of 24 bits, one per hour of the day.
type hourBucket uint32

func (h hourBucket) set(hour int) hourBucket   { return h | (1 << uint(hour)) }
func (h hourBucket) isset(hour int) bool        { return h&(1<<uint(hour)) != 0 }

// Violation describes a port seen outside its expected schedule.
type Violation struct {
	Port     scanner.Port
	Hour     int
	Expected string // human-readable expected window
}

func (v Violation) Error() string {
	return fmt.Sprintf("port %s active at hour %02d:00 outside expected schedule (%s)",
		v.Port, v.Hour, v.Expected)
}

type entry struct {
	bucket hourBucket
	count  int
}

// Tracker learns per-port active hours and detects schedule violations.
type Tracker struct {
	mu        sync.Mutex
	schedules map[string]*entry
	minObs    int // minimum observations before schedule is enforced
}

// New returns a Tracker that requires at least minObservations scans before
// enforcing a schedule for any port.
func New(minObservations int) *Tracker {
	if minObservations < 1 {
		minObservations = 1
	}
	return &Tracker{
		schedules: make(map[string]*entry),
		minObs:    minObservations,
	}
}

func portKey(p scanner.Port) string {
	return fmt.Sprintf("%d/%s", p.Number, p.Protocol)
}

// Observe records that the given ports were seen at time t and returns any
// violations for ports active outside their learned schedule.
func (t *Tracker) Observe(ports []scanner.Port, at time.Time) []Violation {
	hour := at.Hour()
	t.mu.Lock()
	defer t.mu.Unlock()

	var violations []Violation
	for _, p := range ports {
		k := portKey(p)
		e, ok := t.schedules[k]
		if !ok {
			e = &entry{}
			t.schedules[k] = e
		}
		if e.count < t.minObs {
			e.bucket = e.bucket.set(hour)
			e.count++
			continue
		}
		if !e.bucket.isset(hour) {
			violations = append(violations, Violation{
				Port:     p,
				Hour:     hour,
				Expected: formatBucket(e.bucket),
			})
		}
	}
	return violations
}

// Reset clears all learned schedules.
func (t *Tracker) Reset() {
	t.mu.Lock()
	t.schedules = make(map[string]*entry)
	t.mu.Unlock()
}

func formatBucket(b hourBucket) string {
	var hours []int
	for h := 0; h < 24; h++ {
		if b.isset(h) {
			hours = append(hours, h)
		}
	}
	if len(hours) == 0 {
		return "none"
	}
	return fmt.Sprintf("hours %v", hours)
}
