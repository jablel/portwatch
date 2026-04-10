// Package metrics tracks runtime counters for portwatch scans.
package metrics

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// Counters holds cumulative statistics gathered during daemon operation.
type Counters struct {
	mu          sync.Mutex
	Scans       int
	PortsFound  int
	AlertsEmit  int
	Errors      int
	LastScan    time.Time
	StartedAt   time.Time
}

// Metrics wraps Counters and provides thread-safe update methods.
type Metrics struct {
	c Counters
}

// New returns a new Metrics instance with StartedAt set to now.
func New() *Metrics {
	return &Metrics{
		c: Counters{StartedAt: time.Now()},
	}
}

// RecordScan increments the scan counter and updates LastScan.
func (m *Metrics) RecordScan(portsFound int) {
	m.c.mu.Lock()
	defer m.c.mu.Unlock()
	m.c.Scans++
	m.c.PortsFound += portsFound
	m.c.LastScan = time.Now()
}

// RecordAlert increments the alerts-emitted counter.
func (m *Metrics) RecordAlert() {
	m.c.mu.Lock()
	defer m.c.mu.Unlock()
	m.c.AlertsEmit++
}

// RecordError increments the error counter.
func (m *Metrics) RecordError() {
	m.c.mu.Lock()
	defer m.c.mu.Unlock()
	m.c.Errors++
}

// Snapshot returns a copy of the current counters.
func (m *Metrics) Snapshot() Counters {
	m.c.mu.Lock()
	defer m.c.mu.Unlock()
	return m.c
}

// Write prints a human-readable summary to w.
func (m *Metrics) Write(w io.Writer) {
	c := m.Snapshot()
	uptime := time.Since(c.StartedAt).Round(time.Second)
	last := "never"
	if !c.LastScan.IsZero() {
		last = c.LastScan.Format(time.RFC3339)
	}
	fmt.Fprintf(w, "uptime=%s scans=%d ports_found=%d alerts=%d errors=%d last_scan=%s\n",
		uptime, c.Scans, c.PortsFound, c.AlertsEmit, c.Errors, last)
}
