package portmatcher_test

import (
	"testing"

	"github.com/user/portwatch/internal/portmatcher"
	"github.com/user/portwatch/internal/scanner"
)

func port(number uint16, proto string) scanner.Port {
	return scanner.Port{Number: number, Protocol: proto}
}

func TestAdd_InvalidPattern(t *testing.T) {
	m := portmatcher.New()
	for _, bad := range []string{"abc", "80:ftp", "9000-8000", "65536"} {
		if err := m.Add(bad); err == nil {
			t.Errorf("expected error for pattern %q, got nil", bad)
		}
	}
}

func TestMatch_ExactPort(t *testing.T) {
	m := portmatcher.New()
	if err := m.Add("80"); err != nil {
		t.Fatal(err)
	}
	if !m.Match(port(80, "tcp")) {
		t.Error("expected port 80/tcp to match")
	}
	if !m.Match(port(80, "udp")) {
		t.Error("expected port 80/udp to match (no protocol filter)")
	}
	if m.Match(port(81, "tcp")) {
		t.Error("expected port 81/tcp not to match")
	}
}

func TestMatch_Range(t *testing.T) {
	m := portmatcher.New()
	if err := m.Add("8000-8080"); err != nil {
		t.Fatal(err)
	}
	for _, n := range []uint16{8000, 8040, 8080} {
		if !m.Match(port(n, "tcp")) {
			t.Errorf("expected port %d to match", n)
		}
	}
	for _, n := range []uint16{7999, 8081} {
		if m.Match(port(n, "tcp")) {
			t.Errorf("expected port %d not to match", n)
		}
	}
}

func TestMatch_ProtocolFilter(t *testing.T) {
	m := portmatcher.New()
	if err := m.Add("443:tcp"); err != nil {
		t.Fatal(err)
	}
	if !m.Match(port(443, "tcp")) {
		t.Error("expected 443/tcp to match")
	}
	if m.Match(port(443, "udp")) {
		t.Error("expected 443/udp not to match")
	}
}

func TestMatch_RangeWithProtocol(t *testing.T) {
	m := portmatcher.New()
	if err := m.Add("5000-5010:udp"); err != nil {
		t.Fatal(err)
	}
	if !m.Match(port(5005, "udp")) {
		t.Error("expected 5005/udp to match")
	}
	if m.Match(port(5005, "tcp")) {
		t.Error("expected 5005/tcp not to match")
	}
}

func TestMatchAny_FiltersSlice(t *testing.T) {
	m := portmatcher.New()
	_ = m.Add("22")
	_ = m.Add("80")
	ports := []scanner.Port{
		port(22, "tcp"),
		port(80, "tcp"),
		port(443, "tcp"),
		port(8080, "tcp"),
	}
	got := m.MatchAny(ports)
	if len(got) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(got))
	}
}

func TestMatchAny_EmptyRulesReturnsNone(t *testing.T) {
	m := portmatcher.New()
	ports := []scanner.Port{port(80, "tcp"), port(443, "tcp")}
	got := m.MatchAny(ports)
	if len(got) != 0 {
		t.Errorf("expected no matches with empty rules, got %d", len(got))
	}
}
