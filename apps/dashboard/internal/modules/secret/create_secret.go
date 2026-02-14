package secret

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func CreateSecret(c *routerx.Context) {
	var app models.App
	if err := c.DB().FindForID(c.AppID(), &app); err != nil {
		c.ServerError("Internal server error", err)
		return
	}

	// create random api key
	apiKey, err := GenerateAPIKey(64)
	if err != nil {
		c.ServerError("Internal server error", err)
		return
	}

	publicID := dbx.GenPublicID("sk")

	hash, err := c.Crypto().Encrypt(apiKey)
	if err != nil {
		c.ServerError("Failed to encrypt secret value", err)
		return
	}

	secret := new(models.Secret)
	secret.PublicID = publicID
	secret.AppID = c.AppID()
	secret.KeyHash = hash
	secret.KeyType = "api_key"

	if err := c.DB().Create(secret); err != nil {
		c.ServerError("Failed to create secret", err)
		return
	}

	key := fmt.Sprintf("%s_%s", publicID, apiKey)

	c.OK(map[string]any{"key": key, "public_id": secret.PublicID})
}

type CreateSecretRequest struct {
	Value string `json:"value"`
}

// GenerateAPIKey creates a cryptographically secure API key
// Returns a hex-encoded string of specified byte length
func GenerateAPIKey(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be positive")
	}

	// Calculate bytes needed (2 hex chars = 1 byte)
	byteLength := (length + 1) / 2
	if length%2 != 0 {
		byteLength++
	}

	bytes := make([]byte, byteLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Encode to hex and trim to exact length
	hexStr := hex.EncodeToString(bytes)
	if len(hexStr) > length {
		hexStr = hexStr[:length]
	}

	return hexStr, nil
}
