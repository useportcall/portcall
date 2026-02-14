package discordx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Client struct {
	WebhookURL string
}

func New(webhookURL string) *Client {
	return &Client{
		WebhookURL: webhookURL,
	}
}

// Send sends a message to the configured Discord webhook.
// If the webhook URL is not set (e.g. in local dev without env var), it logs the message instead.
func (c *Client) Send(content string) error {
	if c.WebhookURL == "" {
		fmt.Printf("[Discord Stub] %s\n", content)
		return nil
	}

	payload := map[string]string{
		"content": content,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal discord payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.WebhookURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create discord request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send discord request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord webhook failed with status: %d", resp.StatusCode)
	}

	return nil
}

// DefaultClient returns a client using the DISCORD_WEBHOOK_URL environment variable.
// Useful for simple cases where a single webhook is enough.
func DefaultClient() *Client {
	return New(os.Getenv("DISCORD_WEBHOOK_URL"))
}
