package ratelimit

import (
	"testing"
	"time"
)

// newFakeLimiter returns a Limiter whose clock can be controlled via the
// returned pointer.
func newFakeLimiter(cooldown time.Duration) (*Limiter, *time.Time) {
	t := time.Now()
	l := New(cooldown)
	l.now = func() time.Time { return t }
	return l, &t
}

func TestAllow_FirstCallAlwaysPasses(t *testing.T) {
	l, _ := newFakeLimiter(5 * time.Second)
	if !l.Allow("tcp:8080") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallWithinCooldownBlocked(t *testing.T) {
	l, _ := newFakeLimiter(5 * time.Second)
	l.Allow("tcp:8080")
	if l.Allow("tcp:8080") {
		t.Fatal("expected second call within cooldown to be blocked")
	}
}

func TestAllow_SecondCallAfterCooldownPasses(t *testing.T) {
	l, clock := newFakeLimiter(5 * time.Second)
	l.Allow("tcp:8080")
	*clock = clock.Add(6 * time.Second)
	if !l.Allow("tcp:8080") {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestAllow_ZeroCooldownAlwaysPasses(t *testing.T) {
	l := New(0)
	for i := 0; i < 5; i++ {
		if !l.Allow("tcp:9000") {
			t.Fatalf("call %d: expected zero-cooldown to always pass", i)
		}
	}
}

func TestAllow_NegativeCooldownAlwaysPasses(t *testing.T) {
	l := New(-1 * time.Second)
	l.Allow("k")
	if !l.Allow("k") {
		t.Fatal("expected negative cooldown to always pass")
	}
}

func TestReset_UnblocksKey(t *testing.T) {
	l, _ := newFakeLimiter(10 * time.Second)
	l.Allow("tcp:443")
	l.Reset("tcp:443")
	if !l.Allow("tcp:443") {
		t.Fatal("expected Reset to unblock the key")
	}
}

func TestResetAll_ClearsAllKeys(t *testing.T) {
	l, _ := newFakeLimiter(10 * time.Second)
	l.Allow("a")
	l.Allow("b")
	l.Allow("c")
	l.ResetAll()
	if l.Len() != 0 {
		t.Fatalf("expected 0 keys after ResetAll, got %d", l.Len())
	}
}

func TestLen_TracksActiveKeys(t *testing.T) {
	l, _ := newFakeLimiter(10 * time.Second)
	l.Allow("x")
	l.Allow("y")
	if got := l.Len(); got != 2 {
		t.Fatalf("expected 2 tracked keys, got %d", got)
	}
}
