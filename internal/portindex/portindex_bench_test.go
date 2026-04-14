package portindex_test

import (
	"fmt"
	"testing"

	"github.com/user/portwatch/internal/portindex"
	"github.com/user/portwatch/internal/scanner"
)

func largeBatch(n int) []scanner.Port {
	ports := make([]scanner.Port, n)
	for i := 0; i < n; i++ {
		ports[i] = scanner.Port{
			Number:   i + 1,
			Protocol: "tcp",
			Tags:     []string{fmt.Sprintf("tag%d", i%10)},
		}
	}
	return ports
}

func BenchmarkBuild_1000Ports(b *testing.B) {
	batch := largeBatch(1000)
	idx := portindex.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		idx.Build(batch)
	}
}

func BenchmarkByNumber(b *testing.B) {
	batch := largeBatch(1000)
	idx := portindex.New()
	idx.Build(batch)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = idx.ByNumber(500)
	}
}

func BenchmarkByTag(b *testing.B) {
	batch := largeBatch(1000)
	idx := portindex.New()
	idx.Build(batch)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = idx.ByTag("tag3")
	}
}
