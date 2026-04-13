package portclassifier_test

import (
	"testing"

	"github.com/user/portwatch/internal/portclassifier"
	"github.com/user/portwatch/internal/scanner"
)

func BenchmarkClassify(b *testing.B) {
	c := portclassifier.New(nil)
	p := scanner.Port{Number: 443, Protocol: "tcp"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.Classify(p)
	}
}

func BenchmarkClassifyAll_100Ports(b *testing.B) {
	c := portclassifier.New(nil)
	ports := make([]scanner.Port, 100)
	for i := range ports {
		ports[i] = scanner.Port{Number: uint16(i * 655), Protocol: "tcp"}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.ClassifyAll(ports)
	}
}
