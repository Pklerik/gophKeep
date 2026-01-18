package httpclient

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Pklerik/gophKeep/internal/models"

	"github.com/go-resty/resty/v2"
)

// HTTPClient provides HTTP methods for server communication.
type HTTPClient struct {
	httpClient *resty.Client
}

// NewHTTPClient creates a new HTTP client.
func NewHTTPClient(baseURL string, timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		httpClient: resty.New().SetTimeout(timeout).SetBaseURL(baseURL),
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

	c.httpClient.SetAuthToken(authResp.Token)
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

	c.httpClient.SetAuthToken(authResp.Token)
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
	resp, err := c.get(fmt.Sprintf("/api/v1/secrets/get?id=%s", id))
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
	resp, err := c.get("/api/v1/secrets")
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

// post sends a POST request.
func (c *HTTPClient) post(endpoint string, data interface{}) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.R().SetBody(body).Post(endpoint)
	if err != nil {
		return nil, err
	}

	respBody := resp.Body()

	if resp.StatusCode() >= 400 {
		var errResp models.ErrorResponse
		json.Unmarshal(respBody, &errResp)
		return nil, fmt.Errorf("server error: %s - %s", errResp.Code, errResp.Message)
	}

	return respBody, nil
}

// post sends a GET request.
func (c *HTTPClient) get(endpoint string) ([]byte, error) {
	resp, err := c.httpClient.R().Get(endpoint)
	if err != nil {
		return nil, err
	}

	respBody := resp.Body()

	if resp.StatusCode() >= 400 {
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

	resp, err := c.httpClient.R().SetHeader("Content-Type", "application/json").SetBody(body).Put(endpoint)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	respBody := resp.Body()

	if resp.StatusCode() >= 400 {
		var errResp models.ErrorResponse
		json.Unmarshal(respBody, &errResp)
		return nil, fmt.Errorf("server error: %s - %s", errResp.Code, errResp.Message)
	}

	return respBody, nil
}

// deleteWithAuth sends a DELETE request with authentication.
func (c *HTTPClient) deleteWithAuth(endpoint string) ([]byte, error) {
	resp, err := c.httpClient.R().Delete(endpoint)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	respBody := resp.Body()
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode() >= 400 {
		var errResp models.ErrorResponse
		json.Unmarshal(respBody, &errResp)
		return nil, fmt.Errorf("server error: %s - %s", errResp.Code, errResp.Message)
	}

	return respBody, nil
}
