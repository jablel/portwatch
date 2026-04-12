package labelstore_test

import (
	"path/filepath"
	"sync"
	"testing"

	"github.com/user/portwatch/internal/labelstore"
)

func TestLabelStore_ConcurrentSetAndGet(t *testing.T) {
	s := labelstore.New("")
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			p := makePort("tcp", 1024+n)
			s.Set(p, "worker")
			_, _ = s.Get(p)
		}(i)
	}
	wg.Wait()
}

func TestLabelStore_OverwritePreservesLatest(t *testing.T) {
	s := labelstore.New("")
	p := makePort("tcp", 8080)
	s.Set(p, "first")
	s.Set(p, "second")
	got, ok := s.Get(p)
	if !ok || got != "second" {
		t.Fatalf("got %q ok=%v, want %q", got, ok, "second")
	}
}

func TestLabelStore_PersistAndReload(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "labels.json")

	s := labelstore.New(path)
	ports := []struct {
		proto string
		port  int
		label string
	}{
		{"tcp", 80, "http"},
		{"tcp", 443, "https"},
		{"udp", 123, "ntp"},
	}
	for _, tc := range ports {
		s.Set(makePort(tc.proto, tc.port), tc.label)
	}
	if err := s.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	s2 := labelstore.New(path)
	if err := s2.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	for _, tc := range ports {
		got, ok := s2.Get(makePort(tc.proto, tc.port))
		if !ok || got != tc.label {
			t.Errorf("%s:%d: got %q ok=%v, want %q", tc.proto, tc.port, got, ok, tc.label)
		}
	}
}
