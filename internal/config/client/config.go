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
	ServerURL      string
	ClientDBPath   string
	LoggerConfig   zap.Config
	Timeout        time.Duration
	LogLevel       string
	EncryptionKey  []byte
	EncryptionSalt []byte
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
		EncryptionKey:  []byte("don't use base encryption password"), // In production, use a secure, random key
		EncryptionSalt: []byte("don't use base salt"),                // In production, use a secure, random salt
	}

	// Load from environment variables
	if url := os.Getenv("SERVER_URL"); url != "" {
		cfg.ServerURL = url
	}
	if dbPath := os.Getenv("CLIENT_DB_PATH"); dbPath != "" {
		cfg.ClientDBPath = dbPath
	}

	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	if encryptionKey != "" {
		cfg.EncryptionKey = []byte(encryptionKey)
	}

	encryptionSalt := os.Getenv("ENCRYPTION_SALT")
	if encryptionSalt != "" {
		cfg.EncryptionSalt = []byte(encryptionSalt)
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
