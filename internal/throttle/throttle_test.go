package throttle_test

import (
	"testing"
	"time"

	"portwatch/internal/throttle"
)

func TestAllow_FirstCallAlwaysPasses(t *testing.T) {
	th := throttle.New(time.Second)
	if !th.Allow("port:tcp:8080") {
		t.Fatal("expected first call to Allow to return true")
	}
}

func TestAllow_SecondCallWithinIntervalBlocked(t *testing.T) {
	th := throttle.New(time.Hour)
	th.Allow("port:tcp:8080")
	if th.Allow("port:tcp:8080") {
		t.Fatal("expected second call within interval to return false")
	}
}

func TestAllow_DifferentKeysAreIndependent(t *testing.T) {
	th := throttle.New(time.Hour)
	th.Allow("port:tcp:8080")
	if !th.Allow("port:tcp:9090") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestAllow_ZeroIntervalAlwaysPasses(t *testing.T) {
	th := throttle.New(0)
	for i := 0; i < 5; i++ {
		if !th.Allow("port:tcp:8080") {
			t.Fatalf("expected zero-interval throttle to always allow (iteration %d)", i)
		}
	}
}

func TestAllow_NegativeIntervalAlwaysPasses(t *testing.T) {
	th := throttle.New(-time.Second)
	if !th.Allow("key") || !th.Allow("key") {
		t.Fatal("expected negative interval to always allow")
	}
}

func TestReset_AllowsImmediateRetry(t *testing.T) {
	th := throttle.New(time.Hour)
	th.Allow("port:tcp:8080")
	th.Reset("port:tcp:8080")
	if !th.Allow("port:tcp:8080") {
		t.Fatal("expected Allow to return true after Reset")
	}
}

func TestLen_TracksKeys(t *testing.T) {
	th := throttle.New(time.Hour)
	if th.Len() != 0 {
		t.Fatalf("expected 0 keys, got %d", th.Len())
	}
	th.Allow("a")
	th.Allow("b")
	if th.Len() != 2 {
		t.Fatalf("expected 2 keys, got %d", th.Len())
	}
	th.Reset("a")
	if th.Len() != 1 {
		t.Fatalf("expected 1 key after reset, got %d", th.Len())
	}
}
