package resources

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient is the underlying HTTP client for the SDK
type HTTPClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(apiKey, baseURL string) *HTTPClient {
	return &HTTPClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetBaseURL changes the base URL
func (h *HTTPClient) SetBaseURL(url string) {
	h.baseURL = url
}

// APIError represents an error from the API
type APIError struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	Status  int    `json:"status,omitempty"`
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return e.Message
}

// Get performs a GET request
func (h *HTTPClient) Get(ctx context.Context, path string, result interface{}) error {
	return h.doRequest(ctx, http.MethodGet, path, nil, result)
}

// Post performs a POST request
func (h *HTTPClient) Post(ctx context.Context, path string, body interface{}, result interface{}) error {
	return h.doRequest(ctx, http.MethodPost, path, body, result)
}

// Put performs a PUT request
func (h *HTTPClient) Put(ctx context.Context, path string, body interface{}, result interface{}) error {
	return h.doRequest(ctx, http.MethodPut, path, body, result)
}

// Delete performs a DELETE request
func (h *HTTPClient) Delete(ctx context.Context, path string, result interface{}) error {
	return h.doRequest(ctx, http.MethodDelete, path, nil, result)
}

// Patch performs a PATCH request
func (h *HTTPClient) Patch(ctx context.Context, path string, body interface{}, result interface{}) error {
	return h.doRequest(ctx, http.MethodPatch, path, body, result)
}

func (h *HTTPClient) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	url := h.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", h.apiKey)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			return &APIError{
				Message: string(respBody),
				Status:  resp.StatusCode,
			}
		}
		apiErr.Status = resp.StatusCode
		return &apiErr
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// DataWrapper is a generic wrapper for API responses that have a "data" field
type DataWrapper[T any] struct {
	Data T `json:"data"`
}
