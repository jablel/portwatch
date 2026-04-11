package eventlog

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/state"
)

func TestForPort_ReturnsMatchingEvents(t *testing.T) {
	l := New(20)
	now := time.Now()
	l.Append(makeEvent("added", 80, "tcp", now))
	l.Append(makeEvent("added", 443, "tcp", now))
	l.Append(makeEvent("removed", 80, "tcp", now.Add(time.Minute)))
	got := l.ForPort(state.Port{Number: 80, Protocol: "tcp"})
	if len(got) != 2 {
		t.Fatalf("expected 2 events for port 80, got %d", len(got))
	}
}

func TestForPort_ProtocolMismatch_ReturnsNone(t *testing.T) {
	l := New(10)
	l.Append(makeEvent("added", 53, "tcp", time.Now()))
	got := l.ForPort(state.Port{Number: 53, Protocol: "udp"})
	if len(got) != 0 {
		t.Errorf("expected 0 events for udp/53, got %d", len(got))
	}
}

func TestCountByKind(t *testing.T) {
	l := New(20)
	now := time.Now()
	l.Append(makeEvent("added", 80, "tcp", now))
	l.Append(makeEvent("added", 443, "tcp", now))
	l.Append(makeEvent("removed", 8080, "tcp", now))
	if got := l.CountByKind("added"); got != 2 {
		t.Errorf("expected 2 added, got %d", got)
	}
	if got := l.CountByKind("removed"); got != 1 {
		t.Errorf("expected 1 removed, got %d", got)
	}
}

func TestBetween_ReturnsEventsInRange(t *testing.T) {
	l := New(10)
	now := time.Now()
	l.Append(makeEvent("added", 80, "tcp", now.Add(-3*time.Hour)))
	l.Append(makeEvent("added", 443, "tcp", now.Add(-1*time.Hour)))
	l.Append(makeEvent("added", 8080, "tcp", now))
	from := now.Add(-2 * time.Hour)
	to := now.Add(-30 * time.Minute)
	got := l.Between(from, to)
	if len(got) != 1 {
		t.Fatalf("expected 1 event in range, got %d", len(got))
	}
	if got[0].Port.Number != 443 {
		t.Errorf("expected port 443, got %d", got[0].Port.Number)
	}
}

func TestLatest_ReturnsNilWhenEmpty(t *testing.T) {
	l := New(10)
	if l.Latest() != nil {
		t.Error("expected nil from empty log")
	}
}

func TestLatest_ReturnsMostRecent(t *testing.T) {
	l := New(10)
	now := time.Now()
	l.Append(makeEvent("added", 80, "tcp", now))
	l.Append(makeEvent("added", 9000, "tcp", now.Add(time.Second)))
	if got := l.Latest(); got == nil || got.Port.Number != 9000 {
		t.Errorf("expected latest port 9000, got %+v", got)
	}
}
