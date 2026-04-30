package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the agent configuration
type Config struct {
	ServerID      string           `yaml:"server_id"`
	Endpoint      string           `yaml:"endpoint"`
	Interval      time.Duration    `yaml:"interval"`
	Metrics       MetricsConfig    `yaml:"metrics"`
	Auth          AuthConfig       `yaml:"auth"`
	Ports         []PortCheck      `yaml:"ports"`
	LatencyChecks []LatencyCheck   `yaml:"latency_checks"`
}

// MetricsConfig defines which metrics to collect
type MetricsConfig struct {
	CPU         bool `yaml:"cpu"`
	Memory      bool `yaml:"memory"`
	Disk        bool `yaml:"disk"`
	Network     bool `yaml:"network"`
	Connections bool `yaml:"connections"`
}

// AuthConfig defines authentication settings
type AuthConfig struct {
	TokenFile string `yaml:"token_file"`
}

// PortCheck represents a port monitoring configuration
type PortCheck struct {
	Port    int    `yaml:"port"`
	Name    string `yaml:"name"`
	Host    string `yaml:"host"`
	Timeout string `yaml:"timeout"`
}

// LatencyCheck represents a latency monitoring configuration
type LatencyCheck struct {
	Host    string `yaml:"host"`
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Port    int    `yaml:"port"`
	Timeout string `yaml:"timeout"`
	URL     string `yaml:"url"`
}

// Load reads and parses the configuration file
func Load(filename string) (*Config, error) {
	// Set default values
	config := &Config{
		ServerID: "",
		Endpoint: "https://api.yourapp.com",
		Interval: 5 * time.Second,
		Metrics: MetricsConfig{
			CPU:         true,
			Memory:      true,
			Disk:        true,
			Network:     true,
			Connections: false, // Optional by default
		},
		Auth: AuthConfig{
			TokenFile: "./agent_token",
		},
		Ports:         []PortCheck{},
		LatencyChecks: []LatencyCheck{},
	}

	// Read config file
	data, err := os.ReadFile(filename)
	if err != nil {
		// If config file doesn't exist, create it with defaults
		if os.IsNotExist(err) {
			fmt.Printf("Config file %s not found, creating with defaults...\n", filename)
			if err := Save(config, filename); err != nil {
				return nil, fmt.Errorf("failed to create default config: %w", err)
			}
			return config, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// Save writes the configuration to a file
func Save(config *Config, filename string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Endpoint == "" {
		return fmt.Errorf("endpoint cannot be empty")
	}

	if c.Interval <= 0 {
		return fmt.Errorf("interval must be positive")
	}

	if c.Auth.TokenFile == "" {
		return fmt.Errorf("token_file cannot be empty")
	}

	return nil
}