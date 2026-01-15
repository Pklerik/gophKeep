// Package service provides business logic for user and secret management.
package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Pklerik/gophKeep/internal/cryptography"
	"github.com/Pklerik/gophKeep/internal/models"
	"github.com/Pklerik/gophKeep/internal/repository"
	"github.com/google/uuid"
)

// UserService handles user-related business logic.
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService creates a new user service.
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// RegisterUser registers a new user.
func (s *UserService) RegisterUser(username, password string) (*models.User, error) {
	if username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}

	// Check if user already exists
	existing, err := s.repo.GetUserByUsername(username)
	if err == nil && existing != nil {
		return nil, errors.New("user already exists")
	}

	id := uuid.New().String()
	hash := cryptography.HashPassword(password)

	err = s.repo.CreateUser(id, username, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return &models.User{
		ID:        id,
		Username:  username,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// AuthenticateUser authenticates a user with username and password.
func (s *UserService) AuthenticateUser(username, password string) (*models.User, error) {
	if username == "" || password == "" {
		return nil, errors.New("username and password are required")
	}

	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	if !cryptography.VerifyPassword(password, user.PasswordHash) {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

// GetUser retrieves a user by ID.
func (s *UserService) GetUser(userID string) (*models.User, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// SecretService handles secret-related business logic.
type SecretService struct {
	repo *repository.SecretRepository
}

// NewSecretService creates a new secret service.
func NewSecretService(repo *repository.SecretRepository) *SecretService {
	return &SecretService{repo: repo}
}

// CreateSecret creates a new secret for a user.
func (s *SecretService) CreateSecret(userID string, secretType models.SecretType,
	title, data, metadata string) (*models.Secret, error) {
	if userID == "" || title == "" || data == "" {
		return nil, errors.New("user_id, title, and data are required")
	}

	if secretType == "" {
		return nil, errors.New("secret type is required")
	}

	id := uuid.New().String()
	now := time.Now()

	secret := &models.Secret{
		ID:        id,
		UserID:    userID,
		Type:      secretType,
		Title:     title,
		Data:      []byte(data),
		Metadata:  metadata,
		Version:   1,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := s.repo.CreateSecret(secret)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret: %w", err)
	}

	return secret, nil
}

// GetSecret retrieves a secret by ID.
func (s *SecretService) GetSecret(id, userID string) (*models.Secret, error) {
	secret, err := s.repo.GetSecretByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}

	if secret.UserID != userID {
		return nil, errors.New("access denied")
	}

	return secret, nil
}

// ListSecrets lists all secrets for a user.
func (s *SecretService) ListSecrets(userID string) ([]models.Secret, error) {
	secrets, err := s.repo.GetUserSecrets(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}

	if secrets == nil {
		return []models.Secret{}, nil
	}

	return secrets, nil
}

// UpdateSecret updates a secret.
func (s *SecretService) UpdateSecret(id, userID string, secretType models.SecretType,
	title, data, metadata string) (*models.Secret, error) {
	secret, err := s.GetSecret(id, userID)
	if err != nil {
		return nil, err
	}

	secret.Type = secretType
	secret.Title = title
	secret.Data = []byte(data)
	secret.Metadata = metadata

	err = s.repo.UpdateSecret(secret)
	if err != nil {
		return nil, fmt.Errorf("failed to update secret: %w", err)
	}

	return secret, nil
}

// DeleteSecret deletes a secret.
func (s *SecretService) DeleteSecret(id, userID string) error {
	err := s.repo.DeleteSecret(id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}
	return nil
}

// GetSecretsSince retrieves secrets updated after a certain time.
func (s *SecretService) GetSecretsSince(userID string, since time.Time) ([]models.Secret, error) {
	secrets, err := s.repo.GetSecretsSince(userID, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get secrets since: %w", err)
	}

	if secrets == nil {
		return []models.Secret{}, nil
	}

	return secrets, nil
}
