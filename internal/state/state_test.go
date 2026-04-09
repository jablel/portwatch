package state_test

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

func makePorts(nums ...int) []scanner.Port {
	ports := make([]scanner.Port, len(nums))
	for i, n := range nums {
		ports[i] = scanner.Port{Number: n, Protocol: "tcp"}
	}
	return ports
}

func TestCompare_NoChanges(t *testing.T) {
	prev := makePorts(80, 443)
	curr := makePorts(80, 443)
	diff := state.Compare(prev, curr)
	if diff.HasChanges() {
		t.Errorf("expected no changes, got added=%v removed=%v", diff.Added, diff.Removed)
	}
}

func TestCompare_PortAdded(t *testing.T) {
	prev := makePorts(80)
	curr := makePorts(80, 8080)
	diff := state.Compare(prev, curr)
	if len(diff.Added) != 1 || diff.Added[0].Number != 8080 {
		t.Errorf("expected port 8080 added, got %v", diff.Added)
	}
	if len(diff.Removed) != 0 {
		t.Errorf("expected no removed ports, got %v", diff.Removed)
	}
}

func TestCompare_PortRemoved(t *testing.T) {
	prev := makePorts(80, 443)
	curr := makePorts(80)
	diff := state.Compare(prev, curr)
	if len(diff.Removed) != 1 || diff.Removed[0].Number != 443 {
		t.Errorf("expected port 443 removed, got %v", diff.Removed)
	}
	if len(diff.Added) != 0 {
		t.Errorf("expected no added ports, got %v", diff.Added)
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmp, err := os.CreateTemp("", "portwatch-state-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	ports := makePorts(22, 80, 443)
	if err := state.Save(tmp.Name(), ports); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	snap, err := state.Load(tmp.Name())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(snap.Ports) != len(ports) {
		t.Errorf("expected %d ports, got %d", len(ports), len(snap.Ports))
	}
	if snap.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := state.Load("/nonexistent/portwatch-state.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
