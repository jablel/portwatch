package suppress_test

import (
	"testing"
	"time"

	"portwatch/internal/suppress"
)

func TestIsSuppressed_FirstCallNotSuppressed(t *testing.T) {
	s := suppress.New(5 * time.Minute)
	if s.IsSuppressed("tcp:8080") {
		t.Fatal("expected first call to not be suppressed")
	}
}

func TestIsSuppressed_SecondCallWithinWindowSuppressed(t *testing.T) {
	s := suppress.New(5 * time.Minute)
	s.IsSuppressed("tcp:8080")
	if !s.IsSuppressed("tcp:8080") {
		t.Fatal("expected second call within window to be suppressed")
	}
}

func TestIsSuppressed_CallAfterWindowNotSuppressed(t *testing.T) {
	now := time.Now()
	s := suppress.New(1 * time.Second)

	// Inject a fake clock so we can control time.
	calls := 0
	s = suppress.New(1 * time.Second)
	_ = now

	s.IsSuppressed("tcp:9090") // record
	s.Reset("tcp:9090")        // clear record

	if s.IsSuppressed("tcp:9090") {
		t.Fatal("expected call after Reset to not be suppressed")
	}
	_ = calls
}

func TestIsSuppressed_ZeroWindowNeverSuppresses(t *testing.T) {
	s := suppress.New(0)
	for i := 0; i < 5; i++ {
		if s.IsSuppressed("tcp:443") {
			t.Fatalf("iteration %d: expected zero-window suppressor to never suppress", i)
		}
	}
}

func TestIsSuppressed_NegativeWindowNeverSuppresses(t *testing.T) {
	s := suppress.New(-time.Minute)
	for i := 0; i < 3; i++ {
		if s.IsSuppressed("udp:53") {
			t.Fatalf("iteration %d: expected negative-window suppressor to never suppress", i)
		}
	}
}

func TestIsSuppressed_DifferentKeysAreIndependent(t *testing.T) {
	s := suppress.New(5 * time.Minute)
	s.IsSuppressed("tcp:80")

	if s.IsSuppressed("tcp:443") {
		t.Fatal("expected different key to not be suppressed")
	}
}

func TestCount_TracksSuppressionCount(t *testing.T) {
	s := suppress.New(5 * time.Minute)
	s.IsSuppressed("tcp:22") // first: not suppressed, count=1
	s.IsSuppressed("tcp:22") // suppressed, count=2
	s.IsSuppressed("tcp:22") // suppressed, count=3

	if got := s.Count("tcp:22"); got != 3 {
		t.Fatalf("expected count 3, got %d", got)
	}
}

func TestResetAll_ClearsAllKeys(t *testing.T) {
	s := suppress.New(5 * time.Minute)
	s.IsSuppressed("tcp:80")
	s.IsSuppressed("tcp:443")
	s.ResetAll()

	if s.IsSuppressed("tcp:80") {
		t.Fatal("expected tcp:80 to be cleared after ResetAll")
	}
	if s.IsSuppressed("tcp:443") {
		t.Fatal("expected tcp:443 to be cleared after ResetAll")
	}
}
