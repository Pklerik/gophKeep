// Package repository provides data access layer for the application.
package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	config "github.com/Pklerik/gophKeep/internal/config/server"
	"github.com/Pklerik/gophKeep/internal/logger"
	"github.com/Pklerik/gophKeep/internal/migrations"
	"github.com/Pklerik/gophKeep/internal/models"
)

var ErrCollectingDBConf = errors.New("error collecting database configuration")

type UserRepositoryInterface interface {
	CreateUser(id, username, passwordHash string) error
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
}

type SecretRepositoryInterface interface {
	CreateSecret(secret *models.Secret) error
	GetSecretByID(id string) (*models.Secret, error)
	GetUserSecrets(userID string) ([]models.Secret, error)
	UpdateSecret(secret *models.Secret) error
	DeleteSecret(id, userID string) error
	GetSecretsSince(userID string, since time.Time) ([]models.Secret, error)
}

// ConnectDB connecting to DB.
func ConnectDB(cfg config.Config) (*sql.DB, error) {
	dbConf := cfg.DatabaseURL

	if err := dbConf.Valid(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCollectingDBConf, err)
	}

	if os.Getenv("GOOSE_DRIVER") == "" {
		if err := os.Setenv("GOOSE_DRIVER", dbConf.Dialect); err != nil {
			return nil, fmt.Errorf("cant set env variable: %w", err)
		}
	}

	if os.Getenv("GOOSE_DBSTRING") == "" {
		if err := os.Setenv("GOOSE_DBSTRING", dbConf.GetConnString()); err != nil {
			return nil, fmt.Errorf("cant set env variable: %w", err)
		}
	}

	if os.Getenv("GOOSE_MIGRATION_DIR") == "" {
		if err := os.Setenv("GOOSE_MIGRATION_DIR", migrations.MigrationDir()); err != nil {
			return nil, fmt.Errorf("cant set env variable: %w", err)
		}
	}

	logger.Sugar.Infof("ConnString: Database: %s, User: %s, Options: %v",
		dbConf.Database,
		dbConf.User,
		dbConf.Options,
	)

	ps := dbConf.GetConnString()

	db, err := sql.Open("pgx", ps)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to DB: %w", err)
	}

	return db, nil
}
