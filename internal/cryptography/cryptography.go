// Package cryptography provides encryption, decryption, and hashing utilities
// for secure data handling in the GophKeeper application.
package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

// EncryptAES encrypts data using AES-256-GCM.
// The returned value is base64 encoded and includes the nonce.
func EncryptAES(plaintext, password, salt string) (string, error) {
	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES decrypts data encrypted with EncryptAES.
// The input should be base64 encoded.
func DecryptAES(ciphertextB64, password, salt string) (string, error) {
	key := deriveKey(password, salt)

	data, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesgcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// HashPassword creates a password hash using PBKDF2.
func HashPassword(password string, salt []byte) string {
	hash := pbkdf2.Key([]byte(password), salt, 4096, 32, sha256.New)
	return base64.StdEncoding.EncodeToString(hash)
}

// VerifyPassword checks if the provided password matches the hash.
func VerifyPassword(password, hash string, salt []byte) bool {
	return HashPassword(password, salt) == hash
}

// deriveKey derives a 32-byte encryption key from a password using PBKDF2.
func deriveKey(password, salt string) []byte {
	return pbkdf2.Key([]byte(password), []byte(salt), 4096, 32, sha256.New)
}
