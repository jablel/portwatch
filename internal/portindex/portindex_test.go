package portindex_test

import (
	"testing"

	"github.com/user/portwatch/internal/portindex"
	"github.com/user/portwatch/internal/scanner"
)

func makePorts() []scanner.Port {
	return []scanner.Port{
		{Number: 80, Protocol: "tcp", Tags: []string{"http", "web"}},
		{Number: 443, Protocol: "tcp", Tags: []string{"https", "web"}},
		{Number: 53, Protocol: "udp", Tags: []string{"dns"}},
		{Number: 80, Protocol: "udp", Tags: []string{"http"}},
	}
}

func TestBuild_PopulatesIndex(t *testing.T) {
	idx := portindex.New()
	idx.Build(makePorts())

	if got := idx.Size(); got != 4 {
		t.Fatalf("expected size 4, got %d", got)
	}
}

func TestByNumber_ReturnsMatchingPorts(t *testing.T) {
	idx := portindex.New()
	idx.Build(makePorts())

	ports := idx.ByNumber(80)
	if len(ports) != 2 {
		t.Fatalf("expected 2 ports for 80, got %d", len(ports))
	}
}

func TestByNumber_UnknownReturnsNil(t *testing.T) {
	idx := portindex.New()
	idx.Build(makePorts())

	if got := idx.ByNumber(9999); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestByProtocol_ReturnsTCPPorts(t *testing.T) {
	idx := portindex.New()
	idx.Build(makePorts())

	ports := idx.ByProtocol("tcp")
	if len(ports) != 2 {
		t.Fatalf("expected 2 tcp ports, got %d", len(ports))
	}
}

func TestByTag_ReturnsTaggedPorts(t *testing.T) {
	idx := portindex.New()
	idx.Build(makePorts())

	ports := idx.ByTag("web")
	if len(ports) != 2 {
		t.Fatalf("expected 2 web ports, got %d", len(ports))
	}
}

func TestByTag_UnknownTagReturnsNil(t *testing.T) {
	idx := portindex.New()
	idx.Build(makePorts())

	if got := idx.ByTag("nonexistent"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestBuild_ReplacesExistingIndex(t *testing.T) {
	idx := portindex.New()
	idx.Build(makePorts())
	idx.Build([]scanner.Port{{Number: 22, Protocol: "tcp"}})

	if got := idx.Size(); got != 1 {
		t.Fatalf("expected size 1 after rebuild, got %d", got)
	}
	if got := idx.ByNumber(80); got != nil {
		t.Fatal("expected old entries to be cleared")
	}
}

func TestKey_Format(t *testing.T) {
	p := scanner.Port{Number: 443, Protocol: "tcp"}
	if got := portindex.Key(p); got != "tcp:443" {
		t.Fatalf("expected tcp:443, got %s", got)
	}
}
