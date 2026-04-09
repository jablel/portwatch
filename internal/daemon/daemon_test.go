package daemon_test

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"portwatch/internal/alert"
	"portwatch/internal/config"
	"portwatch/internal/daemon"
)

func testConfig(t *testing.T) *config.Config {
	t.Helper()
	cfg := config.Default()
	cfg.Interval = 50 * time.Millisecond
	cfg.PortRangeStart = 19900
	cfg.PortRangeEnd = 19910
	cfg.StateFile = t.TempDir() + "/state.json"
	return cfg
}

func TestDaemon_RunAndStop(t *testing.T) {
	cfg := testConfig(t)
	var buf bytes.Buffer
	alerter := alert.New(&buf)

	d := daemon.New(cfg, alerter)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	err := d.Run(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}

func TestDaemon_MissingStateFile(t *testing.T) {
	cfg := testConfig(t)
	cfg.StateFile = t.TempDir() + "/nonexistent/state.json"
	var buf bytes.Buffer
	alerter := alert.New(&buf)

	d := daemon.New(cfg, alerter)

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	// Should not panic even if state directory is missing
	err := d.Run(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}

func TestDaemon_WritesStateFile(t *testing.T) {
	cfg := testConfig(t)
	var buf bytes.Buffer
	alerter := alert.New(&buf)

	d := daemon.New(cfg, alerter)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_ = d.Run(ctx)

	if _, err := os.Stat(cfg.StateFile); os.IsNotExist(err) {
		t.Error("expected state file to be written, but it does not exist")
	}
}
