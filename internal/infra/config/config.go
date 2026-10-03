package config

import (
	"adbcon/api"
	"log/slog"
	"net/url"
	"os"
	"time"

	yaml "gopkg.in/yaml.v3"
)

// Config is an aplication configurations.
type Config struct {
	Log    LogConfig    `yaml:"log"`
	Server ServerConfig `yaml:"server"`
}

// LogConfig is a log configurations
type LogConfig struct {
	FilePath string `yaml:"filepath"`
}

// ServerConfig is a server configurations
type ServerConfig struct {
	// Bind is an address to bind server linstener.
	Bind string `yaml:"bind"`
	// Port is a server listen port.
	Port string    `yaml:"port"`
	SSE  SSEConfig `yaml:"sse"`
}

// SSEConfig is a SSE configurations.
type SSEConfig struct {
	// Interval is a SSE send maximum interval.
	// Even if buffer is not filled, send SSE data if expiry this interval.
	Interval int `yanl:"interval"`
	// Limit is a same as buffer size.
	// Send SSE data when buffered line over limit.
	Limit int `yaml:"limit"`
}

// Read reads configurations from file.
func (c *Config) Read(path string) error {
	// read value from openapi difinition
	urlServer, err := url.Parse(api.ServerUrlDefault)
	if err != nil {
		slog.Error("url parse error", "err", err, "url", api.ServerUrlDefault)
		return err
	}
	c.Server.Bind = urlServer.Hostname() // default is a definition of openapi
	c.Server.Port = urlServer.Port()     // default is a definition of openapi

	// check file existing
	// Change spec. : It did not error if file is not existing,
	// because application uses default configuration if config file is not existing.
	if _, err := os.Stat(path); err != nil {
		slog.Info("used default config")
		return nil
	}

	f, err := os.ReadFile(path)
	if err != nil {
		slog.Error("config read file error", "err", err, "path", path)
		return err
	}

	if err := yaml.Unmarshal(f, c); err != nil {
		slog.Error("yaml unmarshal error", "err", err, "path", path)
		return err
	}

	return nil
}

// LogFilePath provides a file path that outputs information log.
// Default value is `adbcon.log`
func (c Config) LogFilePath() string {
	if c.Log.FilePath == "" {
		return "adbcon.log"
	}
	return c.Log.FilePath
}

// ServerBind resolves an address to bind server listener.
// Default value is definition of openapi if not specified in file.
func (c Config) ServerBind() string {
	return c.Server.Bind
}

// ServerPort is a server listen port.
// Default value is definition of openapi if not specified in file.
func (c Config) ServerPort() string {
	return c.Server.Port
}

// ServerSseInterval returns interval time by milliseconds.
func (c Config) ServerSseInterval() time.Duration {
	if c.Server.SSE.Interval == 0 {
		return 100 * time.Millisecond // TODO: supports 0, it measns send immediate
	}
	return time.Duration(c.Server.SSE.Interval) * time.Millisecond
}

// ServerSseLimit returns maximum number to store SSE data in buffer.
func (c Config) ServerSseLimit() int {
	if c.Server.SSE.Limit == 0 {
		return 30 // TODO: supports 0, it measns send immediate
	}
	return c.Server.SSE.Limit
}
