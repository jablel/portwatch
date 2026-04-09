package main

import (
	"os"
	"os/exec"
	"testing"
	"time"
)

// TestMain_BuildAndRun compiles the binary and verifies it starts and
// terminates cleanly on SIGINT within a short window.
func TestMain_BuildAndRun(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping binary integration test in short mode")
	}

	tmpDir := t.TempDir()
	binPath := tmpDir + "/portwatch"

	build := exec.Command("go", "build", "-o", binPath, ".")
	build.Dir = "."
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}

	cmd := exec.Command(binPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start binary: %v", err)
	}

	time.Sleep(150 * time.Millisecond)

	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		t.Fatalf("failed to send interrupt: %v", err)
	}

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case err := <-done:
		if err != nil {
			t.Logf("process exited with: %v (may be normal on some platforms)", err)
		}
	case <-time.After(3 * time.Second):
		_ = cmd.Process.Kill()
		t.Error("binary did not exit within timeout after SIGINT")
	}
}
