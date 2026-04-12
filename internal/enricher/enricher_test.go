package enricher_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/enricher"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/tagger"
)

func newEnricher() *enricher.Enricher {
	t := tagger.New()
	return enricher.New(t)
}

func makePort(number uint16, proto string) scanner.Port {
	return scanner.Port{Number: number, Protocol: proto}
}

func TestEnrich_WellKnownPort(t *testing.T) {
	e := newEnricher()
	ports := []scanner.Port{makePort(80, "tcp")}
	result := e.Enrich(ports)
	if len(result) != 1 {
		t.Fatalf("expected 1 enriched port, got %d", len(result))
	}
	if result[0].ServiceName == "" {
		t.Error("expected non-empty service name for port 80")
	}
	if !strings.Contains(strings.ToLower(result[0].ServiceName), "http") {
		t.Errorf("expected http in service name, got %q", result[0].ServiceName)
	}
}

func TestEnrich_UnknownPort(t *testing.T) {
	e := newEnricher()
	ports := []scanner.Port{makePort(19999, "tcp")}
	result := e.Enrich(ports)
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	// Unknown ports should still be returned, just without a service name.
	if result[0].Port.Number != 19999 {
		t.Errorf("expected port 19999, got %d", result[0].Port.Number)
	}
}

func TestEnrich_EmptySlice(t *testing.T) {
	e := newEnricher()
	result := e.Enrich([]scanner.Port{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestEnrichOne_SSH(t *testing.T) {
	e := newEnricher()
	p := makePort(22, "tcp")
	result := e.EnrichOne(p)
	if result.Port.Number != 22 {
		t.Errorf("expected port 22, got %d", result.Port.Number)
	}
	if !strings.Contains(strings.ToLower(result.ServiceName), "ssh") {
		t.Errorf("expected ssh in service name, got %q", result.ServiceName)
	}
}

func TestEnrichedPort_String_WithServiceName(t *testing.T) {
	ep := enricher.EnrichedPort{
		Port:        makePort(443, "tcp"),
		ServiceName: "https",
	}
	s := ep.String()
	if !strings.Contains(s, "https") {
		t.Errorf("expected 'https' in string output, got %q", s)
	}
}

func TestEnrichedPort_String_NoServiceName(t *testing.T) {
	ep := enricher.EnrichedPort{
		Port: makePort(9999, "udp"),
	}
	s := ep.String()
	if s == "" {
		t.Error("expected non-empty string for port with no service name")
	}
}
