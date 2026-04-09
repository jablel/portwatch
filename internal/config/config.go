package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/user/portwatch/internal/filter"
)

// Config holds the full portwatch configuration.
type Config struct {
	ScanInterval int           `yaml:"scan_interval"` // seconds
	PortRange    PortRange     `yaml:"port_range"`
	StateFile    string        `yaml:"state_file"`
	AlertLog     string        `yaml:"alert_log"`
	Ignore       []FilterEntry `yaml:"ignore"`
}

// PortRange defines the inclusive range of ports to scan.
type PortRange struct {
	Min uint16 `yaml:"min"`
	Max uint16 `yaml:"max"`
}

// FilterEntry maps to a filter.Rule in YAML form.
type FilterEntry struct {
	MinPort   uint16   `yaml:"min_port"`
	MaxPort   uint16   `yaml:"max_port"`
	Protocols []string `yaml:"protocols"`
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		ScanInterval: 60,
		PortRange:    PortRange{Min: 1, Max: 65535},
		StateFile:    "/var/lib/portwatch/state.json",
		AlertLog:     "",
	}
}

// Load reads a YAML config file, falling back to defaults for missing fields.
func Load(path string) (*Config, error) {
	cfg := Default()
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return nil, err
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, cfg.Validate()
}

// Validate checks that the configuration values are sensible.
func (c *Config) Validate() error {
	if c.ScanInterval <= 0 {
		return errors.New("scan_interval must be greater than 0")
	}
	if c.PortRange.Min > c.PortRange.Max {
		return errors.New("port_range min must be <= max")
	}
	return nil
}

// IgnoreFilter builds a filter.Filter from the Ignore entries.
func (c *Config) IgnoreFilter() *filter.Filter {
	var rules []filter.Rule
	for _, e := range c.Ignore {
		rules = append(rules, filter.Rule{
			MinPort:   e.MinPort,
			MaxPort:   e.MaxPort,
			Protocols: e.Protocols,
		})
	}
	return filter.New(rules)
}
