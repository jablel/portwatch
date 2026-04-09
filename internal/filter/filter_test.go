package filter_test

import (
	"testing"

	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/scanner"
)

func makePorts(specs [][2]interface{}) []scanner.Port {
	var ports []scanner.Port
	for _, s := range specs {
		ports = append(ports, scanner.Port{
			Number:   s[0].(uint16),
			Protocol: s[1].(string),
		})
	}
	return ports
}

func TestApply_NoRules_ReturnsAll(t *testing.T) {
	f := filter.New(nil)
	ports := makePorts([][2]interface{}{{uint16(80), "tcp"}, {uint16(443), "tcp"}})
	got := f.Apply(ports)
	if len(got) != len(ports) {
		t.Fatalf("expected %d ports, got %d", len(ports), len(got))
	}
}

func TestApply_MatchesRange(t *testing.T) {
	rules := []filter.Rule{{MinPort: 80, MaxPort: 90, Protocols: nil}}
	f := filter.New(rules)
	ports := makePorts([][2]interface{}{{uint16(80), "tcp"}, {uint16(443), "tcp"}, {uint16(85), "udp"}})
	got := f.Apply(ports)
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(got))
	}
}

func TestApply_MatchesProtocol(t *testing.T) {
	rules := []filter.Rule{{MinPort: 1, MaxPort: 1024, Protocols: []string{"tcp"}}}
	f := filter.New(rules)
	ports := makePorts([][2]interface{}{{uint16(80), "tcp"}, {uint16(53), "udp"}})
	got := f.Apply(ports)
	if len(got) != 1 || got[0].Protocol != "tcp" {
		t.Fatalf("expected only tcp port, got %+v", got)
	}
}

func TestExclude_RemovesMatched(t *testing.T) {
	rules := []filter.Rule{{MinPort: 80, MaxPort: 80, Protocols: nil}}
	f := filter.New(rules)
	ports := makePorts([][2]interface{}{{uint16(80), "tcp"}, {uint16(443), "tcp"}})
	got := f.Exclude(ports)
	if len(got) != 1 || got[0].Number != 443 {
		t.Fatalf("expected port 443 only, got %+v", got)
	}
}

func TestExclude_NoRules_ReturnsAll(t *testing.T) {
	f := filter.New(nil)
	ports := makePorts([][2]interface{}{{uint16(22), "tcp"}})
	got := f.Exclude(ports)
	if len(got) != 1 {
		t.Fatalf("expected 1 port, got %d", len(got))
	}
}
