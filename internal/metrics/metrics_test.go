package metrics

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNew_InitializesCounters(t *testing.T) {
	m := New()
	c := m.Snapshot()
	if c.Scans != 0 || c.PortsFound != 0 || c.AlertsEmit != 0 || c.Errors != 0 {
		t.Fatalf("expected zero counters, got %+v", c)
	}
	if c.StartedAt.IsZero() {
		t.Fatal("StartedAt should be set")
	}
}

func TestRecordScan_IncrementsCounters(t *testing.T) {
	m := New()
	m.RecordScan(5)
	m.RecordScan(3)
	c := m.Snapshot()
	if c.Scans != 2 {
		t.Fatalf("expected 2 scans, got %d", c.Scans)
	}
	if c.PortsFound != 8 {
		t.Fatalf("expected 8 ports found, got %d", c.PortsFound)
	}
	if c.LastScan.IsZero() {
		t.Fatal("LastScan should be set after RecordScan")
	}
}

func TestRecordAlert_Increments(t *testing.T) {
	m := New()
	m.RecordAlert()
	m.RecordAlert()
	if got := m.Snapshot().AlertsEmit; got != 2 {
		t.Fatalf("expected 2 alerts, got %d", got)
	}
}

func TestRecordError_Increments(t *testing.T) {
	m := New()
	m.RecordError()
	if got := m.Snapshot().Errors; got != 1 {
		t.Fatalf("expected 1 error, got %d", got)
	}
}

func TestSnapshot_ReturnsCopy(t *testing.T) {
	m := New()
	m.RecordScan(10)
	a := m.Snapshot()
	m.RecordScan(10)
	b := m.Snapshot()
	if a.Scans == b.Scans {
		t.Fatal("snapshot should be a copy, not a live reference")
	}
}

func TestWrite_ContainsExpectedFields(t *testing.T) {
	m := New()
	m.RecordScan(4)
	m.RecordAlert()
	var buf bytes.Buffer
	m.Write(&buf)
	out := buf.String()
	for _, want := range []string{"scans=1", "ports_found=4", "alerts=1", "errors=0", "last_scan="} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q: %s", want, out)
		}
	}
}

func TestWrite_LastScanNeverWhenNoScans(t *testing.T) {
	m := New()
	var buf bytes.Buffer
	m.Write(&buf)
	if !strings.Contains(buf.String(), "last_scan=never") {
		t.Errorf("expected last_scan=never, got: %s", buf.String())
	}
}

func TestNew_StartedAtIsRecent(t *testing.T) {
	before := time.Now()
	m := New()
	after := time.Now()
	c := m.Snapshot()
	if c.StartedAt.Before(before) || c.StartedAt.After(after) {
		t.Fatalf("StartedAt %v not between %v and %v", c.StartedAt, before, after)
	}
}
