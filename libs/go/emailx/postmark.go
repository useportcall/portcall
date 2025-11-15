package emailx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type postmarkEmailClient struct {
	apiKey string
}

// EmailRequest represents the structure expected by Postmark's Send API
type postmarkRequest struct {
	From     string `json:"From"`
	To       string `json:"To"`
	Subject  string `json:"Subject"`
	HtmlBody string `json:"HtmlBody,omitempty"`
	TextBody string `json:"TextBody,omitempty"`
}

func (c *postmarkEmailClient) Send(content, subject, from string, to []string) error {
	// Postmark expects a comma-separated list of recipients
	toList := ""
	for i, recipient := range to {
		if i > 0 {
			toList += ","
		}
		toList += recipient
	}

	// Prepare the request payload
	payload := postmarkRequest{
		From:     from,
		To:       toList,
		Subject:  subject,
		HtmlBody: content, // You can split into HtmlBody/TextBody if needed
		// TextBody: content, // Uncomment and adjust if you want plain text fallback
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal email payload: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.postmarkapp.com/email", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Postmark-Server-Token", c.apiKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to Postmark: %w", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK {
		// Try to read error message from body
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return fmt.Errorf("postmark API error: status %d, body: %s", resp.StatusCode, buf.String())
	}

	return nil
}
