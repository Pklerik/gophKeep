// Package client provides configuration for the GophKeeper client.
package client

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
)

// Config represents client configuration.
type Config struct {
	ServerURL    string
	ClientDBPath string
	LoggerConfig zap.Config
	Timeout      time.Duration
	LogLevel     string
}

// LoadConfig loads client configuration from environment variables and flags.
func LoadConfig() Config {
	cfg := Config{
		ServerURL:    "http://localhost:8080",
		ClientDBPath: filepath.Join(os.ExpandEnv("$HOME"), ".gophkeeper", "client.db"),
		Timeout:      10 * time.Second,
		LogLevel:     "info",
		LoggerConfig: zap.Config{
			Level:            zap.NewAtomicLevel(),
			Development:      false,
			Encoding:         "json",
			EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		},
	}

	// Load from environment variables
	if url := os.Getenv("SERVER_URL"); url != "" {
		cfg.ServerURL = url
	}
	if dbPath := os.Getenv("CLIENT_DB_PATH"); dbPath != "" {
		cfg.ClientDBPath = dbPath
	}

	// Load from command line flags
	flag.StringVar(&cfg.ServerURL, "u", cfg.ServerURL, "Server URL")
	flag.StringVar(&cfg.ClientDBPath, "d", cfg.ClientDBPath, "Client database path")
	flag.StringVar(&cfg.LogLevel, "l", cfg.LogLevel, "Log level")
	flag.CommandLine.Parse(os.Args[2:])
	// Parse log level and set it in the logger configuration
	lvl, err := zap.ParseAtomicLevel(cfg.LogLevel)
	if err == nil {
		cfg.LoggerConfig.Level.SetLevel(lvl.Level())
	}
	return cfg
}
