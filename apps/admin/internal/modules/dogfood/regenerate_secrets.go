package dogfood

import (
	"fmt"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type RegenerateSecretsRequest struct {
	UpdateK8s bool `json:"update_k8s"`
}

type RegenerateSecretsResponse struct {
	LiveSecret string `json:"live_secret"`
	TestSecret string `json:"test_secret"`
	K8sUpdated bool   `json:"k8s_updated,omitempty"`
	K8sMessage string `json:"k8s_message,omitempty"`
}

// RegenerateSecrets disables existing secrets and creates new ones
// This is useful when secrets need to be added to K8s
func RegenerateSecrets(c *routerx.Context) {
	var body RegenerateSecretsRequest
	_ = c.ShouldBindJSON(&body)

	// Find dogfood account
	var account models.Account
	if err := c.DB().FindFirst(&account, "email = ?", DogfoodAccountEmail); err != nil {
		c.ServerError("Dogfood account not found", err)
		return
	}

	// Find live app
	var liveApp models.App
	if err := c.DB().FindFirst(&liveApp, "account_id = ? AND name = ?", account.ID, DogfoodLiveAppName); err != nil {
		c.ServerError("Dogfood live app not found", err)
		return
	}

	// Find test app
	var testApp models.App
	if err := c.DB().FindFirst(&testApp, "account_id = ? AND name = ?", account.ID, DogfoodTestAppName); err != nil {
		c.ServerError("Dogfood test app not found", err)
		return
	}

	// Disable existing secrets and create new ones
	liveSecret, err := regenerateSecretForApp(c, liveApp.ID)
	if err != nil {
		c.ServerError("Failed to regenerate live app secret", err)
		return
	}

	testSecret, err := regenerateSecretForApp(c, testApp.ID)
	if err != nil {
		c.ServerError("Failed to regenerate test app secret", err)
		return
	}

	response := RegenerateSecretsResponse{
		LiveSecret: liveSecret,
		TestSecret: testSecret,
	}

	// Update K8s secrets if requested
	if body.UpdateK8s {
		k8sUpdated, k8sMsg := UpdateK8sSecretsWithNewKeys(liveSecret, testSecret)
		response.K8sUpdated = k8sUpdated
		response.K8sMessage = k8sMsg
	}

	c.OK(response)
}

func regenerateSecretForApp(c *routerx.Context, appID uint) (string, error) {
	// Disable existing secrets
	now := time.Now()
	var existingSecrets []models.Secret
	if err := c.DB().List(&existingSecrets, "app_id = ? AND disabled_at IS NULL AND key_type = ?", appID, "api_key"); err != nil {
		if !dbx.IsRecordNotFoundError(err) {
			return "", err
		}
	}

	for _, secret := range existingSecrets {
		secret.DisabledAt = &now
		if err := c.DB().Save(&secret); err != nil {
			return "", fmt.Errorf("failed to disable old secret: %w", err)
		}
	}

	// Generate new API key
	apiKey, err := generateAPIKey(64)
	if err != nil {
		return "", err
	}

	publicID := dbx.GenPublicID("sk")

	hash, err := c.Crypto().Encrypt(apiKey)
	if err != nil {
		return "", err
	}

	secret := models.Secret{
		PublicID: publicID,
		AppID:    appID,
		KeyHash:  hash,
		KeyType:  "api_key",
	}

	if err := c.DB().Create(&secret); err != nil {
		return "", err
	}

	// Return full key (only time it's visible)
	return fmt.Sprintf("%s_%s", publicID, apiKey), nil
}
