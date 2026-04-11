package tagger_test

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/tagger"
)

func port(n uint16, proto string) scanner.Port {
	return scanner.Port{Number: n, Protocol: proto}
}

func TestTag_WellKnownHTTP(t *testing.T) {
	tg := tagger.New()
	if got := tg.Tag(port(80, "tcp")); got != "http" {
		t.Fatalf("expected http, got %q", got)
	}
}

func TestTag_WellKnownSSH(t *testing.T) {
	tg := tagger.New()
	if got := tg.Tag(port(22, "tcp")); got != "ssh" {
		t.Fatalf("expected ssh, got %q", got)
	}
}

func TestTag_UnknownPort(t *testing.T) {
	tg := tagger.New()
	if got := tg.Tag(port(9999, "tcp")); got != "unknown:9999" {
		t.Fatalf("unexpected label %q", got)
	}
}

func TestTag_CustomOverridesWellKnown(t *testing.T) {
	tg := tagger.New()
	tg.Define(80, "my-app")
	if got := tg.Tag(port(80, "tcp")); got != "my-app" {
		t.Fatalf("expected my-app, got %q", got)
	}
}

func TestTag_CustomUnknownPort(t *testing.T) {
	tg := tagger.New()
	tg.Define(12345, "custom-svc")
	if got := tg.Tag(port(12345, "tcp")); got != "custom-svc" {
		t.Fatalf("expected custom-svc, got %q", got)
	}
}

func TestTagAll_ReturnsAllLabels(t *testing.T) {
	tg := tagger.New()
	ports := []scanner.Port{
		port(22, "tcp"),
		port(80, "tcp"),
		port(9999, "tcp"),
	}
	labels := tg.TagAll(ports)
	if len(labels) != 3 {
		t.Fatalf("expected 3 labels, got %d", len(labels))
	}
	if labels[port(22, "tcp")] != "ssh" {
		t.Errorf("expected ssh for port 22")
	}
	if labels[port(80, "tcp")] != "http" {
		t.Errorf("expected http for port 80")
	}
	if labels[port(9999, "tcp")] != "unknown:9999" {
		t.Errorf("expected unknown:9999 for port 9999")
	}
}

func TestTagAll_EmptySlice(t *testing.T) {
	tg := tagger.New()
	if got := tg.TagAll(nil); len(got) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(got))
	}
}
