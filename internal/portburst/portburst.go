// Package portburst detects short-lived bursts of new ports appearing within
// a sliding time window. A burst is declared when the number of distinct ports
// added within the window exceeds a configurable threshold.
package portburst

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Burst describes a detected burst event.
type Burst struct {
	At    time.Time
	Ports []scanner.Port
	Count int
}

// Detector tracks port additions and signals bursts.
type Detector struct {
	mu        sync.Mutex
	window    time.Duration
	threshold int
	events    []event
}

type event struct {
	at   time.Time
	port scanner.Port
}

// New returns a Detector that fires when more than threshold distinct ports
// are added within window.
func New(window time.Duration, threshold int) *Detector {
	return &Detector{
		window:    window,
		threshold: threshold,
	}
}

// Record ingests newly added ports and returns a Burst if the threshold is
// exceeded within the current window, or nil otherwise.
func (d *Detector) Record(added []scanner.Port) *Burst {
	if len(added) == 0 {
		return nil
	}

	now := time.Now()

	d.mu.Lock()
	defer d.mu.Unlock()

	// Append new events.
	for _, p := range added {
		d.events = append(d.events, event{at: now, port: p})
	}

	// Evict events outside the window.
	cutoff := now.Add(-d.window)
	start := 0
	for start < len(d.events) && d.events[start].at.Before(cutoff) {
		start++
	}
	d.events = d.events[start:]

	if d.threshold <= 0 || len(d.events) <= d.threshold {
		return nil
	}

	ports := make([]scanner.Port, len(d.events))
	for i, e := range d.events {
		ports[i] = e.port
	}
	return &Burst{At: now, Ports: ports, Count: len(ports)}
}

// Reset clears all recorded events.
func (d *Detector) Reset() {
	d.mu.Lock()
	d.events = d.events[:0]
	d.mu.Unlock()
}
