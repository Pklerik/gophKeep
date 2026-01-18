// Package postgresrepository_test tests the repository package.
package postgresrepository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	dbconfig "github.com/Pklerik/gophKeep/internal/config/db"
	config "github.com/Pklerik/gophKeep/internal/config/server"
	"github.com/Pklerik/gophKeep/internal/migrations"
	"github.com/Pklerik/gophKeep/internal/models"
	"github.com/Pklerik/gophKeep/internal/repository"
	postgresrepository "github.com/Pklerik/gophKeep/internal/repository/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// RepositoryTestSuite represents the test suite for repositories.
type RepositoryTestSuite struct {
	suite.Suite
	userRepo   *postgresrepository.UserRepository
	secretRepo *postgresrepository.SecretRepository
}

// SetupSuite sets up the test suite.
func (suite *RepositoryTestSuite) SetupSuite() {

	cfg := config.Config{
		DatabaseURL: dbconfig.Config{
			Dialect:  "postgres",
			Host:     "localhost",
			Port:     "5432",
			User:     "gophkeeper_test",
			Password: "secure_password",
			Database: "gophkeeper_test",
			Options:  map[string]string{"sslmode": "disable", "search_path": "gophkeeper"},
		},
	}
	db, err := repository.ConnectDB(cfg)
	suite.NoError(err)

	db.Exec("DROP SCHEMA IF EXISTS \"%s\" CASCADE;", cfg.DatabaseURL.GetOptions()["search_path"])

	migrations.MakeMigrations(context.Background(), db, cfg.DatabaseURL)

	suite.userRepo = postgresrepository.NewUserRepository(db)
	suite.secretRepo = postgresrepository.NewSecretRepository(db)

	suite.T().Cleanup(func() {
		db.Exec("DROP SCHEMA IF EXISTS \"%s\" CASCADE;", cfg.DatabaseURL.GetOptions()["search_path"])
		db.Close()
	})
}

// TestCreateUser tests creating a user.
func (suite *RepositoryTestSuite) TestCreateUser() {
	id := uuid.New().String()
	err := suite.userRepo.CreateUser(id, "testuser", "hashedpassword")
	assert.NoError(suite.T(), err)

	user, err := suite.userRepo.GetUserByID(id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "testuser", user.Username)
}

// TestGetUserByUsername tests retrieving a user by username.
func (suite *RepositoryTestSuite) TestGetUserByUsername() {
	id := uuid.New().String()
	suite.userRepo.CreateUser(id, "findme", "hashedpassword")

	user, err := suite.userRepo.GetUserByUsername("findme")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), id, user.ID)
}

// TestGetUserByUsernameNotFound tests retrieving non-existent user.
func (suite *RepositoryTestSuite) TestGetUserByUsernameNotFound() {
	_, err := suite.userRepo.GetUserByUsername("nonexistent")
	assert.Error(suite.T(), err)
}

// TestCreateSecret tests creating a secret.
func (suite *RepositoryTestSuite) TestCreateSecret() {
	userID := uuid.New().String()
	suite.userRepo.CreateUser(userID, "user", "hash")

	secret := &models.Secret{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      models.SecretTypeCredentials,
		Title:     "Test Secret",
		Data:      []byte("secret_data"),
		Metadata:  "test metadata",
		Version:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.secretRepo.CreateSecret(secret)
	assert.NoError(suite.T(), err)

	retrieved, err := suite.secretRepo.GetSecretByID(secret.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), secret.Title, retrieved.Title)
}

// TestGetUserSecrets tests retrieving all user secrets.
func (suite *RepositoryTestSuite) TestGetUserSecrets() {
	userID := uuid.New().String()
	suite.userRepo.CreateUser(userID, "user2", "hash")

	for i := 0; i < 3; i++ {
		secret := &models.Secret{
			ID:        uuid.New().String(),
			UserID:    userID,
			Type:      models.SecretTypeText,
			Title:     fmt.Sprintf("Secret %d", i),
			Data:      []byte("data"),
			Metadata:  "meta",
			Version:   1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		suite.secretRepo.CreateSecret(secret)
	}

	secrets, err := suite.secretRepo.GetUserSecrets(userID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 3, len(secrets))
}

// TestUpdateSecret tests updating a secret.
func (suite *RepositoryTestSuite) TestUpdateSecret() {
	userID := uuid.New().String()
	suite.userRepo.CreateUser(userID, "user3", "hash")

	secret := &models.Secret{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      models.SecretTypeCredentials,
		Title:     "Original",
		Data:      []byte("data"),
		Metadata:  "meta",
		Version:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.secretRepo.CreateSecret(secret)

	secret.Title = "Updated"
	secret.Data = []byte("updated_data")

	err := suite.secretRepo.UpdateSecret(secret)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, secret.Version)

	retrieved, _ := suite.secretRepo.GetSecretByID(secret.ID)
	assert.Equal(suite.T(), "Updated", retrieved.Title)
	assert.Equal(suite.T(), 2, retrieved.Version)
}

// TestDeleteSecret tests deleting a secret.
func (suite *RepositoryTestSuite) TestDeleteSecret() {
	userID := uuid.New().String()
	suite.userRepo.CreateUser(userID, "user4", "hash")

	secret := &models.Secret{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      models.SecretTypeText,
		Title:     "To Delete",
		Data:      []byte("data"),
		Metadata:  "meta",
		Version:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.secretRepo.CreateSecret(secret)

	err := suite.secretRepo.DeleteSecret(secret.ID, userID)
	assert.NoError(suite.T(), err)

	_, err = suite.secretRepo.GetSecretByID(secret.ID)
	assert.Error(suite.T(), err)
}

// TestDeleteSecretWrongUser tests deleting a secret with wrong user ID.
func (suite *RepositoryTestSuite) TestDeleteSecretWrongUser() {
	userID1 := uuid.New().String()
	userID2 := uuid.New().String()
	suite.userRepo.CreateUser(userID1, "user5", "hash")
	suite.userRepo.CreateUser(userID2, "user6", "hash")

	secret := &models.Secret{
		ID:        uuid.New().String(),
		UserID:    userID1,
		Type:      models.SecretTypeText,
		Title:     "Secret",
		Data:      []byte("data"),
		Metadata:  "meta",
		Version:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.secretRepo.CreateSecret(secret)

	err := suite.secretRepo.DeleteSecret(secret.ID, userID2)
	assert.Error(suite.T(), err)
}

// TestRepositoryTestSuite runs the test suite.
func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
