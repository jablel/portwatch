package ratelimit

import (
	"testing"
	"time"
)

func newFakeLimiter(cooldown time.Duration) (*Limiter, *time.Time) {
	current := time.Now()
	l := New(cooldown)
	l.now = func() time.Time { return current }
	return l, &current
}

func TestAllow_FirstCallAlwaysPasses(t *testing.T) {
	l, _ := newFakeLimiter(5 * time.Second)
	if !l.Allow("port:80/tcp") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallWithinCooldownBlocked(t *testing.T) {
	l, _ := newFakeLimiter(5 * time.Second)
	l.Allow("port:80/tcp")
	if l.Allow("port:80/tcp") {
		t.Fatal("expected second call within cooldown to be blocked")
	}
}

func TestAllow_SecondCallAfterCooldownPasses(t *testing.T) {
	l, ts := newFakeLimiter(5 * time.Second)
	l.Allow("port:80/tcp")
	*ts = ts.Add(6 * time.Second)
	if !l.Allow("port:80/tcp") {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestAllow_ZeroCooldownAlwaysPasses(t *testing.T) {
	l, _ := newFakeLimiter(0)
	l.Allow("k")
	if !l.Allow("k") {
		t.Fatal("expected zero cooldown to always allow")
	}
}

func TestAllow_NegativeCooldownAlwaysPasses(t *testing.T) {
	l, _ := newFakeLimiter(-time.Second)
	l.Allow("k")
	if !l.Allow("k") {
		t.Fatal("expected negative cooldown to always allow")
	}
}

func TestAllow_DifferentKeysAreIndependent(t *testing.T) {
	l, _ := newFakeLimiter(5 * time.Second)
	l.Allow("port:80/tcp")
	if !l.Allow("port:443/tcp") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestReset_AllowsKeyImmediately(t *testing.T) {
	l, _ := newFakeLimiter(5 * time.Second)
	l.Allow("port:80/tcp")
	l.Reset("port:80/tcp")
	if !l.Allow("port:80/tcp") {
		t.Fatal("expected key to be allowed after reset")
	}
}

func TestRemaining_ReturnsZeroWhenNotSeen(t *testing.T) {
	l, _ := newFakeLimiter(5 * time.Second)
	if r := l.Remaining("port:80/tcp"); r != 0 {
		t.Fatalf("expected 0, got %v", r)
	}
}

func TestRemaining_ReturnsPositiveWithinCooldown(t *testing.T) {
	l, ts := newFakeLimiter(10 * time.Second)
	l.Allow("port:80/tcp")
	*ts = ts.Add(3 * time.Second)
	if r := l.Remaining("port:80/tcp"); r != 7*time.Second {
		t.Fatalf("expected 7s remaining, got %v", r)
	}
}

func TestResetAll_ClearsAllKeys(t *testing.T) {
	l, _ := newFakeLimiter(5 * time.Second)
	l.Allow("a")
	l.Allow("b")
	l.ResetAll()
	if !l.Allow("a") || !l.Allow("b") {
		t.Fatal("expected all keys to be cleared after ResetAll")
	}
}
