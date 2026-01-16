package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Pklerik/gophKeep/internal/models"
)

// HTTPClient provides HTTP methods for server communication.
type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

// NewHTTPClient creates a new HTTP client.
func NewHTTPClient(baseURL string, timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// Register registers a new user.
func (c *HTTPClient) Register(username, password string) (*models.AuthResponse, error) {
	req := models.AuthRequest{
		Username: username,
		Password: password,
	}

	resp, err := c.post("/api/v1/auth/register", req)
	if err != nil {
		return nil, err
	}

	var authResp models.AuthResponse
	if err := json.Unmarshal(resp, &authResp); err != nil {
		return nil, fmt.Errorf("failed to parse registration response: %w", err)
	}

	c.token = authResp.Token
	return &authResp, nil
}

// Login logs in a user.
func (c *HTTPClient) Login(username, password string) (*models.AuthResponse, error) {
	req := models.AuthRequest{
		Username: username,
		Password: password,
	}

	resp, err := c.post("/api/v1/auth/login", req)
	if err != nil {
		return nil, err
	}

	var authResp models.AuthResponse
	if err := json.Unmarshal(resp, &authResp); err != nil {
		return nil, fmt.Errorf("failed to parse login response: %w", err)
	}

	c.token = authResp.Token
	return &authResp, nil
}

// CreateSecret creates a new secret.
func (c *HTTPClient) CreateSecret(secretType models.SecretType, title, data, metadata string) (*models.Secret, error) {
	req := map[string]interface{}{
		"type":     secretType,
		"title":    title,
		"data":     data,
		"metadata": metadata,
	}

	resp, err := c.post("/api/v1/secrets", req)
	if err != nil {
		return nil, err
	}

	var secret models.Secret
	if err := json.Unmarshal(resp, &secret); err != nil {
		return nil, fmt.Errorf("failed to parse secret response: %w", err)
	}

	return &secret, nil
}

// GetSecret retrieves a secret by ID.
func (c *HTTPClient) GetSecret(id string) (*models.Secret, error) {
	resp, err := c.getWithAuth(fmt.Sprintf("/api/v1/secrets/get?id=%s", id))
	if err != nil {
		return nil, err
	}

	var secret models.Secret
	if err := json.Unmarshal(resp, &secret); err != nil {
		return nil, fmt.Errorf("failed to parse secret response: %w", err)
	}

	return &secret, nil
}

// ListSecrets lists all secrets for the authenticated user.
func (c *HTTPClient) ListSecrets() ([]models.Secret, error) {
	resp, err := c.getWithAuth("/api/v1/secrets")
	if err != nil {
		return nil, err
	}

	var secrets []models.Secret
	if err := json.Unmarshal(resp, &secrets); err != nil {
		return nil, fmt.Errorf("failed to parse secrets list: %w", err)
	}

	return secrets, nil
}

// UpdateSecret updates an existing secret.
func (c *HTTPClient) UpdateSecret(id string, secretType models.SecretType, title, data, metadata string) (*models.Secret, error) {
	req := map[string]interface{}{
		"type":     secretType,
		"title":    title,
		"data":     data,
		"metadata": metadata,
	}

	resp, err := c.putWithAuth(fmt.Sprintf("/api/v1/secrets/update?id=%s", id), req)
	if err != nil {
		return nil, err
	}

	var secret models.Secret
	if err := json.Unmarshal(resp, &secret); err != nil {
		return nil, fmt.Errorf("failed to parse secret response: %w", err)
	}

	return &secret, nil
}

// DeleteSecret deletes a secret.
func (c *HTTPClient) DeleteSecret(id string) error {
	_, err := c.deleteWithAuth(fmt.Sprintf("/api/v1/secrets/delete?id=%s", id))
	return err
}

// SetToken sets the authentication token.
func (c *HTTPClient) SetToken(token string) {
	c.token = token
}

// post sends a POST request.
func (c *HTTPClient) post(endpoint string, data interface{}) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp models.ErrorResponse
		json.Unmarshal(respBody, &errResp)
		return nil, fmt.Errorf("server error: %s - %s", errResp.Code, errResp.Message)
	}

	return respBody, nil
}

// getWithAuth sends a GET request with authentication.
func (c *HTTPClient) getWithAuth(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp models.ErrorResponse
		json.Unmarshal(respBody, &errResp)
		return nil, fmt.Errorf("server error: %s - %s", errResp.Code, errResp.Message)
	}

	return respBody, nil
}

// putWithAuth sends a PUT request with authentication.
func (c *HTTPClient) putWithAuth(endpoint string, data interface{}) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("PUT", c.baseURL+endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp models.ErrorResponse
		json.Unmarshal(respBody, &errResp)
		return nil, fmt.Errorf("server error: %s - %s", errResp.Code, errResp.Message)
	}

	return respBody, nil
}

// deleteWithAuth sends a DELETE request with authentication.
func (c *HTTPClient) deleteWithAuth(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("DELETE", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errResp models.ErrorResponse
		json.Unmarshal(respBody, &errResp)
		return nil, fmt.Errorf("server error: %s - %s", errResp.Code, errResp.Message)
	}

	return respBody, nil
}
