// Package client provides the client entry point for the GophKeeper application.
package client

import (
	"log"

	"github.com/Pklerik/gophKeep/internal/app/client"
	config "github.com/Pklerik/gophKeep/internal/config/client"
	"github.com/Pklerik/gophKeep/internal/logger"
)

// Run starts the client application.
func Run() {
	cfg := config.LoadConfig()
	if err := logger.Initialize(cfg.LoggerConfig); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	client.StartApp(cfg)
}
