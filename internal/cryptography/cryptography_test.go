// Package cryptography provides encryption, decryption, and hashing utilities.
package cryptography

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEncryptDecrypt tests encryption and decryption.
func TestEncryptDecrypt(t *testing.T) {
	plaintext := "secret message"
	password := "mypassword"
	salt := "somesalt"

	encrypted, err := EncryptAES(plaintext, password, salt)
	assert.NoError(t, err)
	assert.NotEqual(t, plaintext, encrypted)

	decrypted, err := DecryptAES(encrypted, password, salt)
	assert.NoError(t, err)
	assert.Equal(t, plaintext, decrypted)
}

// TestDecryptWrongPassword tests decryption with wrong password.
func TestDecryptWrongPassword(t *testing.T) {
	plaintext := "secret"
	password := "correct"
	salt := "somesalt"

	encrypted, _ := EncryptAES(plaintext, password, salt)
	_, err := DecryptAES(encrypted, "wrong", salt)
	assert.Error(t, err)
}

// TestHashPassword tests password hashing.
func TestHashPassword(t *testing.T) {
	password := "mypassword"
	salt := []byte("somesalt")
	hash := HashPassword(password, salt)

	assert.NotEqual(t, password, hash)
	assert.True(t, VerifyPassword(password, hash, salt))
	assert.False(t, VerifyPassword("wrongpass", hash, salt))
}
