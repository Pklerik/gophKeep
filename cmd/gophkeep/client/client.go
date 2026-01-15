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
	if err := logger.Initialize("info"); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	cfg := config.LoadConfig()
	client.StartApp(cfg)
}
