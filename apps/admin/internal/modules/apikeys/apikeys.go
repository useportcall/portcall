package apikeys

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/useportcall/portcall/libs/go/routerx"
)

// GeneratedAPIKey represents a generated API key
type GeneratedAPIKey struct {
	Key       string    `json:"key"`
	CreatedAt time.Time `json:"created_at"`
	Usage     string    `json:"usage"`
}

// GenerateAPIKeyRequest specifies what kind of key to generate
type GenerateAPIKeyRequest struct {
	Usage string `json:"usage"` // "dashboard", "admin", etc.
}

// GenerateAPIKeyResponse returns the generated key
type GenerateAPIKeyResponse struct {
	Key     GeneratedAPIKey `json:"key"`
	Message string          `json:"message"`
	Notes   []string        `json:"notes"`
}

// generateSecureKey creates a cryptographically secure API key
func generateSecureKey(prefix string, length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return prefix + hex.EncodeToString(bytes), nil
}

// GenerateAdminAPIKey generates a new admin API key for use with the dashboard
// This key can be used to authenticate dashboard-to-admin API calls
func GenerateAdminAPIKey(c *routerx.Context) {
	var body GenerateAPIKeyRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		body.Usage = "dashboard" // Default to dashboard usage
	}

	// Generate a secure key
	key, err := generateSecureKey("admin_", 32)
	if err != nil {
		c.ServerError("Failed to generate key", err)
		return
	}

	response := GenerateAPIKeyResponse{
		Key: GeneratedAPIKey{
			Key:       key,
			CreatedAt: time.Now().UTC(),
			Usage:     body.Usage,
		},
		Message: "API key generated successfully",
		Notes: []string{
			"Store this key securely - it won't be shown again",
			"Add to ADMIN_API_KEY environment variable for admin service",
			"For dashboard usage, set ADMIN_API_KEY in dashboard service",
			"Use with X-Admin-API-Key header in requests",
		},
	}

	log.Printf("Generated new admin API key for usage: %s", body.Usage)

	c.OK(response)
}

// ValidateAPIKey checks if the current API key is valid
func ValidateAPIKey(c *routerx.Context) {
	adminAPIKey := os.Getenv("ADMIN_API_KEY")
	providedKey := c.GetHeader("X-Admin-API-Key")

	if adminAPIKey == "" {
		c.OK(map[string]any{
			"valid":   false,
			"reason":  "ADMIN_API_KEY not configured on server",
			"message": "No admin API key is configured on this server",
		})
		return
	}

	if providedKey == "" {
		c.OK(map[string]any{
			"valid":   false,
			"reason":  "no key provided",
			"message": "No X-Admin-API-Key header provided",
		})
		return
	}

	if providedKey != adminAPIKey {
		c.OK(map[string]any{
			"valid":   false,
			"reason":  "key mismatch",
			"message": "The provided API key does not match",
		})
		return
	}

	c.OK(map[string]any{
		"valid":   true,
		"message": "API key is valid",
	})
}

// GetCurrentKeyInfo returns info about the currently configured key (masked)
func GetCurrentKeyInfo(c *routerx.Context) {
	adminAPIKey := os.Getenv("ADMIN_API_KEY")

	if adminAPIKey == "" {
		c.OK(map[string]any{
			"configured": false,
			"message":    "No admin API key is configured",
		})
		return
	}

	// Mask the key, show only first 10 and last 4 characters
	masked := adminAPIKey
	if len(adminAPIKey) > 14 {
		masked = fmt.Sprintf("%s...%s", adminAPIKey[:10], adminAPIKey[len(adminAPIKey)-4:])
	}

	c.OK(map[string]any{
		"configured": true,
		"masked_key": masked,
		"prefix":     "admin_",
	})
}
