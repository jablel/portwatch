package portmap_test

import (
	"testing"

	"github.com/user/portwatch/internal/portmap"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(proto string, number int) scanner.Port {
	return scanner.Port{Proto: proto, Number: number}
}

func TestRegister_LookupService_Found(t *testing.T) {
	m := portmap.New()
	p := makePort("tcp", 80)
	m.Register(p, "http")

	svc, ok := m.LookupService(p)
	if !ok {
		t.Fatal("expected service to be found")
	}
	if svc != "http" {
		t.Fatalf("expected \"http\", got %q", svc)
	}
}

func TestLookupService_NotFound(t *testing.T) {
	m := portmap.New()
	_, ok := m.LookupService(makePort("tcp", 9999))
	if ok {
		t.Fatal("expected not found")
	}
}

func TestLookupPorts_ReturnsMappedPorts(t *testing.T) {
	m := portmap.New()
	m.Register(makePort("tcp", 80), "http")
	m.Register(makePort("tcp", 8080), "http")

	ports := m.LookupPorts("http")
	if len(ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(ports))
	}
}

func TestLookupPorts_UnknownServiceReturnsEmpty(t *testing.T) {
	m := portmap.New()
	ports := m.LookupPorts("unknown")
	if len(ports) != 0 {
		t.Fatalf("expected empty slice, got %d entries", len(ports))
	}
}

func TestRegister_OverwriteMovesService(t *testing.T) {
	m := portmap.New()
	p := makePort("tcp", 443)
	m.Register(p, "https")
	m.Register(p, "tls")

	svc, ok := m.LookupService(p)
	if !ok || svc != "tls" {
		t.Fatalf("expected \"tls\", got %q (ok=%v)", svc, ok)
	}
	if ports := m.LookupPorts("https"); len(ports) != 0 {
		t.Fatalf("old service should have no ports, got %d", len(ports))
	}
}

func TestLen_ReflectsRegistrations(t *testing.T) {
	m := portmap.New()
	if m.Len() != 0 {
		t.Fatal("expected empty map")
	}
	m.Register(makePort("tcp", 22), "ssh")
	m.Register(makePort("udp", 53), "dns")
	if m.Len() != 2 {
		t.Fatalf("expected 2, got %d", m.Len())
	}
}
