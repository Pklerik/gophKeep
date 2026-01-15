// Package handler provides HTTP request handlers for the GophKeeper server.
package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Pklerik/gophKeep/internal/auth"
	"github.com/Pklerik/gophKeep/internal/models"
	"github.com/Pklerik/gophKeep/internal/service"
)

// Handler contains all the handlers and their dependencies.
type Handler struct {
	userService   *service.UserService
	secretService *service.SecretService
}

// NewHandler creates a new handler.
func NewHandler(userService *service.UserService, secretService *service.SecretService) *Handler {
	return &Handler{
		userService:   userService,
		secretService: secretService,
	}
}

// Register handles user registration.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only POST is allowed")
		return
	}

	var req models.AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	user, err := h.userService.RegisterUser(req.Username, req.Password)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "REGISTRATION_FAILED", err.Error())
		return
	}

	token, expiresAt, err := auth.GenerateToken(user.ID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "TOKEN_GENERATION_FAILED", "Failed to generate token")
		return
	}

	resp := models.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}

	writeJSONResponse(w, http.StatusCreated, resp)
}

// Login handles user login.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only POST is allowed")
		return
	}

	var req models.AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	user, err := h.userService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, "AUTHENTICATION_FAILED", err.Error())
		return
	}

	token, expiresAt, err := auth.GenerateToken(user.ID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "TOKEN_GENERATION_FAILED", "Failed to generate token")
		return
	}

	resp := models.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}

	writeJSONResponse(w, http.StatusOK, resp)
}

// CreateSecret handles secret creation.
func (h *Handler) CreateSecret(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only POST is allowed")
		return
	}

	userID, err := h.extractUserID(r)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to read request body")
		return
	}

	var req struct {
		Type     models.SecretType `json:"type"`
		Title    string            `json:"title"`
		Data     string            `json:"data"`
		Metadata string            `json:"metadata"`
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	secret, err := h.secretService.CreateSecret(userID, req.Type, req.Title, req.Data, req.Metadata)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "CREATION_FAILED", err.Error())
		return
	}

	writeJSONResponse(w, http.StatusCreated, secret)
}

// GetSecret retrieves a secret by ID.
func (h *Handler) GetSecret(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only GET is allowed")
		return
	}

	userID, err := h.extractUserID(r)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	secretID := r.URL.Query().Get("id")
	if secretID == "" {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Secret ID is required")
		return
	}

	secret, err := h.secretService.GetSecret(secretID, userID)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, "NOT_FOUND", err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, secret)
}

// ListSecrets retrieves all secrets for the authenticated user.
func (h *Handler) ListSecrets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only GET is allowed")
		return
	}

	userID, err := h.extractUserID(r)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	secrets, err := h.secretService.ListSecrets(userID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, secrets)
}

// UpdateSecret updates an existing secret.
func (h *Handler) UpdateSecret(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only PUT is allowed")
		return
	}

	userID, err := h.extractUserID(r)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	secretID := r.URL.Query().Get("id")
	if secretID == "" {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Secret ID is required")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to read request body")
		return
	}

	var req struct {
		Type     models.SecretType `json:"type"`
		Title    string            `json:"title"`
		Data     string            `json:"data"`
		Metadata string            `json:"metadata"`
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	secret, err := h.secretService.UpdateSecret(secretID, userID, req.Type, req.Title, req.Data, req.Metadata)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "UPDATE_FAILED", err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, secret)
}

// DeleteSecret deletes a secret.
func (h *Handler) DeleteSecret(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeErrorResponse(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only DELETE is allowed")
		return
	}

	userID, err := h.extractUserID(r)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
		return
	}

	secretID := r.URL.Query().Get("id")
	if secretID == "" {
		writeErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Secret ID is required")
		return
	}

	err = h.secretService.DeleteSecret(secretID, userID)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, "DELETE_FAILED", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// extractUserID extracts user ID from the JWT token in the request.
func (h *Handler) extractUserID(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", error(nil) // return error as error interface
	}

	tokenString, err := auth.ExtractToken(authHeader)
	if err != nil {
		return "", err
	}

	claims, err := auth.VerifyToken(tokenString)
	if err != nil {
		return "", err
	}

	return claims.UserID, nil
}

// writeJSONResponse writes a JSON response to the client.
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// writeErrorResponse writes an error response to the client.
func writeErrorResponse(w http.ResponseWriter, statusCode int, code, message string) {
	resp := models.ErrorResponse{
		Code:    code,
		Message: message,
	}
	writeJSONResponse(w, statusCode, resp)
}
