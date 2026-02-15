package harness

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"
)

// Discord mode constants.
const (
	DiscordModeLive    = "live"
	DiscordModeCapture = "capture"
)

// DiscordCapture is a local mock Discord webhook server that captures
// posted messages for test assertions.
type DiscordCapture struct {
	mu       sync.Mutex
	messages []string
	server   *httptest.Server
}

// NewDiscordCapture creates a mock Discord webhook server.
func NewDiscordCapture(t *testing.T) *DiscordCapture {
	t.Helper()
	d := &DiscordCapture{}
	d.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var payload struct {
			Content string `json:"content"`
		}
		_ = json.NewDecoder(r.Body).Decode(&payload)
		d.mu.Lock()
		d.messages = append(d.messages, payload.Content)
		d.mu.Unlock()
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(d.server.Close)
	return d
}

// URL returns the mock server's URL.
func (d *DiscordCapture) URL() string { return d.server.URL }

// Count returns how many messages have been captured.
func (d *DiscordCapture) Count() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.messages)
}

// Messages returns a copy of all captured messages.
func (d *DiscordCapture) Messages() []string {
	d.mu.Lock()
	defer d.mu.Unlock()
	out := make([]string, len(d.messages))
	copy(out, d.messages)
	return out
}

// Reset clears all captured messages.
func (d *DiscordCapture) Reset() {
	d.mu.Lock()
	d.messages = nil
	d.mu.Unlock()
}

// WaitForCount blocks until the capture has at least want messages.
func (d *DiscordCapture) WaitForCount(t *testing.T, want int, timeout time.Duration) []string {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if d.Count() >= want {
			return d.Messages()
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("discord webhook did not reach %d events in %s (got %d)", want, timeout, d.Count())
	return nil
}

// SetupDiscord configures Discord webhook behavior on the harness.
func SetupDiscord(t *testing.T, h *Harness) {
	t.Helper()
	if os.Getenv("E2E_DISCORD_LIVE") == "1" {
		h.DiscordMode = DiscordModeLive
		signup := os.Getenv("DISCORD_WEBHOOK_URL_SIGNUP")
		billing := os.Getenv("DISCORD_WEBHOOK_URL_BILLING")
		if signup == "" || billing == "" {
			t.Fatal("E2E_DISCORD_LIVE=1 requires DISCORD_WEBHOOK_URL_SIGNUP and DISCORD_WEBHOOK_URL_BILLING")
		}
		return
	}

	h.DiscordMode = DiscordModeCapture
	h.SignupWebhook = NewDiscordCapture(t)
	h.BillingWebhook = NewDiscordCapture(t)
	t.Setenv("DISCORD_WEBHOOK_URL_SIGNUP", h.SignupWebhook.URL())
	t.Setenv("DISCORD_WEBHOOK_URL_BILLING", h.BillingWebhook.URL())

	h.SnapshotWebhookURL = os.Getenv("DISCORD_WEBHOOK_URL_SNAPSHOTS")
}

// CaptureForKind returns the DiscordCapture for the given kind ("signup" or "billing").
func CaptureForKind(h *Harness, kind string) *DiscordCapture {
	switch kind {
	case "signup":
		return h.SignupWebhook
	case "billing":
		return h.BillingWebhook
	default:
		return nil
	}
}
