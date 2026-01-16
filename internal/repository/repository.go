// Package repository provides data access layer for the application.
package repository

import (
	"time"

	"github.com/Pklerik/gophKeep/internal/models"
)

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
