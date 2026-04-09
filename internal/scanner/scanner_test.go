package scanner_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/scanner"
)

// startListener opens a TCP listener on a random port and returns the port number and a stop function.
func startListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return port, func() { ln.Close() }
}

func TestScan_DetectsOpenPort(t *testing.T) {
	port, stop := startListener(t)
	defer stop()

	s := scanner.New("127.0.0.1", port, port)
	s.Timeout = 200 * time.Millisecond

	ports, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if len(ports) != 1 {
		t.Fatalf("expected 1 open port, got %d", len(ports))
	}
	if ports[0].Number != port {
		t.Errorf("expected port %d, got %d", port, ports[0].Number)
	}
	if ports[0].Protocol != "tcp" {
		t.Errorf("expected protocol tcp, got %s", ports[0].Protocol)
	}
}

func TestScan_NoOpenPorts(t *testing.T) {
	// Use a port range unlikely to have anything listening in CI.
	s := scanner.New("127.0.0.1", 19800, 19810)
	s.Timeout = 100 * time.Millisecond

	ports, err := s.Scan()
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if len(ports) != 0 {
		t.Errorf("expected 0 open ports, got %d", len(ports))
	}
}

func TestPort_String(t *testing.T) {
	p := scanner.Port{Protocol: "tcp", Number: 8080, Address: "127.0.0.1"}
	want := fmt.Sprintf("127.0.0.1:8080 (tcp)")
	if p.String() != want {
		t.Errorf("Port.String() = %q, want %q", p.String(), want)
	}
}
