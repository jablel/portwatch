package windowcounter

import (
	"testing"
	"time"
)

func TestAdd_FirstCallReturnsOne(t *testing.T) {
	c := New(time.Second)
	if got := c.Add("tcp:80"); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestAdd_AccumulatesWithinWindow(t *testing.T) {
	c := New(time.Second)
	c.Add("tcp:80")
	c.Add("tcp:80")
	if got := c.Add("tcp:80"); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestAdd_DifferentKeysAreIndependent(t *testing.T) {
	c := New(time.Second)
	c.Add("tcp:80")
	c.Add("tcp:80")
	if got := c.Add("udp:53"); got != 1 {
		t.Fatalf("expected 1 for udp:53, got %d", got)
	}
}

func TestCount_DoesNotAddEvent(t *testing.T) {
	c := New(time.Second)
	c.Add("tcp:443")
	if got := c.Count("tcp:443"); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
	if got := c.Count("tcp:443"); got != 1 {
		t.Fatalf("expected still 1, got %d", got)
	}
}

func TestCount_UnknownKeyReturnsZero(t *testing.T) {
	c := New(time.Second)
	if got := c.Count("tcp:9999"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestReset_ClearsKey(t *testing.T) {
	c := New(time.Second)
	c.Add("tcp:80")
	c.Add("tcp:80")
	c.Reset("tcp:80")
	if got := c.Count("tcp:80"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestAdd_ZeroWindowCountsAll(t *testing.T) {
	c := New(0)
	for i := 0; i < 5; i++ {
		c.Add("tcp:22")
	}
	if got := c.Count("tcp:22"); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestAdd_ExpiredEntriesEvicted(t *testing.T) {
	c := New(50 * time.Millisecond)
	c.Add("tcp:8080")
	c.Add("tcp:8080")
	time.Sleep(80 * time.Millisecond)
	if got := c.Add("tcp:8080"); got != 1 {
		t.Fatalf("expected 1 after window expired, got %d", got)
	}
}
