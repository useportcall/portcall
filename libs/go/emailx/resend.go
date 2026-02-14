package emailx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type resendEmailClient struct {
	apiKey     string
	apiURL     string
	httpClient *http.Client
}

func resendAPIURL() string {
	url := os.Getenv("RESEND_API_URL")
	if url == "" {
		return "https://api.resend.com/emails"
	}
	return url
}

func (c *resendEmailClient) resolveAPIURL() string {
	if c.apiURL != "" {
		return c.apiURL
	}
	return resendAPIURL()
}

func (c *resendEmailClient) resolveClient() *http.Client {
	if c.httpClient != nil {
		return c.httpClient
	}
	return &http.Client{Timeout: 10 * time.Second}
}

// ResendRequest represents the structure expected by Resend's Send API
type resendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Html    string   `json:"html,omitempty"`
	Text    string   `json:"text,omitempty"`
}

func (c *resendEmailClient) Send(content, subject, from string, to []string) error {
	// Prepare the request payload
	payload := resendRequest{
		From:    from,
		To:      to,
		Subject: subject,
		Html:    content,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal email payload: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", c.resolveAPIURL(), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.resolveClient().Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to Resend: %w", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		// Try to read error message from body
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return fmt.Errorf("resend API error: status %d, body: %s", resp.StatusCode, buf.String())
	}

	return nil
}
