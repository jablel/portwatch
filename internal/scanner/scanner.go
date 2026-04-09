package scanner

import (
	"fmt"
	"net"
	"time"
)

// Port represents an open port with its protocol and state.
type Port struct {
	Protocol string
	Number   int
	Address  string
}

// String returns a human-readable representation of the port.
func (p Port) String() string {
	return fmt.Sprintf("%s:%d (%s)", p.Address, p.Number, p.Protocol)
}

// Scanner probes a host for open TCP/UDP ports within a given range.
type Scanner struct {
	Host    string
	MinPort int
	MaxPort int
	Timeout time.Duration
}

// New creates a Scanner with sensible defaults.
func New(host string, minPort, maxPort int) *Scanner {
	return &Scanner{
		Host:    host,
		MinPort: minPort,
		MaxPort: maxPort,
		Timeout: 500 * time.Millisecond,
	}
}

// Scan iterates over the port range and returns all open TCP ports.
func (s *Scanner) Scan() ([]Port, error) {
	var open []Port

	for port := s.MinPort; port <= s.MaxPort; port++ {
		addr := fmt.Sprintf("%s:%d", s.Host, port)
		conn, err := net.DialTimeout("tcp", addr, s.Timeout)
		if err != nil {
			continue
		}
		conn.Close()
		open = append(open, Port{
			Protocol: "tcp",
			Number:   port,
			Address:  s.Host,
		})
	}

	return open, nil
}
