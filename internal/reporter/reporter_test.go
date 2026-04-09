package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/reporter"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

func makeDiff(added, removed []scanner.Port) state.Diff {
	return state.Diff{Added: added, Removed: removed}
}

func TestReport_TextFormat_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	if err := r.Report(makeDiff(nil, nil)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No changes detected") {
		t.Errorf("expected no-changes message, got: %s", buf.String())
	}
}

func TestReport_TextFormat_Added(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	added := []scanner.Port{{Number: 8080, Protocol: "tcp"}}
	if err := r.Report(makeDiff(added, nil)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "ADDED") {
		t.Errorf("expected ADDED in output, got: %s", buf.String())
	}
}

func TestReport_TextFormat_Removed(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	removed := []scanner.Port{{Number: 22, Protocol: "tcp"}}
	if err := r.Report(makeDiff(nil, removed)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "REMOVED") {
		t.Errorf("expected REMOVED in output, got: %s", buf.String())
	}
}

func TestReport_CSVFormat(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatCSV)
	added := []scanner.Port{{Number: 443, Protocol: "tcp"}}
	removed := []scanner.Port{{Number: 80, Protocol: "tcp"}}
	if err := r.Report(makeDiff(added, removed)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 CSV lines, got %d: %s", len(lines), buf.String())
	}
	if !strings.Contains(lines[0], "ADDED") {
		t.Errorf("first line should contain ADDED: %s", lines[0])
	}
	if !strings.Contains(lines[1], "REMOVED") {
		t.Errorf("second line should contain REMOVED: %s", lines[1])
	}
}

func TestNew_DefaultsToStdoutAndText(t *testing.T) {
	r := reporter.New(nil, "")
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}
