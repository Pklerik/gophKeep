// Package service_test tests the service package.
package service

import (
	"errors"
	"testing"

	"github.com/Pklerik/gophKeep/internal/cryptography"
	"github.com/Pklerik/gophKeep/internal/models"
	"github.com/Pklerik/gophKeep/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var testSalt = []byte("test-salt")

// ServiceTestSuite represents the test suite for services.
type ServiceTestSuite struct {
	suite.Suite
	userService    *UserService
	secretService  *SecretService
	ctrl           *gomock.Controller
	mockUserRepo   *mocks.MockUserRepositoryInterface
	mockSecretRepo *mocks.MockSecretRepositoryInterface
}

// SetupTest creates a fresh gomock controller and mocks for each test.
func (suite *ServiceTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockUserRepo = mocks.NewMockUserRepositoryInterface(suite.ctrl)
	suite.mockSecretRepo = mocks.NewMockSecretRepositoryInterface(suite.ctrl)

	suite.userService = NewUserService(suite.mockUserRepo, testSalt)
	suite.secretService = NewSecretService(suite.mockSecretRepo)
}

// TearDownTest finishes gomock controller.
func (suite *ServiceTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}

// TestRegisterUser tests user registration.
func (suite *ServiceTestSuite) TestRegisterUser() {
	username := "testuser"
	password := "password123"

	suite.mockUserRepo.EXPECT().GetUserByUsername(username).Return(nil, errors.New("not found"))
	suite.mockUserRepo.EXPECT().CreateUser(gomock.Any(), username, gomock.Any()).Return(nil)

	user, err := suite.userService.RegisterUser(username, password)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), username, user.Username)
}

// TestAuthenticateUser tests user authentication.
func (suite *ServiceTestSuite) TestAuthenticateUser() {
	username := "authuser"
	password := "password123"
	hash := cryptography.HashPassword(password, testSalt)

	userModel := &models.User{ID: "u1", Username: username, PasswordHash: hash}
	suite.mockUserRepo.EXPECT().GetUserByUsername(username).Return(userModel, nil)

	user, err := suite.userService.AuthenticateUser(username, password)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), username, user.Username)
}

// TestCreateSecret tests creating a secret.
func (suite *ServiceTestSuite) TestCreateSecret() {
	userID := "user1"

	suite.mockSecretRepo.EXPECT().CreateSecret(gomock.AssignableToTypeOf(&models.Secret{})).DoAndReturn(
		func(s *models.Secret) error {
			assert.Equal(suite.T(), "My Password", s.Title)
			assert.Equal(suite.T(), userID, s.UserID)
			return nil
		},
	)

	secret, err := suite.secretService.CreateSecret(userID, models.SecretTypeCredentials,
		"My Password", "login:password", "website.com")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), secret)
	assert.Equal(suite.T(), "My Password", secret.Title)
}

// TestListSecrets tests listing user secrets.
func (suite *ServiceTestSuite) TestListSecrets() {
	userID := "listuser"

	expected := []models.Secret{
		{ID: "s1", UserID: userID, Title: "Secret 1"},
		{ID: "s2", UserID: userID, Title: "Secret 2"},
	}

	suite.mockSecretRepo.EXPECT().GetUserSecrets(userID).Return(expected, nil)

	secrets, err := suite.secretService.ListSecrets(userID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, len(secrets))
}

// TestGetSecret tests retrieving a specific secret.
func (suite *ServiceTestSuite) TestGetSecret() {
	userID := "uget"
	created := &models.Secret{ID: "sget", UserID: userID, Title: "My Secret"}

	suite.mockSecretRepo.EXPECT().GetSecretByID(created.ID).Return(created, nil)

	retrieved, err := suite.secretService.GetSecret(created.ID, userID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), created.ID, retrieved.ID)
	assert.Equal(suite.T(), "My Secret", retrieved.Title)
}

// TestUpdateSecret tests updating a secret.
func (suite *ServiceTestSuite) TestUpdateSecret() {
	userID := "upuser"
	created := &models.Secret{ID: "sup", UserID: userID, Title: "Original", Version: 1}

	suite.mockSecretRepo.EXPECT().GetSecretByID(created.ID).Return(created, nil)
	suite.mockSecretRepo.EXPECT().UpdateSecret(gomock.AssignableToTypeOf(&models.Secret{})).DoAndReturn(
		func(s *models.Secret) error {
			assert.Equal(suite.T(), "Updated", s.Title)
			assert.Equal(suite.T(), userID, s.UserID)
			return nil
		},
	)

	updated, err := suite.secretService.UpdateSecret(created.ID, userID,
		models.SecretTypeText, "Updated", "new_data", "new_meta")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated", updated.Title)
	assert.Equal(suite.T(), 1, updated.Version)
}

// TestDeleteSecret tests deleting a secret.
func (suite *ServiceTestSuite) TestDeleteSecret() {
	userID := "deluser"
	created := &models.Secret{ID: "sdel", UserID: userID, Title: "To Delete"}

	suite.mockSecretRepo.EXPECT().DeleteSecret(created.ID, userID).Return(nil)
	suite.mockSecretRepo.EXPECT().GetSecretByID(created.ID).Return(nil, errors.New("not found"))

	err := suite.secretService.DeleteSecret(created.ID, userID)
	assert.NoError(suite.T(), err)

	_, err = suite.secretService.GetSecret(created.ID, userID)
	assert.Error(suite.T(), err)
}

// TestServiceTestSuite runs the test suite.
func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
