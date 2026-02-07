// Package handler provides HTTP request handlers for the GophKeeper server.
package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Pklerik/gophKeep/internal/auth"
	"github.com/Pklerik/gophKeep/internal/handler/server/dictionary"
	"github.com/Pklerik/gophKeep/internal/models"
	"github.com/Pklerik/gophKeep/internal/service"
	"github.com/go-chi/chi/v5"
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
		writeErrorResponse(w, http.StatusMethodNotAllowed, dictionary.MethodNotAllowed, fmt.Sprintf(dictionary.OnlyAllowed, http.MethodPost))
		return
	}

	var req models.AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, dictionary.InvalidRequest, "Invalid request body")
		return
	}

	user, err := h.userService.RegisterUser(req.Username, req.Password)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, dictionary.RegistrationFailed, err.Error())
		return
	}

	token, expiresAt, err := auth.GenerateToken(user.ID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, dictionary.TokenGenerationFailed, "Failed to generate token")
		return
	}

	resp := models.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}
	w.Header().Set("Authorization", "Bearer "+token)

	writeJSONResponse(w, http.StatusCreated, resp)
}

// Login handles user login.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, dictionary.MethodNotAllowed, fmt.Sprintf(dictionary.OnlyAllowed, http.MethodPost))
		return
	}

	var req models.AuthRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, dictionary.InvalidRequest, "Invalid request body")
		return
	}

	user, err := h.userService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, dictionary.AuthenticationFailed, err.Error())
		return
	}

	token, expiresAt, err := auth.GenerateToken(user.ID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, dictionary.TokenGenerationFailed, "Failed to generate token")
		return
	}

	resp := models.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}

	w.Header().Set("Authorization", "Bearer "+token)

	writeJSONResponse(w, http.StatusOK, resp)
}

// CreateSecret handles secret creation.
func (h *Handler) CreateSecret(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(w, http.StatusMethodNotAllowed, dictionary.MethodNotAllowed, fmt.Sprintf(dictionary.OnlyAllowed, http.MethodPost))
		return
	}

	userID, err := h.extractUserID(r)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, dictionary.Unauthorized, err.Error())
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, dictionary.InvalidRequest, "Failed to read request body")
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
		writeErrorResponse(w, http.StatusBadRequest, dictionary.InvalidRequest, "Invalid request body")
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
		writeErrorResponse(w, http.StatusMethodNotAllowed, dictionary.MethodNotAllowed, fmt.Sprintf(dictionary.OnlyAllowed, http.MethodGet))
		return
	}

	userID, err := h.extractUserID(r)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, dictionary.Unauthorized, err.Error())
		return
	}

	secretID := chi.URLParam(r, "id")
	if secretID == "" {
		writeErrorResponse(w, http.StatusBadRequest, dictionary.InvalidRequest, "Secret ID is required")
		return
	}

	secret, err := h.secretService.GetSecret(secretID, userID)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, dictionary.NotFound, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, secret)
}

// ListSecrets retrieves all secrets for the authenticated user.
func (h *Handler) ListSecrets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, http.StatusMethodNotAllowed, dictionary.MethodNotAllowed, fmt.Sprintf(dictionary.OnlyAllowed, http.MethodGet))
		return
	}

	userID, err := h.extractUserID(r)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, dictionary.Unauthorized, err.Error())
		return
	}

	secrets, err := h.secretService.ListSecrets(userID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, dictionary.ListFailed, err.Error())
		return
	}

	writeJSONResponse(w, http.StatusOK, secrets)
}

// UpdateSecret updates an existing secret.
func (h *Handler) UpdateSecret(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeErrorResponse(w, http.StatusMethodNotAllowed, dictionary.MethodNotAllowed, fmt.Sprintf(dictionary.OnlyAllowed, http.MethodPut))
		return
	}

	userID, err := h.extractUserID(r)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, dictionary.Unauthorized, err.Error())
		return
	}

	secretID := chi.URLParam(r, "id")
	if secretID == "" {
		writeErrorResponse(w, http.StatusBadRequest, dictionary.InvalidRequest, "Secret ID is required")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, dictionary.InvalidRequest, "Failed to read request body")
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
		writeErrorResponse(w, http.StatusBadRequest, dictionary.InvalidRequest, "Invalid request body")
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
		writeErrorResponse(w, http.StatusMethodNotAllowed, dictionary.MethodNotAllowed, fmt.Sprintf(dictionary.OnlyAllowed, http.MethodDelete))
		return
	}

	userID, err := h.extractUserID(r)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, dictionary.Unauthorized, err.Error())
		return
	}

	secretID := chi.URLParam(r, "id")
	if secretID == "" {
		writeErrorResponse(w, http.StatusBadRequest, dictionary.InvalidRequest, "Secret ID is required")
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
