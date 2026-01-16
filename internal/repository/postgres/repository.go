// Package postgresrepository provides data access layer for the application.
package postgresrepository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Pklerik/gophKeep/internal/models"
)

// UserRepository handles user-related database operations.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user in the database.
func (r *UserRepository) CreateUser(id, username, passwordHash string) error {
	query := `INSERT INTO users (id, username, password_hash) VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, id, username, passwordHash)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByUsername retrieves a user by username.
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, password_hash, created_at, updated_at FROM users WHERE username = ?`
	row := r.db.QueryRow(query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID.
func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	query := `SELECT id, username, password_hash, created_at, updated_at FROM users WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// SecretRepository handles secret-related database operations.
type SecretRepository struct {
	db *sql.DB
}

// NewSecretRepository creates a new secret repository.
func NewSecretRepository(db *sql.DB) *SecretRepository {
	return &SecretRepository{db: db}
}

// CreateSecret creates a new secret in the database.
func (r *SecretRepository) CreateSecret(secret *models.Secret) error {
	query := `INSERT INTO secrets (id, user_id, type, title, data, metadata, version) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, secret.ID, secret.UserID, secret.Type, secret.Title,
		secret.Data, secret.Metadata, secret.Version)
	if err != nil {
		return fmt.Errorf("failed to create secret: %w", err)
	}
	return nil
}

// GetSecretByID retrieves a secret by ID.
func (r *SecretRepository) GetSecretByID(id string) (*models.Secret, error) {
	query := `SELECT id, user_id, type, title, data, metadata, version, created_at, updated_at 
	          FROM secrets WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var secret models.Secret
	err := row.Scan(&secret.ID, &secret.UserID, &secret.Type, &secret.Title, &secret.Data,
		&secret.Metadata, &secret.Version, &secret.CreatedAt, &secret.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("secret not found")
		}
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}

	return &secret, nil
}

// GetUserSecrets retrieves all secrets for a user.
func (r *SecretRepository) GetUserSecrets(userID string) ([]models.Secret, error) {
	query := `SELECT id, user_id, type, title, data, metadata, version, created_at, updated_at 
	          FROM secrets WHERE user_id = ? ORDER BY updated_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user secrets: %w", err)
	}
	defer rows.Close()

	var secrets []models.Secret
	for rows.Next() {
		var secret models.Secret
		err := rows.Scan(&secret.ID, &secret.UserID, &secret.Type, &secret.Title, &secret.Data,
			&secret.Metadata, &secret.Version, &secret.CreatedAt, &secret.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan secret: %w", err)
		}
		secrets = append(secrets, secret)
	}

	return secrets, rows.Err()
}

// UpdateSecret updates an existing secret.
func (r *SecretRepository) UpdateSecret(secret *models.Secret) error {
	secret.UpdatedAt = time.Now()
	secret.Version++

	query := `UPDATE secrets SET type = ?, title = ?, data = ?, metadata = ?, version = ?, updated_at = ? 
	          WHERE id = ? AND user_id = ?`
	result, err := r.db.Exec(query, secret.Type, secret.Title, secret.Data, secret.Metadata,
		secret.Version, secret.UpdatedAt, secret.ID, secret.UserID)
	if err != nil {
		return fmt.Errorf("failed to update secret: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return errors.New("secret not found")
	}

	return nil
}

// DeleteSecret deletes a secret.
func (r *SecretRepository) DeleteSecret(id, userID string) error {
	query := `DELETE FROM secrets WHERE id = ? AND user_id = ?`
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return errors.New("secret not found")
	}

	return nil
}

// GetSecretsSince retrieves secrets updated after a certain time.
func (r *SecretRepository) GetSecretsSince(userID string, since time.Time) ([]models.Secret, error) {
	query := `SELECT id, user_id, type, title, data, metadata, version, created_at, updated_at 
	          FROM secrets WHERE user_id = ? AND updated_at > ? ORDER BY updated_at DESC`
	rows, err := r.db.Query(query, userID, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get secrets since: %w", err)
	}
	defer rows.Close()

	var secrets []models.Secret
	for rows.Next() {
		var secret models.Secret
		err := rows.Scan(&secret.ID, &secret.UserID, &secret.Type, &secret.Title, &secret.Data,
			&secret.Metadata, &secret.Version, &secret.CreatedAt, &secret.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan secret: %w", err)
		}
		secrets = append(secrets, secret)
	}

	return secrets, rows.Err()
}
