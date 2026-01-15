// Package server provides router configuration for the GophKeeper server.
package server

import (
	"context"
	"fmt"
	"net/http"

	config "github.com/Pklerik/gophKeep/internal/config/server"
	"github.com/Pklerik/gophKeep/internal/dbinit"
	"github.com/Pklerik/gophKeep/internal/handler"
	"github.com/Pklerik/gophKeep/internal/logger"
	"github.com/Pklerik/gophKeep/internal/repository"
	"github.com/Pklerik/gophKeep/internal/service"
	"github.com/gorilla/mux"
)

// ServerRouter represents the server router.
type ServerRouter struct{}

// NewServerRouter creates a new server router.
func NewServerRouter() *ServerRouter {
	return &ServerRouter{}
}

// ConfigureRouter configures the HTTP router with all routes.
func (sr *ServerRouter) ConfigureRouter(ctx context.Context, cfg config.Config) (http.Handler, error) {
	// Initialize database
	db, err := dbinit.InitDB(cfg.DatabasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Create repositories
	userRepo := repository.NewUserRepository(db)
	secretRepo := repository.NewSecretRepository(db)

	// Create services
	userService := service.NewUserService(userRepo)
	secretService := service.NewSecretService(secretRepo)

	// Create handlers
	h := handler.NewHandler(userService, secretService)

	// Create router
	router := mux.NewRouter()

	// Authentication routes
	router.HandleFunc("/api/v1/auth/register", h.Register).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/auth/login", h.Login).Methods(http.MethodPost)

	// Secret routes
	router.HandleFunc("/api/v1/secrets", h.ListSecrets).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/secrets", h.CreateSecret).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/secrets/get", h.GetSecret).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/secrets/update", h.UpdateSecret).Methods(http.MethodPut)
	router.HandleFunc("/api/v1/secrets/delete", h.DeleteSecret).Methods(http.MethodDelete)

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok"}`)
	}).Methods(http.MethodGet)

	logger.Sugar.Info("Router configured successfully")

	return router, nil
}
