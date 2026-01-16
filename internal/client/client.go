// Package client provides HTTP client for communicating with the GophKeeper server.
package client

import (
	"github.com/Pklerik/gophKeep/internal/models"
)

type IClient interface {
	Register(username, password string) (*models.AuthResponse, error)
	Login(username, password string) (*models.AuthResponse, error)
	CreateSecret(secretType models.SecretType, title, data, metadata string) (*models.Secret, error)
	GetSecret(id string) (*models.Secret, error)
	ListSecrets() ([]models.Secret, error)
	UpdateSecret(id string, secretType models.SecretType, title, data, metadata string) (*models.Secret, error)
	DeleteSecret(id string) error
	SetToken(token string)
}
