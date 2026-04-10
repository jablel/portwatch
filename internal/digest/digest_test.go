package digest_test

import (
	"testing"

	"portwatch/internal/digest"
	"portwatch/internal/scanner"
)

func makePorts(specs ...string) []scanner.Port {
	ports := make([]scanner.Port, 0, len(specs))
	for _, s := range specs {
		var proto string
		var num uint16
		if _, err := fmt.Sscanf(s, "%5s/%d", &proto, &num); err == nil {
			ports = append(ports, scanner.Port{Protocol: proto, Number: num})
		}
	}
	return ports
}

import "fmt"

func TestCompute_EmptySlice(t *testing.T) {
	d := digest.Compute(nil)
	if d != digest.Empty {
		t.Errorf("expected Empty digest, got %s", d)
	}
}

func TestCompute_Deterministic(t *testing.T) {
	ports := []scanner.Port{
		{Protocol: "tcp", Number: 80},
		{Protocol: "tcp", Number: 443},
	}
	d1 := digest.Compute(ports)
	d2 := digest.Compute(ports)
	if d1 != d2 {
		t.Errorf("digest not deterministic: %s != %s", d1, d2)
	}
}

func TestCompute_OrderIndependent(t *testing.T) {
	a := []scanner.Port{
		{Protocol: "tcp", Number: 80},
		{Protocol: "tcp", Number: 443},
	}
	b := []scanner.Port{
		{Protocol: "tcp", Number: 443},
		{Protocol: "tcp", Number: 80},
	}
	if digest.Compute(a) != digest.Compute(b) {
		t.Error("digest should be order-independent")
	}
}

func TestCompute_DifferentPortsProduceDifferentDigests(t *testing.T) {
	a := []scanner.Port{{Protocol: "tcp", Number: 80}}
	b := []scanner.Port{{Protocol: "tcp", Number: 8080}}
	if digest.Compute(a) == digest.Compute(b) {
		t.Error("different ports should produce different digests")
	}
}

func TestEqual(t *testing.T) {
	ports := []scanner.Port{{Protocol: "udp", Number: 53}}
	d := digest.Compute(ports)
	if !digest.Equal(d, d) {
		t.Error("Equal should return true for same digest")
	}
	if digest.Equal(d, digest.Empty) {
		t.Error("Equal should return false for different digests")
	}
}

func TestDigest_String(t *testing.T) {
	d := digest.Empty
	if d.String() == "" {
		t.Error("String() should not be empty")
	}
}
