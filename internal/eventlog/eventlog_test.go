package eventlog

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/state"
)

func makeEvent(kind string, number int, protocol string, t time.Time) Event {
	return Event{
		Timestamp: t,
		Kind:      kind,
		Port:      state.Port{Number: number, Protocol: protocol},
	}
}

func TestAppend_AddsEvents(t *testing.T) {
	l := New(10)
	l.Append(makeEvent("added", 80, "tcp", time.Now()))
	l.Append(makeEvent("removed", 443, "tcp", time.Now()))
	if got := len(l.All()); got != 2 {
		t.Fatalf("expected 2 events, got %d", got)
	}
}

func TestAppend_EvictsOldestWhenFull(t *testing.T) {
	l := New(3)
	base := time.Now()
	for i := 0; i < 4; i++ {
		l.Append(makeEvent("added", 8000+i, "tcp", base.Add(time.Duration(i)*time.Second)))
	}
	all := l.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 events after eviction, got %d", len(all))
	}
	if all[0].Port.Number != 8001 {
		t.Errorf("expected oldest surviving port 8001, got %d", all[0].Port.Number)
	}
}

func TestSince_FiltersEvents(t *testing.T) {
	l := New(10)
	now := time.Now()
	l.Append(makeEvent("added", 80, "tcp", now.Add(-2*time.Hour)))
	l.Append(makeEvent("added", 443, "tcp", now.Add(-30*time.Minute)))
	l.Append(makeEvent("added", 8080, "tcp", now))
	got := l.Since(now.Add(-1 * time.Hour))
	if len(got) != 2 {
		t.Fatalf("expected 2 events, got %d", len(got))
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "events.jsonl")
	l := New(10)
	now := time.Now().UTC().Truncate(time.Millisecond)
	l.Append(makeEvent("added", 22, "tcp", now))
	l.Append(makeEvent("removed", 3306, "tcp", now.Add(time.Minute)))
	if err := l.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	l2, err := Load(path, 10)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	all := l2.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 events after load, got %d", len(all))
	}
	if all[0].Port.Number != 22 || all[1].Port.Number != 3306 {
		t.Errorf("unexpected ports after round-trip: %+v", all)
	}
}

func TestLoad_MissingFile_ReturnsEmpty(t *testing.T) {
	l, err := Load(filepath.Join(os.TempDir(), "no_such_file.jsonl"), 10)
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(l.All()) != 0 {
		t.Errorf("expected empty log, got %d events", len(l.All()))
	}
}
