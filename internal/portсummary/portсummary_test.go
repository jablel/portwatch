package portсummary

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

func makePorts(specs [][2]string) []scanner.Port {
	out := make([]scanner.Port, len(specs))
	for i, s := range specs {
		out[i] = scanner.Port{Proto: s[0], Number: 0}
		_ = s[1]
	}
	return out
}

func TestBuild_TotalsAreCorrect(t *testing.T) {
	ports := []scanner.Port{
		{Proto: "tcp", Number: 80},
		{Proto: "tcp", Number: 443},
		{Proto: "udp", Number: 53},
	}
	diff := state.Diff{
		Added:   []scanner.Port{{Proto: "tcp", Number: 8080}},
		Removed: []scanner.Port{},
	}
	b := New()
	s := b.Build(ports, diff)
	if s.Total != 3 {
		t.Errorf("expected Total=3, got %d", s.Total)
	}
	if s.ByProto["tcp"] != 2 {
		t.Errorf("expected tcp=2, got %d", s.ByProto["tcp"])
	}
	if s.ByProto["udp"] != 1 {
		t.Errorf("expected udp=1, got %d", s.ByProto["udp"])
	}
	if s.Added != 1 {
		t.Errorf("expected Added=1, got %d", s.Added)
	}
	if s.Removed != 0 {
		t.Errorf("expected Removed=0, got %d", s.Removed)
	}
}

func TestBuild_EmptyPorts(t *testing.T) {
	b := New()
	s := b.Build(nil, state.Diff{})
	if s.Total != 0 {
		t.Errorf("expected Total=0, got %d", s.Total)
	}
}

func TestWrite_ContainsExpectedLines(t *testing.T) {
	ports := []scanner.Port{
		{Proto: "tcp", Number: 22},
	}
	diff := state.Diff{
		Added:   []scanner.Port{{Proto: "tcp", Number: 22}},
		Removed: []scanner.Port{{Proto: "tcp", Number: 80}},
	}
	b := New()
	s := b.Build(ports, diff)
	var buf bytes.Buffer
	b.Write(&buf, s)
	out := buf.String()
	for _, want := range []string{"Total ports:", "tcp:", "Added:", "Removed:"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n%s", want, out)
		}
	}
}
