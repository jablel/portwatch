package aggregator_test

import (
	"fmt"
	"testing"

	"github.com/user/portwatch/internal/aggregator"
	"github.com/user/portwatch/internal/scanner"
)

func BenchmarkMerge_TenSources(b *testing.B) {
	a := aggregator.New()
	for i := 0; i < 10; i++ {
		ports := make([]scanner.Port, 100)
		for j := range ports {
			ports[j] = scanner.Port{Number: uint16(i*100 + j + 1), Protocol: "tcp"}
		}
		a.Update(fmt.Sprintf("src%d", i), ports)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = a.Merge()
	}
}

func BenchmarkUpdate_SingleSource(b *testing.B) {
	a := aggregator.New()
	ports := makePorts(80, 443, 8080, 8443, 9090)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a.Update("src", ports)
	}
}
