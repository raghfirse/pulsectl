package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Endpoint represents a single HTTP endpoint to poll.
type Endpoint struct {
	Name     string        `yaml:"name"`
	URL      string        `yaml:"url"`
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
}

// Config holds the full application configuration.
type Config struct {
	Endpoints []Endpoint `yaml:"endpoints"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if len(c.Endpoints) == 0 {
		return fmt.Errorf("at least one endpoint must be defined")
	}
	for i, ep := range c.Endpoints {
		if ep.Name == "" {
			return fmt.Errorf("endpoint[%d]: name is required", i)
		}
		if ep.URL == "" {
			return fmt.Errorf("endpoint[%d]: url is required", i)
		}
		if ep.Interval <= 0 {
			c.Endpoints[i].Interval = 30 * time.Second
		}
		if ep.Timeout <= 0 {
			c.Endpoints[i].Timeout = 5 * time.Second
		}
	}
	return nil
}
