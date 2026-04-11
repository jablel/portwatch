package healthcheck_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/healthcheck"
)

func TestStatus_UnhealthyBeforeFirstScan(t *testing.T) {
	c := healthcheck.New(5 * time.Second)
	s := c.Status()
	if s.Healthy {
		t.Fatal("expected unhealthy before first scan")
	}
	if s.ScanCount != 0 {
		t.Fatalf("expected 0 scans, got %d", s.ScanCount)
	}
}

func TestStatus_HealthyAfterScan(t *testing.T) {
	c := healthcheck.New(5 * time.Second)
	c.RecordScan()
	s := c.Status()
	if !s.Healthy {
		t.Fatal("expected healthy after scan")
	}
	if s.ScanCount != 1 {
		t.Fatalf("expected 1 scan, got %d", s.ScanCount)
	}
}

func TestStatus_UnhealthyWhenStale(t *testing.T) {
	c := healthcheck.New(1 * time.Millisecond)
	c.RecordScan()
	time.Sleep(10 * time.Millisecond)
	s := c.Status()
	if s.Healthy {
		t.Fatal("expected unhealthy after staleness threshold exceeded")
	}
}

func TestRecordError_IncrementsCounter(t *testing.T) {
	c := healthcheck.New(5 * time.Second)
	c.RecordError()
	c.RecordError()
	s := c.Status()
	if s.ErrorCount != 2 {
		t.Fatalf("expected 2 errors, got %d", s.ErrorCount)
	}
}

func TestStatus_UptimeGrowsOverTime(t *testing.T) {
	c := healthcheck.New(5 * time.Second)
	time.Sleep(5 * time.Millisecond)
	s := c.Status()
	if s.Uptime < 5*time.Millisecond {
		t.Fatalf("expected uptime >= 5ms, got %s", s.Uptime)
	}
}

func TestStatus_String_Healthy(t *testing.T) {
	c := healthcheck.New(5 * time.Second)
	c.RecordScan()
	s := c.Status()
	got := s.String()
	if got == "" {
		t.Fatal("expected non-empty string")
	}
	if s.Healthy && len(got) == 0 {
		t.Fatal("string should not be empty")
	}
}

func TestStatus_String_Unhealthy(t *testing.T) {
	c := healthcheck.New(5 * time.Second)
	s := c.Status()
	got := s.String()
	if got == "" {
		t.Fatal("expected non-empty string for unhealthy status")
	}
}
