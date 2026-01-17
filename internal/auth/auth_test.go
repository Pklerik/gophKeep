package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Ensure secret key is set for the package tests once.
	SetSecretKey("s3cr3t-tests")
}

func TestGenerateAndVerifyToken(t *testing.T) {
	tok, exp, err := GenerateToken("user-1")
	assert.NoError(t, err)
	assert.NotEmpty(t, tok)
	assert.True(t, exp.After(time.Now()))

	claims, err := VerifyToken(tok)
	assert.NoError(t, err)
	assert.Equal(t, "user-1", claims.UserID)
}

func TestVerifyToken_InvalidSignature(t *testing.T) {
	// Create token signed with a different secret
	claims := &Claims{
		UserID: "user-2",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "gophkeeper",
		},
	}

	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tkn.SignedString([]byte("other-secret"))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	_, err = VerifyToken(tokenString)
	assert.Error(t, err)
}

func TestVerifyInvalidToken(t *testing.T) {
	_, err := VerifyToken("invalid-token-string")
	assert.Error(t, err)
}

func TestVerifyToken_UnexpectedSigningMethod(t *testing.T) {
	// create RSA key to sign token with RS256 (non-HMAC)
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate rsa key: %v", err)
	}

	claims := &Claims{
		UserID: "user-rs",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "gophkeeper",
		},
	}

	tkn := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := tkn.SignedString(priv)
	if err != nil {
		t.Fatalf("failed to sign rs token: %v", err)
	}

	_, err = VerifyToken(tokenString)
	assert.Error(t, err)
}

func TestExtractToken(t *testing.T) {
	v, err := ExtractToken("Bearer ABCDEF")
	assert.NoError(t, err)
	assert.Equal(t, "ABCDEF", v)

	_, err = ExtractToken("Bear ABC")
	assert.Error(t, err)

	_, err = ExtractToken("")
	assert.Error(t, err)
}
