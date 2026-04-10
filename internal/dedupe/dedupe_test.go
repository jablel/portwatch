package dedupe_test

import (
	"testing"
	"time"

	"portwatch/internal/dedupe"
	"portwatch/internal/scanner"
)

func makePort(proto string, number int) scanner.Port {
	return scanner.Port{Protocol: proto, Number: number}
}

func TestIsDuplicate_FirstCallNotDuplicate(t *testing.T) {
	d := dedupe.New(10 * time.Second)
	p := makePort("tcp", 8080)
	if d.IsDuplicate(p) {
		t.Fatal("expected first call to not be a duplicate")
	}
}

func TestIsDuplicate_SecondCallWithinWindowIsDuplicate(t *testing.T) {
	d := dedupe.New(10 * time.Second)
	p := makePort("tcp", 8080)
	d.IsDuplicate(p)
	if !d.IsDuplicate(p) {
		t.Fatal("expected second call within window to be a duplicate")
	}
}

func TestIsDuplicate_CallAfterWindowNotDuplicate(t *testing.T) {
	now := time.Now()
	d := dedupe.New(5 * time.Second)

	// Inject a fake clock that advances beyond the window.
	type clockSetter interface {
		SetNow(func() time.Time)
	}

	// Use the exported nowFunc indirectly by manipulating time via Reset + re-add.
	p := makePort("tcp", 9090)
	d.IsDuplicate(p) // record at "now"

	// Simulate time passing by using a fresh deduper with a tiny window
	// and verifying Evict clears the entry.
	d2 := dedupe.New(1 * time.Millisecond)
	d2.IsDuplicate(p)
	time.Sleep(5 * time.Millisecond)
	d2.Evict()
	if d2.IsDuplicate(p) {
		t.Fatal("expected entry to be cleared after eviction")
	}
	_ = now
}

func TestIsDuplicate_DifferentPortsAreIndependent(t *testing.T) {
	d := dedupe.New(10 * time.Second)
	p1 := makePort("tcp", 80)
	p2 := makePort("tcp", 443)

	d.IsDuplicate(p1)
	if d.IsDuplicate(p2) {
		t.Fatal("expected different port to not be a duplicate")
	}
}

func TestIsDuplicate_ZeroWindowNeverDedupes(t *testing.T) {
	d := dedupe.New(0)
	p := makePort("udp", 53)
	d.IsDuplicate(p)
	if d.IsDuplicate(p) {
		t.Fatal("expected zero window to never suppress")
	}
}

func TestReset_ClearsAllEntries(t *testing.T) {
	d := dedupe.New(10 * time.Second)
	p := makePort("tcp", 22)
	d.IsDuplicate(p)
	d.Reset()
	if d.IsDuplicate(p) {
		t.Fatal("expected Reset to clear all entries")
	}
}
