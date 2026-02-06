// Package server provides configuration for the GophKeeper server.
package server

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/Pklerik/gophKeep/internal/auth"
	dbconf "github.com/Pklerik/gophKeep/internal/config/db"
	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"
)

// Config represents server configuration.
type Config struct {
	ServerAddress string        `json:"server_address" env:"SERVER_ADDRESS"`
	DatabasePath  string        `json:"database_path" env:"DATABASE_PATH"`
	DatabaseURL   dbconf.Config `json:"database_url" env:"DATABASE_URL"`
	LoggerConfig  zap.Config    `json:"logger_config" env:"LOGGER_CONFIG"`
	secretKey     string        `json:"-" env:"SECRET_KEY"`
	LogLevel      string        `json:"log_level" env:"LOG_LEVEL"`
	Timeout       time.Duration `json:"timeout" env:"TIMEOUT"`
	TLS           bool          `json:"enable_https" env:"ENABLE_HTTPS"`
	CertFile      string        `json:"cert_file" env:"CERT_FILE"`
	KeyFile       string        `json:"key_file" env:"KEY_FILE"`
	PasswordSalt  []byte        `json:"-" env:"PASSWORD_SALT"`
}

// LoadConfig loads server configuration from environment variables and flags.
func LoadConfig() Config {
	cfg := Config{
		ServerAddress: "localhost:8080",
		DatabasePath:  "./gophkeeper.db",
		secretKey:     "default-secret-key",
		LogLevel:      "info",
		Timeout:       30 * time.Second,
		PasswordSalt:  []byte("default-password-salt"),
		DatabaseURL: dbconf.Config{
			Dialect: "sqlite3",
		},
		LoggerConfig: zap.Config{
			Level:            zap.NewAtomicLevel(),
			Development:      false,
			Encoding:         "json",
			EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		},
	}

	// Load from command line flags
	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "Server address")
	flag.StringVar(&cfg.DatabasePath, "d", cfg.DatabasePath, "Database path")
	flag.Var(&cfg.DatabaseURL, "dsn", "Database URL")
	flag.DurationVar(&cfg.Timeout, "t", cfg.Timeout, "Set timeout")
	flag.StringVar(&cfg.secretKey, "k", cfg.secretKey, "Secret key")
	flag.StringVar(&cfg.LogLevel, "l", cfg.LogLevel, "Log level")
	flag.BoolVar(&cfg.TLS, "s", cfg.TLS, "Enable TLS")
	flag.StringVar(&cfg.CertFile, "c", "", "Certificate file (for TLS)")
	flag.StringVar(&cfg.KeyFile, "p", "", "Key file (for TLS)")
	stringSalt := string(cfg.PasswordSalt)
	flag.StringVar(&stringSalt, "ps", stringSalt, "Password salt")
	cfg.PasswordSalt = []byte(stringSalt)
	flag.CommandLine.Parse(os.Args[2:])

	// Load from environment variables rewriting flag values if exists.
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	lvl, err := zap.ParseAtomicLevel(cfg.LogLevel)
	if err == nil {
		cfg.LoggerConfig.Level.SetLevel(lvl.Level())
	}

	// Set secret key for authentication.
	auth.SetSecretKey(cfg.secretKey)

	return cfg
}
