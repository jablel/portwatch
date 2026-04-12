package profiler_test

import (
	"testing"
	"time"

	"portwatch/internal/profiler"
)

func TestRecord_SingleSample(t *testing.T) {
	p := profiler.New()
	p.Record(profiler.Sample{Key: "tcp:80", Duration: 10 * time.Millisecond, RecordedAt: time.Now()})

	st, ok := p.Get("tcp:80")
	if !ok {
		t.Fatal("expected stats for tcp:80")
	}
	if st.Count != 1 {
		t.Errorf("count: got %d, want 1", st.Count)
	}
	if st.Mean() != 10*time.Millisecond {
		t.Errorf("mean: got %v, want 10ms", st.Mean())
	}
}

func TestRecord_UpdatesMinMax(t *testing.T) {
	p := profiler.New()
	for _, d := range []time.Duration{30 * time.Millisecond, 5 * time.Millisecond, 20 * time.Millisecond} {
		p.Record(profiler.Sample{Key: "udp:53", Duration: d, RecordedAt: time.Now()})
	}

	st, _ := p.Get("udp:53")
	if st.Min != 5*time.Millisecond {
		t.Errorf("min: got %v, want 5ms", st.Min)
	}
	if st.Max != 30*time.Millisecond {
		t.Errorf("max: got %v, want 30ms", st.Max)
	}
	if st.Count != 3 {
		t.Errorf("count: got %d, want 3", st.Count)
	}
}

func TestGet_MissingKey(t *testing.T) {
	p := profiler.New()
	_, ok := p.Get("tcp:9999")
	if ok {
		t.Error("expected no stats for unknown key")
	}
}

func TestAll_ReturnsAllKeys(t *testing.T) {
	p := profiler.New()
	p.Record(profiler.Sample{Key: "tcp:22", Duration: 1 * time.Millisecond, RecordedAt: time.Now()})
	p.Record(profiler.Sample{Key: "tcp:443", Duration: 2 * time.Millisecond, RecordedAt: time.Now()})

	all := p.All()
	if len(all) != 2 {
		t.Errorf("len: got %d, want 2", len(all))
	}
}

func TestReset_ClearsData(t *testing.T) {
	p := profiler.New()
	p.Record(profiler.Sample{Key: "tcp:80", Duration: 5 * time.Millisecond, RecordedAt: time.Now()})
	p.Reset()

	if all := p.All(); len(all) != 0 {
		t.Errorf("expected empty after reset, got %d entries", len(all))
	}
}

func TestMean_ZeroCount(t *testing.T) {
	st := profiler.Stats{}
	if st.Mean() != 0 {
		t.Errorf("mean of empty stats should be 0, got %v", st.Mean())
	}
}
