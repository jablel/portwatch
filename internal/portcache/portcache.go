// Package portcache provides a time-bounded in-memory cache for port scan
// results, allowing consumers to retrieve the most recent scan without
// triggering a new one.
package portcache

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry holds a cached set of ports together with the time it was stored.
type Entry struct {
	Ports     []scanner.Port
	StoredAt  time.Time
}

// Cache stores the latest port scan result and expires it after a configurable
// TTL. A zero TTL means entries never expire.
type Cache struct {
	mu    sync.RWMutex
	entry *Entry
	ttl   time.Duration
	now   func() time.Time
}

// New returns a Cache with the given TTL. Pass 0 for no expiry.
func New(ttl time.Duration) *Cache {
	return &Cache{ttl: ttl, now: time.Now}
}

// Set stores ports as the current cached value.
func (c *Cache) Set(ports []scanner.Port) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cp := make([]scanner.Port, len(ports))
	copy(cp, ports)
	c.entry = &Entry{Ports: cp, StoredAt: c.now()}
}

// Get returns the cached entry and true when a valid (non-expired) entry
// exists, otherwise it returns nil and false.
func (c *Cache) Get() (*Entry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.entry == nil {
		return nil, false
	}
	if c.ttl > 0 && c.now().Sub(c.entry.StoredAt) > c.ttl {
		return nil, false
	}
	return c.entry, true
}

// Invalidate removes the current cached entry.
func (c *Cache) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entry = nil
}

// Age returns how long ago the entry was stored. It returns -1 when the cache
// is empty.
func (c *Cache) Age() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.entry == nil {
		return -1
	}
	return c.now().Sub(c.entry.StoredAt)
}
