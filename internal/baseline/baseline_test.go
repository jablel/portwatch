package baseline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/scanner"
)

func makePorts(nums ...int) []scanner.Port {
	ports := make([]scanner.Port, len(nums))
	for i, n := range nums {
		ports[i] = scanner.Port{Number: n, Protocol: "tcp"}
	}
	return ports
}

func TestNew_IsEmpty(t *testing.T) {
	b := baseline.New()
	if len(b.Ports) != 0 {
		t.Fatalf("expected empty ports, got %d", len(b.Ports))
	}
	if b.CreatedAt.IsZero() {
		t.Fatal("expected non-zero CreatedAt")
	}
}

func TestSet_UpdatesPorts(t *testing.T) {
	b := baseline.New()
	before := b.UpdatedAt
	time.Sleep(2 * time.Millisecond)
	b.Set(makePorts(80, 443))
	if len(b.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(b.Ports))
	}
	if !b.UpdatedAt.After(before) {
		t.Fatal("expected UpdatedAt to advance after Set")
	}
}

func TestContains_ReturnsTrueForKnownPort(t *testing.T) {
	b := baseline.New()
	b.Set(makePorts(22, 80))
	p := scanner.Port{Number: 80, Protocol: "tcp"}
	if !b.Contains(p) {
		t.Fatal("expected Contains to return true for port 80")
	}
}

func TestContains_ReturnsFalseForUnknownPort(t *testing.T) {
	b := baseline.New()
	b.Set(makePorts(22))
	p := scanner.Port{Number: 9999, Protocol: "tcp"}
	if b.Contains(p) {
		t.Fatal("expected Contains to return false for port 9999")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	b := baseline.New()
	b.Set(makePorts(22, 80, 443))
	if err := b.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded.Ports) != 3 {
		t.Fatalf("expected 3 ports, got %d", len(loaded.Ports))
	}
}

func TestLoad_MissingFile_ReturnsErrNoBaseline(t *testing.T) {
	_, err := baseline.Load(filepath.Join(t.TempDir(), "missing.json"))
	if err != baseline.ErrNoBaseline {
		t.Fatalf("expected ErrNoBaseline, got %v", err)
	}
}

func TestLoad_CorruptFile_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o600)
	_, err := baseline.Load(path)
	if err == nil {
		t.Fatal("expected error for corrupt JSON")
	}
}
