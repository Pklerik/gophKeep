// Package server provides configuration for the GophKeeper server.
package server

import (
	"flag"
	"os"
	"time"
)

// Config represents server configuration.
type Config struct {
	ServerAddress string        `json:"server_address"`
	DatabasePath  string        `json:"database_path"`
	SecretKey     string        `json:"secret_key"`
	LogLevel      string        `json:"log_level"`
	Timeout       time.Duration `json:"timeout"`
	TLS           bool          `json:"enable_https"`
	CertFile      string        `json:"cert_file"`
	KeyFile       string        `json:"key_file"`
}

// LoadConfig loads server configuration from environment variables and flags.
func LoadConfig() Config {
	cfg := Config{
		ServerAddress: "localhost:8080",
		DatabasePath:  "./gophkeeper.db",
		SecretKey:     "default-secret-key",
		LogLevel:      "info",
		Timeout:       30 * time.Second,
		TLS:           false,
	}

	// Load from environment variables
	if addr := os.Getenv("SERVER_ADDRESS"); addr != "" {
		cfg.ServerAddress = addr
	}
	if dbPath := os.Getenv("DATABASE_PATH"); dbPath != "" {
		cfg.DatabasePath = dbPath
	}
	if key := os.Getenv("SECRET_KEY"); key != "" {
		cfg.SecretKey = key
	}
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}
	if tls := os.Getenv("ENABLE_HTTPS"); tls == "true" {
		cfg.TLS = true
	}

	// Load from command line flags
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "Server address")
	flag.StringVar(&cfg.DatabasePath, "d", cfg.DatabasePath, "Database path")
	flag.StringVar(&cfg.SecretKey, "k", cfg.SecretKey, "Secret key")
	flag.StringVar(&cfg.LogLevel, "l", cfg.LogLevel, "Log level")
	flag.BoolVar(&cfg.TLS, "s", cfg.TLS, "Enable TLS")
	flag.StringVar(&cfg.CertFile, "c", "", "Certificate file (for TLS)")
	flag.StringVar(&cfg.KeyFile, "p", "", "Key file (for TLS)")
	flag.Parse()

	return cfg
}
