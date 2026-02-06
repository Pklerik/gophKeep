// Package server provides router configuration for the GophKeeper server.
package server

import (
	"context"
	"fmt"
	"net/http"

	config "github.com/Pklerik/gophKeep/internal/config/server"
	handler "github.com/Pklerik/gophKeep/internal/handler/server"
	"github.com/Pklerik/gophKeep/internal/logger"
	"github.com/Pklerik/gophKeep/internal/middleware"
	"github.com/Pklerik/gophKeep/internal/migrations"

	"github.com/Pklerik/gophKeep/internal/repository"
	postgresrepository "github.com/Pklerik/gophKeep/internal/repository/postgres"
	"github.com/Pklerik/gophKeep/internal/service"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
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
	db, err := repository.ConnectDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = migrations.MakeMigrations(ctx, db, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make migrations: %w", err)
	}

	// Create repositories
	var (
		userRepo   repository.UserRepositoryInterface   = postgresrepository.NewUserRepository(db)
		secretRepo repository.SecretRepositoryInterface = postgresrepository.NewSecretRepository(db)
	)

	// Create services
	userService := service.NewUserService(userRepo, cfg.PasswordSalt)
	secretService := service.NewSecretService(secretRepo)

	// Create handlers
	h := handler.NewHandler(userService, secretService)

	// Create router
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(
			chimiddleware.RequestID,
			chimiddleware.RealIP,
			chimiddleware.Logger,
			chimiddleware.Recoverer,
			chimiddleware.Compress(5),
			chimiddleware.Timeout(cfg.Timeout),
		)
		r.Route("/", func(r chi.Router) {
			r.Route("/api/v1", func(r chi.Router) {
				r.Route("/auth", func(r chi.Router) {
					r.Post("/register", h.Register)
					r.Post("/login", h.Login)
				})
				r.Group(func(r chi.Router) {
					r.Use(
						middleware.AuthMiddleware(userService),
					)
					r.Route("/secrets", func(r chi.Router) {
						r.Get("/", h.ListSecrets)
						r.Post("/", h.CreateSecret)
						r.Get("/{id}", h.GetSecret)
						r.Put("/{id}", h.UpdateSecret)
						r.Delete("/{id}", h.DeleteSecret)
					})
				})
			})
			r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, `{"status":"ok"}`)
			})
			r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, `{"version":"1.0.0"}`)
			})
		})
	})

	logger.Sugar.Info("Router configured successfully")

	return r, nil
}
