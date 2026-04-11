package rollup_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/rollup"
)

func BenchmarkAdd_ZeroWindow(b *testing.B) {
	done := make(chan rollup.Diff, b.N)
	r := rollup.New(0, func(d rollup.Diff) { done <- d })
	d := makeDiff([]uint16{8080}, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Add(d)
	}
}

func BenchmarkAdd_WithWindow(b *testing.B) {
	r := rollup.New(10*time.Second, func(_ rollup.Diff) {})
	d := makeDiff([]uint16{8080}, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Add(d)
	}
	r.Flush()
}

func BenchmarkFlush_NoPending(b *testing.B) {
	r := rollup.New(10*time.Second, func(_ rollup.Diff) {})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Flush()
	}
}
