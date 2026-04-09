package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/state"
)

func TestNotify_PortAdded(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	diff := state.Diff{
		Added: []state.Port{{Number: 8080, Protocol: "tcp"}},
	}
	n.Notify(diff)

	out := buf.String()
	if !strings.Contains(out, "ALERT") {
		t.Errorf("expected ALERT level in output, got: %s", out)
	}
	if !strings.Contains(out, "port opened") {
		t.Errorf("expected 'port opened' in output, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port number in output, got: %s", out)
	}
}

func TestNotify_PortRemoved(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	diff := state.Diff{
		Removed: []state.Port{{Number: 22, Protocol: "tcp"}},
	}
	n.Notify(diff)

	out := buf.String()
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN level in output, got: %s", out)
	}
	if !strings.Contains(out, "port closed") {
		t.Errorf("expected 'port closed' in output, got: %s", out)
	}
	if !strings.Contains(out, "22") {
		t.Errorf("expected port number in output, got: %s", out)
	}
}

func TestNotify_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	n.Notify(state.Diff{})

	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff, got: %s", buf.String())
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	// Ensure New(nil) does not panic
	n := alert.New(nil)
	if n == nil {
		t.Fatal("expected non-nil Notifier")
	}
}
