package history_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/state"
)

func makePorts(numbers ...int) []state.Port {
	ports := make([]state.Port, len(numbers))
	for i, n := range numbers {
		ports[i] = state.Port{Number: n, Protocol: "tcp"}
	}
	return ports
}

func TestAdd_AppendsEntry(t *testing.T) {
	h := history.New(10)
	h.Add(makePorts(80, 443))
	if len(h.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(h.Entries))
	}
	if h.Entries[0].Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestAdd_EvictsOldestWhenFull(t *testing.T) {
	h := history.New(3)
	for i := 0; i < 5; i++ {
		h.Add(makePorts(i))
	}
	if len(h.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(h.Entries))
	}
	// Oldest retained entry should have port 2
	if h.Entries[0].Ports[0].Number != 2 {
		t.Errorf("expected port 2, got %d", h.Entries[0].Ports[0].Number)
	}
}

func TestLatest_ReturnsNilWhenEmpty(t *testing.T) {
	h := history.New(10)
	if h.Latest() != nil {
		t.Error("expected nil for empty history")
	}
}

func TestLatest_ReturnsMostRecent(t *testing.T) {
	h := history.New(10)
	h.Add(makePorts(80))
	time.Sleep(time.Millisecond)
	h.Add(makePorts(443))
	latest := h.Latest()
	if latest == nil {
		t.Fatal("expected non-nil latest")
	}
	if latest.Ports[0].Number != 443 {
		t.Errorf("expected port 443, got %d", latest.Ports[0].Number)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	h := history.New(10)
	h.Add(makePorts(22, 80))
	h.Add(makePorts(443))

	if err := h.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := history.Load(path, 10)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(loaded.Entries))
	}
}

func TestLoad_MissingFile_ReturnsEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "no-such-file.json")
	h, err := history.Load(path, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(h.Entries) != 0 {
		t.Errorf("expected empty history, got %d entries", len(h.Entries))
	}
}

func TestLoad_CorruptFile_ReturnsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	_ = os.WriteFile(path, []byte("not json{"), 0o644)
	_, err := history.Load(path, 10)
	if err == nil {
		t.Error("expected error for corrupt file")
	}
}
