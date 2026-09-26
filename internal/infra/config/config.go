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
	Port string `yaml:"port"`
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
