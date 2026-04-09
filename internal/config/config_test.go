package config

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestDefault_Values(t *testing.T) {
	cfg := Default()

	if cfg.Host != DefaultHost {
		t.Errorf("expected host %q, got %q", DefaultHost, cfg.Host)
	}
	if cfg.Interval != DefaultInterval {
		t.Errorf("expected interval %v, got %v", DefaultInterval, cfg.Interval)
	}
	if cfg.PortRangeStart != 1 {
		t.Errorf("expected port_range_start 1, got %d", cfg.PortRangeStart)
	}
	if cfg.PortRangeEnd != 65535 {
		t.Errorf("expected port_range_end 65535, got %d", cfg.PortRangeEnd)
	}
}

func TestLoad_ReadsFile(t *testing.T) {
	input := &Config{
		Host:           "127.0.0.1",
		PortRangeStart: 1024,
		PortRangeEnd:   2048,
		Interval:       10 * time.Second,
		StatePath:      "/tmp/test_state.json",
		AlertOnStart:   true,
	}

	f, err := os.CreateTemp("", "portwatch_cfg_*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())

	if err := json.NewEncoder(f).Encode(input); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	f.Close()

	cfg, err := Load(f.Name())
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Host != input.Host {
		t.Errorf("host: expected %q, got %q", input.Host, cfg.Host)
	}
	if cfg.PortRangeStart != input.PortRangeStart {
		t.Errorf("port_range_start: expected %d, got %d", input.PortRangeStart, cfg.PortRangeStart)
	}
	if cfg.AlertOnStart != input.AlertOnStart {
		t.Errorf("alert_on_start: expected %v, got %v", input.AlertOnStart, cfg.AlertOnStart)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/portwatch.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := Default()
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidate_InvalidRange(t *testing.T) {
	cfg := Default()
	cfg.PortRangeStart = 9000
	cfg.PortRangeEnd = 1000

	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for inverted range, got nil")
	}
}

func TestValidate_ZeroInterval(t *testing.T) {
	cfg := Default()
	cfg.Interval = 0

	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for zero interval, got nil")
	}
}

func TestValidate_OutOfBoundsPort(t *testing.T) {
	cfg := Default()
	cfg.PortRangeEnd = 70000

	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for port > 65535, got nil")
	}
}
