// Package models contains the data structures and types used throughout the GophKeeper application.
// It defines the core models for users, secrets, and API requests/responses.
//
// Key Types:
// - User: Represents a registered user in the system
// - Secret: Represents encrypted secret data stored by a user
// - AuthRequest: User authentication request
// - AuthResponse: Authentication response with JWT token
//
// Secret Types:
// - SecretTypeCredentials: Login/password pairs
// - SecretTypeText: Plain text data
// - SecretTypeBinary: Binary data
// - SecretTypeCard: Credit card information
package models

import "time"

// User represents a registered user in the system.
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Secret represents encrypted secret data stored by a user.
type Secret struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Type      SecretType `json:"type"`
	Title     string     `json:"title"`
	Data      []byte     `json:"-"` // Encrypted data
	Metadata  string     `json:"metadata"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Version   int        `json:"version"`
}

// SecretType defines the type of secret data.
type SecretType string

const (
	SecretTypeCredentials SecretType = "credentials"
	SecretTypeText        SecretType = "text"
	SecretTypeBinary      SecretType = "binary"
	SecretTypeCard        SecretType = "card"
)

// Credentials represents login/password pair data.
type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// BankCard represents credit card data.
type BankCard struct {
	CardNumber string `json:"card_number"`
	Holder     string `json:"holder"`
	ExpiryDate string `json:"expiry_date"`
	CVV        string `json:"cvv"`
}

// AuthRequest represents user authentication request.
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse represents successful authentication response.
type AuthResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// SyncRequest represents data synchronization request from client.
type SyncRequest struct {
	ClientVersion int      `json:"client_version"`
	Secrets       []Secret `json:"secrets"`
}

// SyncResponse represents data synchronization response from server.
type SyncResponse struct {
	ServerVersion int      `json:"server_version"`
	Secrets       []Secret `json:"secrets"`
}

// ErrorResponse represents an API error response.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
