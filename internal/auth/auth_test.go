// Package auth provides authentication and JWT token handling.
package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGenerateToken tests token generation.
func TestGenerateToken(t *testing.T) {
	userID := "test-user"
	token, expiresAt, err := GenerateToken(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.False(t, expiresAt.IsZero())
}

// TestVerifyToken tests token verification.
func TestVerifyToken(t *testing.T) {
	userID := "test-user"
	token, _, _ := GenerateToken(userID)

	claims, err := VerifyToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
}

// TestVerifyInvalidToken tests invalid token verification.
func TestVerifyInvalidToken(t *testing.T) {
	_, err := VerifyToken("invalid")
	assert.Error(t, err)
}

// TestExtractToken tests token extraction from header.
func TestExtractToken(t *testing.T) {
	token, err := ExtractToken("Bearer valid-token")
	assert.NoError(t, err)
	assert.Equal(t, "valid-token", token)
}

// TestExtractTokenInvalid tests invalid header format.
func TestExtractTokenInvalid(t *testing.T) {
	_, err := ExtractToken("InvalidFormat")
	assert.Error(t, err)
}
