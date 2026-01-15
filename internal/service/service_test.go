// Package service_test tests the service package.
package service

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Pklerik/gophKeep/internal/dbinit"
	"github.com/Pklerik/gophKeep/internal/models"
	"github.com/Pklerik/gophKeep/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ServiceTestSuite represents the test suite for services.
type ServiceTestSuite struct {
	suite.Suite
	db            *sql.DB
	userService   *UserService
	secretService *SecretService
	userRepo      *repository.UserRepository
	secretRepo    *repository.SecretRepository
}

// SetupSuite sets up the test suite.
func (suite *ServiceTestSuite) SetupSuite() {
	tmpFile := fmt.Sprintf("%d.db", time.Now().UnixNano())
	db, err := dbinit.InitDB(tmpFile)
	suite.NoError(err)
	suite.db = db

	suite.userRepo = repository.NewUserRepository(db)
	suite.secretRepo = repository.NewSecretRepository(db)

	suite.userService = NewUserService(suite.userRepo)
	suite.secretService = NewSecretService(suite.secretRepo)

	suite.T().Cleanup(func() {
		db.Close()
		os.Remove(tmpFile)
	})
}

// TestRegisterUser tests user registration.
func (suite *ServiceTestSuite) TestRegisterUser() {
	user, err := suite.userService.RegisterUser("testuser", "password123")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), "testuser", user.Username)
}

// TestAuthenticateUser tests user authentication.
func (suite *ServiceTestSuite) TestAuthenticateUser() {
	suite.userService.RegisterUser("authuser", "password123")

	user, err := suite.userService.AuthenticateUser("authuser", "password123")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), "authuser", user.Username)
}

// TestCreateSecret tests creating a secret.
func (suite *ServiceTestSuite) TestCreateSecret() {
	user, _ := suite.userService.RegisterUser("secretuser", "password123")

	secret, err := suite.secretService.CreateSecret(user.ID, models.SecretTypeCredentials,
		"My Password", "login:password", "website.com")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), secret)
	assert.Equal(suite.T(), "My Password", secret.Title)
}

// TestListSecrets tests listing user secrets.
func (suite *ServiceTestSuite) TestListSecrets() {
	user, _ := suite.userService.RegisterUser("listuser", "password123")

	suite.secretService.CreateSecret(user.ID, models.SecretTypeCredentials,
		"Secret 1", "data1", "meta1")
	suite.secretService.CreateSecret(user.ID, models.SecretTypeText,
		"Secret 2", "data2", "meta2")

	secrets, err := suite.secretService.ListSecrets(user.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, len(secrets))
}

// TestGetSecret tests retrieving a specific secret.
func (suite *ServiceTestSuite) TestGetSecret() {
	user, _ := suite.userService.RegisterUser("getuser", "password123")
	created, _ := suite.secretService.CreateSecret(user.ID, models.SecretTypeCredentials,
		"My Secret", "secret_data", "metadata")

	retrieved, err := suite.secretService.GetSecret(created.ID, user.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), created.ID, retrieved.ID)
	assert.Equal(suite.T(), "My Secret", retrieved.Title)
}

// TestUpdateSecret tests updating a secret.
func (suite *ServiceTestSuite) TestUpdateSecret() {
	user, _ := suite.userService.RegisterUser("updateuser", "password123")
	created, _ := suite.secretService.CreateSecret(user.ID, models.SecretTypeCredentials,
		"Original", "data", "meta")

	updated, err := suite.secretService.UpdateSecret(created.ID, user.ID,
		models.SecretTypeText, "Updated", "new_data", "new_meta")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated", updated.Title)
	assert.Equal(suite.T(), 2, updated.Version)
}

// TestDeleteSecret tests deleting a secret.
func (suite *ServiceTestSuite) TestDeleteSecret() {
	user, _ := suite.userService.RegisterUser("deleteuser", "password123")
	created, _ := suite.secretService.CreateSecret(user.ID, models.SecretTypeCredentials,
		"To Delete", "data", "meta")

	err := suite.secretService.DeleteSecret(created.ID, user.ID)
	assert.NoError(suite.T(), err)

	_, err = suite.secretService.GetSecret(created.ID, user.ID)
	assert.Error(suite.T(), err)
}

// TestServiceTestSuite runs the test suite.
func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
