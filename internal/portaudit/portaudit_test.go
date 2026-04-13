package portaudit_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/portaudit"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(num int, proto string) scanner.Port {
	return scanner.Port{Number: num, Protocol: proto}
}

func TestRecord_AppendsEntry(t *testing.T) {
	l := portaudit.New(0)
	l.Record(portaudit.KindAdded, makePort(80, "tcp"), "test")
	all := l.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(all))
	}
	if all[0].Kind != portaudit.KindAdded {
		t.Errorf("expected KindAdded, got %s", all[0].Kind)
	}
	if all[0].Port.Number != 80 {
		t.Errorf("expected port 80, got %d", all[0].Port.Number)
	}
}

func TestRecord_EvictsWhenFull(t *testing.T) {
	l := portaudit.New(2)
	l.Record(portaudit.KindAdded, makePort(80, "tcp"), "")
	l.Record(portaudit.KindAdded, makePort(443, "tcp"), "")
	l.Record(portaudit.KindRemoved, makePort(22, "tcp"), "")
	all := l.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries after eviction, got %d", len(all))
	}
	if all[0].Port.Number != 443 {
		t.Errorf("expected oldest evicted; first entry port 443, got %d", all[0].Port.Number)
	}
}

func TestSince_FiltersOldEntries(t *testing.T) {
	l := portaudit.New(0)
	l.Record(portaudit.KindAdded, makePort(80, "tcp"), "")
	cutoff := time.Now()
	l.Record(portaudit.KindAdded, makePort(443, "tcp"), "")
	result := l.Since(cutoff)
	if len(result) != 1 {
		t.Fatalf("expected 1 entry after cutoff, got %d", len(result))
	}
	if result[0].Port.Number != 443 {
		t.Errorf("expected port 443, got %d", result[0].Port.Number)
	}
}

func TestSince_ZeroCutoffReturnsAll(t *testing.T) {
	l := portaudit.New(0)
	l.Record(portaudit.KindAdded, makePort(80, "tcp"), "")
	l.Record(portaudit.KindRemoved, makePort(22, "tcp"), "")
	if got := l.Since(time.Time{}); len(got) != 2 {
		t.Errorf("expected 2 entries, got %d", len(got))
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.json")

	l := portaudit.New(0)
	l.Record(portaudit.KindAdded, makePort(8080, "tcp"), "ci")
	l.Record(portaudit.KindRemoved, makePort(9090, "udp"), "")

	if err := l.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	l2 := portaudit.New(0)
	if err := l2.Load(path); err != nil {
		t.Fatalf("Load: %v", err)
	}
	all := l2.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	if all[1].Actor != "" {
		t.Errorf("expected empty actor, got %q", all[1].Actor)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	l := portaudit.New(0)
	err := l.Load(filepath.Join(t.TempDir(), "nope.json"))
	if !os.IsNotExist(err) {
		t.Errorf("expected not-exist error, got %v", err)
	}
}
