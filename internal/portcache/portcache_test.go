package portcache

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makePorts(nums ...int) []scanner.Port {
	var out []scanner.Port
	for _, n := range nums {
		out = append(out, scanner.Port{Number: n, Protocol: "tcp"})
	}
	return out
}

func TestGet_EmptyCache(t *testing.T) {
	c := New(5 * time.Second)
	_, ok := c.Get()
	if ok {
		t.Fatal("expected empty cache to return false")
	}
}

func TestSet_ThenGet_ReturnsPorts(t *testing.T) {
	c := New(5 * time.Second)
	ports := makePorts(80, 443)
	c.Set(ports)
	entry, ok := c.Get()
	if !ok {
		t.Fatal("expected entry to be present")
	}
	if len(entry.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(entry.Ports))
	}
}

func TestGet_ExpiredEntry_ReturnsFalse(t *testing.T) {
	c := New(10 * time.Millisecond)
	now := time.Now()
	c.now = func() time.Time { return now }
	c.Set(makePorts(22))

	// advance clock past TTL
	c.now = func() time.Time { return now.Add(20 * time.Millisecond) }
	_, ok := c.Get()
	if ok {
		t.Fatal("expected expired entry to return false")
	}
}

func TestGet_ZeroTTL_NeverExpires(t *testing.T) {
	c := New(0)
	now := time.Now()
	c.now = func() time.Time { return now }
	c.Set(makePorts(8080))

	c.now = func() time.Time { return now.Add(24 * time.Hour) }
	_, ok := c.Get()
	if !ok {
		t.Fatal("expected zero-TTL cache to never expire")
	}
}

func TestInvalidate_ClearsEntry(t *testing.T) {
	c := New(5 * time.Second)
	c.Set(makePorts(3306))
	c.Invalidate()
	_, ok := c.Get()
	if ok {
		t.Fatal("expected cache to be empty after Invalidate")
	}
}

func TestAge_EmptyCacheReturnsNegativeOne(t *testing.T) {
	c := New(5 * time.Second)
	if c.Age() != -1 {
		t.Fatal("expected -1 for empty cache")
	}
}

func TestAge_ReturnsTimeSinceSet(t *testing.T) {
	c := New(5 * time.Second)
	now := time.Now()
	c.now = func() time.Time { return now }
	c.Set(makePorts(9200))
	c.now = func() time.Time { return now.Add(2 * time.Second) }
	if got := c.Age(); got != 2*time.Second {
		t.Fatalf("expected 2s age, got %v", got)
	}
}

func TestSet_MutationOfOriginalSlice_DoesNotAffectCache(t *testing.T) {
	c := New(5 * time.Second)
	ports := makePorts(80)
	c.Set(ports)
	ports[0].Number = 9999
	entry, _ := c.Get()
	if entry.Ports[0].Number == 9999 {
		t.Fatal("cache should store a copy, not a reference")
	}
}
