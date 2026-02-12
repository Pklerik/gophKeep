// Package postgresrepository_test tests the repository package.
package postgresrepository_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

	pwd, err := os.Getwd()
	assert.NoError(suite.T(), err)
	composePath := filepath.Join(pwd, "testdata", "docker-compose-tests.yml")
	fmt.Printf("Using docker-compose file: %s\n", composePath)
	cmd := exec.Command("docker-compose", "-f", composePath, "up", "-d")
	err = cmd.Run()
	assert.NoError(suite.T(), err)

	// Wait for the database to be ready
	for i := 0; i < 10; i++ {
		conn, err := repository.ConnectDB(cfg)
		if err == nil {
			conn.Close()
			time.Sleep(2 * time.Second)
			break
		}
		time.Sleep(2 * time.Second)
	}

	db, err := repository.ConnectDB(cfg)
	suite.NoError(err)

	db.Exec("DROP SCHEMA IF EXISTS \"%s\" CASCADE;", cfg.DatabaseURL.GetOptions()["search_path"])

	migrations.MakeMigrations(context.Background(), db, cfg.DatabaseURL)

	suite.userRepo = postgresrepository.NewUserRepository(db)
	suite.secretRepo = postgresrepository.NewSecretRepository(db)

	suite.T().Cleanup(func() {

		db.Exec("DROP SCHEMA IF EXISTS \"%s\" CASCADE;", cfg.DatabaseURL.GetOptions()["search_path"])
		cmd := exec.Command("docker-compose", "-f", composePath, "down", "-v")
		err = cmd.Run()
		assert.NoError(suite.T(), err)
		db.Close()
	})
}

// TestCreateUser tests creating a user.
func (suite *RepositoryTestSuite) TestCreateUser() {
	id := uuid.New().String()
	username := fmt.Sprintf("testuser_%s", id[:8])
	err := suite.userRepo.CreateUser(id, username, "hashedpassword")
	assert.NoError(suite.T(), err)

	user, err := suite.userRepo.GetUserByID(id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), username, user.Username)
}

// TestGetUserByUsername tests retrieving a user by username.
func (suite *RepositoryTestSuite) TestGetUserByUsername() {
	id := uuid.New().String()
	username := fmt.Sprintf("findme_%s", uuid.New().String()[:8])
	suite.userRepo.CreateUser(id, username, "hashedpassword")

	user, err := suite.userRepo.GetUserByUsername(username)
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
	username := fmt.Sprintf("user_create_%s", uuid.New().String()[:8])
	suite.userRepo.CreateUser(userID, username, "hash")

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
	username := fmt.Sprintf("user_secrets_%s", uuid.New().String()[:8])
	suite.userRepo.CreateUser(userID, username, "hash")

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
	username := fmt.Sprintf("user_update_%s", uuid.New().String()[:8])
	suite.userRepo.CreateUser(userID, username, "hash")

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
	username := fmt.Sprintf("user_delete_%s", uuid.New().String()[:8])
	suite.userRepo.CreateUser(userID, username, "hash")

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
	username1 := fmt.Sprintf("user_del_wrong1_%s", uuid.New().String()[:8])
	username2 := fmt.Sprintf("user_del_wrong2_%s", uuid.New().String()[:8])
	suite.userRepo.CreateUser(userID1, username1, "hash")
	suite.userRepo.CreateUser(userID2, username2, "hash")

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
