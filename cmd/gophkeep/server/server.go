// Package server provides the server entry point for the GophKeeper application.
package server

import (
	"log"

	"github.com/Pklerik/gophKeep/internal/app/server"
	config "github.com/Pklerik/gophKeep/internal/config/server"
	"github.com/Pklerik/gophKeep/internal/logger"
)

// Run starts the server application with the given configuration.
func Run() {

	cfg := config.LoadConfig()

	if err := logger.Initialize(cfg.LoggerConfig); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	server.StartApp(cfg)
}
