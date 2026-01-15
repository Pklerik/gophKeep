// Package client provides configuration for the GophKeeper client.
package client

import (
	"flag"
	"os"
	"path/filepath"
	"time"
)

// Config represents client configuration.
type Config struct {
	ServerURL    string
	ClientDBPath string
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

	return cfg
}
