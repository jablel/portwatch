// Package healthcheck provides a simple liveness probe for the portwatch daemon.
// It tracks the last successful scan time and exposes a status summary.
package healthcheck

import (
	"fmt"
	"sync"
	"time"
)

// Status represents the current health of the daemon.
type Status struct {
	Healthy      bool
	LastScan     time.Time
	ScanCount    int64
	ErrorCount   int64
	Uptime       time.Duration
}

// Checker tracks daemon liveness.
type Checker struct {
	mu         sync.RWMutex
	startedAt  time.Time
	lastScan   time.Time
	scanCount  int64
	errorCount int64
	maxStaleness time.Duration
}

// New creates a Checker with the given max staleness threshold.
// If the last scan is older than maxStaleness, the daemon is considered unhealthy.
func New(maxStaleness time.Duration) *Checker {
	return &Checker{
		startedAt:    time.Now(),
		maxStaleness: maxStaleness,
	}
}

// RecordScan marks a successful scan tick.
func (c *Checker) RecordScan() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastScan = time.Now()
	c.scanCount++
}

// RecordError increments the error counter.
func (c *Checker) RecordError() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.errorCount++
}

// Status returns a snapshot of the current health status.
func (c *Checker) Status() Status {
	c.mu.RLock()
	defer c.mu.RUnlock()

	healthy := !c.lastScan.IsZero() &&
		time.Since(c.lastScan) <= c.maxStaleness

	return Status{
		Healthy:    healthy,
		LastScan:   c.lastScan,
		ScanCount:  c.scanCount,
		ErrorCount: c.errorCount,
		Uptime:     time.Since(c.startedAt),
	}
}

// String returns a human-readable summary.
func (s Status) String() string {
	state := "healthy"
	if !s.Healthy {
		state = "unhealthy"
	}
	return fmt.Sprintf("status=%s scans=%d errors=%d uptime=%s",
		state, s.ScanCount, s.ErrorCount, s.Uptime.Round(time.Second))
}
