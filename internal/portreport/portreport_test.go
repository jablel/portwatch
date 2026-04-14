package portreport_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/portreport"
	"github.com/user/portwatch/internal/scanner"
)

// --- stubs ---

type stubClassifier struct{ label string }

func (s *stubClassifier) Classify(_ scanner.Port) string { return s.label }

type stubTrencher struct{ label string }

func (s *stubTrencher) Trend(_ scanner.Port) string { return s.label }

type stubLifecycler struct{ label string }

func (s *stubLifecycler) State(_ scanner.Port) string { return s.label }

func makePort(number int, proto string) scanner.Port {
	return scanner.Port{Number: number, Proto: proto}
}

// --- tests ---

func TestBuild_PopulatesEntries(t *testing.T) {
	r := portreport.New(
		&stubClassifier{"system"},
		&stubTrencher{"rising"},
		&stubLifecycler{"active"},
	)
	ports := []scanner.Port{makePort(80, "tcp"), makePort(443, "tcp")}
	entries := r.Build(ports)

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if e.Class != "system" {
			t.Errorf("expected class=system, got %s", e.Class)
		}
		if e.Trend != "rising" {
			t.Errorf("expected trend=rising, got %s", e.Trend)
		}
		if e.State != "active" {
			t.Errorf("expected state=active, got %s", e.State)
		}
	}
}

func TestBuild_EmptyPorts(t *testing.T) {
	r := portreport.New(&stubClassifier{}, &stubTrencher{}, &stubLifecycler{})
	entries := r.Build(nil)
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestWrite_ContainsHeader(t *testing.T) {
	r := portreport.New(
		&stubClassifier{"dynamic"},
		&stubTrencher{"stable"},
		&stubLifecycler{"new"},
	)
	entries := r.Build([]scanner.Port{makePort(9000, "udp")})

	var buf bytes.Buffer
	if err := r.Write(&buf, entries); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "PORT") || !strings.Contains(out, "TREND") {
		t.Errorf("output missing header columns: %q", out)
	}
	if !strings.Contains(out, "9000") {
		t.Errorf("output missing port number: %q", out)
	}
}

func TestWrite_EmptyEntries_OnlyHeader(t *testing.T) {
	r := portreport.New(&stubClassifier{}, &stubTrencher{}, &stubLifecycler{})
	var buf bytes.Buffer
	if err := r.Write(&buf, nil); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected only header line, got %d lines", len(lines))
	}
}
