package history_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
)

func TestSince_FiltersOldEntries(t *testing.T) {
	h := history.New(10)
	h.Add(makePorts(22))
	cutoff := time.Now().UTC()
	time.Sleep(2 * time.Millisecond)
	h.Add(makePorts(80))

	result := h.Since(cutoff)
	if len(result) != 1 {
		t.Fatalf("expected 1 entry after cutoff, got %d", len(result))
	}
	if result[0].Ports[0].Number != 80 {
		t.Errorf("expected port 80, got %d", result[0].Ports[0].Number)
	}
}

func TestSince_ReturnsAllWhenCutoffIsZero(t *testing.T) {
	h := history.New(10)
	h.Add(makePorts(22))
	h.Add(makePorts(80))

	result := h.Since(time.Time{})
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestPortSeen_ReturnsTrueWhenPresent(t *testing.T) {
	h := history.New(10)
	h.Add(makePorts(443))

	if !h.PortSeen(443, "tcp") {
		t.Error("expected port 443/tcp to be seen")
	}
}

func TestPortSeen_ReturnsFalseWhenAbsent(t *testing.T) {
	h := history.New(10)
	h.Add(makePorts(80))

	if h.PortSeen(9999, "tcp") {
		t.Error("expected port 9999/tcp to not be seen")
	}
}

func TestUniquePortsInRange_DeduplicatesAcrossEntries(t *testing.T) {
	h := history.New(10)
	h.Add(makePorts(80, 443))
	h.Add(makePorts(80, 8080))

	ports := h.UniquePortsInRange(80, 1000)
	if len(ports) != 3 {
		t.Errorf("expected 3 unique ports, got %d", len(ports))
	}
}

func TestUniquePortsInRange_ExcludesOutOfRange(t *testing.T) {
	h := history.New(10)
	h.Add(makePorts(22, 80, 443))

	ports := h.UniquePortsInRange(80, 443)
	for _, p := range ports {
		if p.Number < 80 || p.Number > 443 {
			t.Errorf("port %d is outside range [80, 443]", p.Number)
		}
	}
	if len(ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(ports))
	}
}
