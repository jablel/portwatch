package ratelimit_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/ratelimit"
)

// TestLimiter_SuppressesFlappingPort simulates a port that toggles state
// repeatedly and verifies that only the first event within each cooldown
// window passes through.
func TestLimiter_SuppressesFlappingPort(t *testing.T) {
	cooldown := 10 * time.Second
	l := ratelimit.New(cooldown)

	key := "added:9090/tcp"
	allowed := 0

	// Simulate 5 rapid events — only the first should pass.
	for i := 0; i < 5; i++ {
		if l.Allow(key) {
			allowed++
		}
	}

	if allowed != 1 {
		t.Fatalf("expected 1 allowed event, got %d", allowed)
	}
}

// TestLimiter_MultipleKeysDoNotInterfere ensures that rate limiting one key
// does not affect unrelated keys (different ports).
func TestLimiter_MultipleKeysDoNotInterfere(t *testing.T) {
	l := ratelimit.New(30 * time.Second)

	ports := []string{"80/tcp", "443/tcp", "22/tcp", "8080/tcp"}
	for _, p := range ports {
		if !l.Allow(p) {
			t.Errorf("expected first call for %s to be allowed", p)
		}
	}

	// Second pass — all should be blocked.
	for _, p := range ports {
		if l.Allow(p) {
			t.Errorf("expected second call for %s to be blocked", p)
		}
	}
}

// TestLimiter_ResetAllUnblocksEverything verifies that ResetAll allows
// all previously throttled keys through again.
func TestLimiter_ResetAllUnblocksEverything(t *testing.T) {
	l := ratelimit.New(60 * time.Second)
	keys := []string{"a", "b", "c"}
	for _, k := range keys {
		l.Allow(k)
	}
	l.ResetAll()
	for _, k := range keys {
		if !l.Allow(k) {
			t.Errorf("expected key %q to be allowed after ResetAll", k)
		}
	}
}

// TestLimiter_AllowsAfterCooldownExpires verifies that a key becomes
// allowed again once its cooldown window has elapsed.
func TestLimiter_AllowsAfterCooldownExpires(t *testing.T) {
	cooldown := 50 * time.Millisecond
	l := ratelimit.New(cooldown)

	key := "added:3000/tcp"

	if !l.Allow(key) {
		t.Fatal("expected first call to be allowed")
	}
	if l.Allow(key) {
		t.Fatal("expected second call to be blocked within cooldown")
	}

	time.Sleep(cooldown + 10*time.Millisecond)

	if !l.Allow(key) {
		t.Fatal("expected call to be allowed after cooldown expired")
	}
}
