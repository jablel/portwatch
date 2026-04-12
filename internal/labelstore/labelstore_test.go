package labelstore_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/labelstore"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(proto string, number int) scanner.Port {
	return scanner.Port{Proto: proto, Number: number}
}

func TestSet_AndGet(t *testing.T) {
	s := labelstore.New("")
	p := makePort("tcp", 8080)
	s.Set(p, "dev-proxy")
	got, ok := s.Get(p)
	if !ok {
		t.Fatal("expected label to be present")
	}
	if got != "dev-proxy" {
		t.Fatalf("got %q, want %q", got, "dev-proxy")
	}
}

func TestGet_MissingReturnsNotFound(t *testing.T) {
	s := labelstore.New("")
	_, ok := s.Get(makePort("tcp", 9999))
	if ok {
		t.Fatal("expected label to be absent")
	}
}

func TestDelete_RemovesLabel(t *testing.T) {
	s := labelstore.New("")
	p := makePort("udp", 53)
	s.Set(p, "dns")
	s.Delete(p)
	_, ok := s.Get(p)
	if ok {
		t.Fatal("expected label to be deleted")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "labels.json")

	s1 := labelstore.New(path)
	s1.Set(makePort("tcp", 443), "https")
	s1.Set(makePort("tcp", 22), "ssh")
	if err := s1.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	s2 := labelstore.New(path)
	if err := s2.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	for _, tc := range []struct {
		proto  string
		port   int
		label  string
	}{
		{"tcp", 443, "https"},
		{"tcp", 22, "ssh"},
	} {
		got, ok := s2.Get(makePort(tc.proto, tc.port))
		if !ok || got != tc.label {
			t.Errorf("port %d: got %q ok=%v, want %q", tc.port, got, ok, tc.label)
		}
	}
}

func TestLoad_MissingFile_IsNoop(t *testing.T) {
	dir := t.TempDir()
	s := labelstore.New(filepath.Join(dir, "no-such-file.json"))
	if err := s.Load(); err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
}

func TestLoad_CorruptFile_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "labels.json")
	_ = os.WriteFile(path, []byte("not-json{"), 0o644)
	s := labelstore.New(path)
	if err := s.Load(); err == nil {
		t.Fatal("expected error for corrupt file")
	}
}
