package portpin

import (
	"fmt"
	"testing"

	"portwatch/internal/scanner"
)

func BenchmarkCheck_100Ports(b *testing.B) {
	p := New()
	ports := make([]scanner.Port, 100)
	for i := 0; i < 100; i++ {
		port := scanner.Port{Number: i + 1, Protocol: "tcp"}
		ports[i] = port
		p.Pin(port)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p.Check(ports)
	}
}

func BenchmarkPin_Sequential(b *testing.B) {
	p := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Pin(scanner.Port{Number: i % 65535, Protocol: fmt.Sprintf("proto%d", i%2)})
	}
}
