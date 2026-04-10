package snapshot_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func makePorts(nums ...int) []scanner.Port {
	ports := make([]scanner.Port, len(nums))
	for i, n := range nums {
		ports[i] = scanner.Port{Number: n, Protocol: "tcp"}
	}
	return ports
}

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	store, err := snapshot.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := store.Save("baseline", makePorts(80, 443)); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "baseline.json")); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}

func TestLoad_ReturnsSavedPorts(t *testing.T) {
	dir := t.TempDir()
	store, _ := snapshot.New(dir)
	want := makePorts(22, 8080)
	_ = store.Save("test", want)

	snap, err := store.Load("test")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(snap.Ports) != len(want) {
		t.Errorf("got %d ports, want %d", len(snap.Ports), len(want))
	}
	if snap.Name != "test" {
		t.Errorf("got name %q, want %q", snap.Name, "test")
	}
	if snap.CapturedAt.IsZero() {
		t.Error("CapturedAt should not be zero")
	}
}

func TestLoad_MissingSnapshot(t *testing.T) {
	dir := t.TempDir()
	store, _ := snapshot.New(dir)
	_, err := store.Load("nonexistent")
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected ErrNotExist, got %v", err)
	}
}

func TestDelete_RemovesFile(t *testing.T) {
	dir := t.TempDir()
	store, _ := snapshot.New(dir)
	_ = store.Save("tmp", makePorts(9000))
	if err := store.Delete("tmp"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "tmp.json")); !os.IsNotExist(err) {
		t.Error("file should have been deleted")
	}
}

func TestDelete_MissingIsNoOp(t *testing.T) {
	dir := t.TempDir()
	store, _ := snapshot.New(dir)
	if err := store.Delete("ghost"); err != nil {
		t.Errorf("Delete missing: expected nil, got %v", err)
	}
}

func TestList_ReturnsNames(t *testing.T) {
	dir := t.TempDir()
	store, _ := snapshot.New(dir)
	_ = store.Save("alpha", makePorts(1))
	_ = store.Save("beta", makePorts(2))

	names, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("got %d names, want 2", len(names))
	}
}

func TestSave_OverwritesPrevious(t *testing.T) {
	dir := t.TempDir()
	store, _ := snapshot.New(dir)
	_ = store.Save("s", makePorts(80))
	_ = store.Save("s", makePorts(443))

	snap, _ := store.Load("s")
	if len(snap.Ports) != 1 || snap.Ports[0].Number != 443 {
		t.Errorf("expected overwritten snapshot with port 443, got %+v", snap.Ports)
	}
	_ = time.Now() // ensure time import used
}
