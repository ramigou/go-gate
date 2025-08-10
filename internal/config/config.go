package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Upstreams []UpstreamConfig `yaml:"upstreams"`
	Routes    []RouteConfig   `yaml:"routes"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type UpstreamConfig struct {
	Name   string `yaml:"name"`
	URL    string `yaml:"url"`
	Weight int    `yaml:"weight"`
	Health HealthConfig `yaml:"health"`
}

type HealthConfig struct {
	Path     string `yaml:"path"`
	Interval string `yaml:"interval"`
	Timeout  string `yaml:"timeout"`
}

type RouteConfig struct {
	Name      string            `yaml:"name"`
	Match     MatchConfig       `yaml:"match"`
	Upstream  string            `yaml:"upstream"`
	Upstreams []string          `yaml:"upstreams"`
}

type MatchConfig struct {
	Host string `yaml:"host"`
	Path string `yaml:"path"`
}

func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", filename, err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", filename, err)
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

func (c *Config) validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if len(c.Upstreams) == 0 {
		return fmt.Errorf("at least one upstream must be configured")
	}

	upstreamNames := make(map[string]bool)
	for _, upstream := range c.Upstreams {
		if upstream.Name == "" {
			return fmt.Errorf("upstream name cannot be empty")
		}
		if upstream.URL == "" {
			return fmt.Errorf("upstream URL cannot be empty")
		}
		if upstreamNames[upstream.Name] {
			return fmt.Errorf("duplicate upstream name: %s", upstream.Name)
		}
		upstreamNames[upstream.Name] = true
	}

	for _, route := range c.Routes {
		if route.Name == "" {
			return fmt.Errorf("route name cannot be empty")
		}
		if route.Upstream == "" && len(route.Upstreams) == 0 {
			return fmt.Errorf("route %s must specify upstream or upstreams", route.Name)
		}
		if route.Upstream != "" && !upstreamNames[route.Upstream] {
			return fmt.Errorf("route %s references unknown upstream: %s", route.Name, route.Upstream)
		}
		for _, upstreamName := range route.Upstreams {
			if !upstreamNames[upstreamName] {
				return fmt.Errorf("route %s references unknown upstream: %s", route.Name, upstreamName)
			}
		}
	}

	return nil
}