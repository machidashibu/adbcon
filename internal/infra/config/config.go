package config

import (
	"adbcon/api"
	"log/slog"
	"net/url"
	"os"

	yaml "gopkg.in/yaml.v3"
)

// Config is an aplication configurations.
type Config struct {
	Server ServerConfig `yaml:"server"`
}

// ServerConfig is a server configurations
type ServerConfig struct {
	// Bind is an address to bind server linstener.
	Bind string `yaml:"bind"`
	// Port is a server listen port.
	Port string `yaml:"port"`
}

// Read reads configurations from file.
func (c *Config) Read(path string) error {
	f, err := os.ReadFile(path)
	if err != nil {
		slog.Error("config read file error", "err", err, "path", path)
		return err
	}

	if err := yaml.Unmarshal(f, &c); err != nil {
		slog.Error("yaml unmarshal error", "err", err, "path", path)
		return err
	}

	// use definition of openapi if not specified both bind and port
	if c.Server.Bind == "" && c.Server.Port == "" {
		// read value from openapi difinition
		urlServer, err := url.Parse(api.ServerUrlDefault)
		if err != nil {
			slog.Error("url parse error", "err", err, "url", api.ServerUrlDefault)
			return err
		}
		c.Server.Bind = urlServer.Hostname()
		c.Server.Port = urlServer.Port()
	}

	return nil
}

// ServerBind resolves an address to bind server listener.
// Default value is definition of openapi if not specified both bind and port.
func (c Config) ServerBind() string {
	return c.Server.Bind
}

// ServerPort is a server listen port.
// Default value is definition of openapi if not specified both bind and port.
func (c Config) ServerPort() string {
	return c.Server.Port
}
