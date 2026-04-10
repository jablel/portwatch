package notifier_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/notifier"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

func makeDiff(added, removed []string) state.Diff {
	toPort := func(ss []string) []scanner.Port {
		var out []scanner.Port
		for _, s := range ss {
			out = append(out, scanner.Port{Proto: "tcp", Number: 0, Raw: s})
		}
		return out
	}
	return state.Diff{Added: toPort(added), Removed: toPort(removed)}
}

func TestParseBackend_Valid(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  notifier.Backend
	}{
		{"stdout", notifier.BackendStdout},
		{"STDOUT", notifier.BackendStdout},
		{"file", notifier.BackendFile},
		{"FILE", notifier.BackendFile},
	} {
		got, err := notifier.ParseBackend(tc.input)
		if err != nil {
			t.Fatalf("ParseBackend(%q) error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseBackend(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestParseBackend_Invalid(t *testing.T) {
	_, err := notifier.ParseBackend("webhook")
	if err == nil {
		t.Fatal("expected error for unknown backend, got nil")
	}
}

func TestNotify_WritesToWriter(t *testing.T) {
	var buf bytes.Buffer
	n, err := notifier.New(notifier.BackendStdout, "")
	if err != nil {
		t.Fatal(err)
	}
	// Redirect internal writer via a helper — test via file backend instead.
	_ = n

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "notify.log")
	fn, err := notifier.New(notifier.BackendFile, path)
	if err != nil {
		t.Fatal(err)
	}

	diff := makeDiff([]string{"tcp:8080"}, []string{"tcp:9090"})
	if err := fn.Notify(diff); err != nil {
		t.Fatalf("Notify error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	out := string(data)
	if !strings.Contains(out, "PORT OPENED") {
		t.Errorf("expected PORT OPENED in output, got: %s", out)
	}
	if !strings.Contains(out, "PORT CLOSED") {
		t.Errorf("expected PORT CLOSED in output, got: %s", out)
	}
	_ = buf
}

func TestNotify_NoChanges_WritesNothing(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "empty.log")
	n, err := notifier.New(notifier.BackendFile, path)
	if err != nil {
		t.Fatal(err)
	}
	if err := n.Notify(state.Diff{}); err != nil {
		t.Fatal(err)
	}
	info, _ := os.Stat(path)
	if info != nil && info.Size() > 0 {
		t.Error("expected no output for empty diff")
	}
}

func TestNew_FileBackend_MissingPath(t *testing.T) {
	_, err := notifier.New(notifier.BackendFile, "")
	if err == nil {
		t.Fatal("expected error for empty file path")
	}
}
