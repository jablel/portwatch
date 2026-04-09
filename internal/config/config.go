package config

import (
	"encoding/json"
	"os"
	"time"
)

// Default values
const (
	DefaultInterval  = 30 * time.Second
	DefaultStatePath = "/tmp/portwatch_state.json"
	DefaultHost      = "localhost"
)

// Config holds the runtime configuration for portwatch.
type Config struct {
	// Host to scan (default: localhost)
	Host string `json:"host"`

	// PortRangeStart is the first port in the scan range (inclusive).
	PortRangeStart int `json:"port_range_start"`

	// PortRangeEnd is the last port in the scan range (inclusive).
	PortRangeEnd int `json:"port_range_end"`

	// Interval between scans.
	Interval time.Duration `json:"interval"`

	// StatePath is the file path used to persist port state.
	StatePath string `json:"state_path"`

	// AlertOnStart controls whether an alert is emitted on the first scan.
	AlertOnStart bool `json:"alert_on_start"`
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		Host:           DefaultHost,
		PortRangeStart: 1,
		PortRangeEnd:   65535,
		Interval:       DefaultInterval,
		StatePath:      DefaultStatePath,
		AlertOnStart:   false,
	}
}

// Load reads a JSON config file from path and returns a Config.
// Fields not present in the file retain their zero values; callers
// should start from Default() and merge if desired.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := Default()
	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Validate returns an error if the configuration contains invalid values.
func (c *Config) Validate() error {
	if c.PortRangeStart < 1 || c.PortRangeStart > 65535 {
		return &ValidationError{Field: "port_range_start", Value: c.PortRangeStart, Msg: "must be between 1 and 65535"}
	}
	if c.PortRangeEnd < 1 || c.PortRangeEnd > 65535 {
		return &ValidationError{Field: "port_range_end", Value: c.PortRangeEnd, Msg: "must be between 1 and 65535"}
	}
	if c.PortRangeStart > c.PortRangeEnd {
		return &ValidationError{Field: "port_range_start", Value: c.PortRangeStart, Msg: "must be less than or equal to port_range_end"}
	}
	if c.Interval <= 0 {
		return &ValidationError{Field: "interval", Value: int(c.Interval), Msg: "must be positive"}
	}
	return nil
}

// ValidationError describes a configuration validation failure.
type ValidationError struct {
	Field string
	Value int
	Msg   string
}

func (e *ValidationError) Error() string {
	return "config: field \"" + e.Field + "\": " + e.Msg
}
