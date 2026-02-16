package dogfood

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

const (
	// K8s configuration
	DefaultNamespace = "portcall"
	SecretName       = "portcall-secrets"
)

type K8sUpdateRequest struct {
	UpdateK8s bool `json:"update_k8s"`
}

type K8sUpdateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// RefreshK8sSecrets updates the Kubernetes secrets with current dogfood configuration
func RefreshK8sSecrets(c *routerx.Context) {
	// Get dogfood status
	var account models.Account
	err := c.DB().FindFirst(&account, "email = ?", DogfoodAccountEmail)
	if err != nil {
		c.NotFound("Dogfood account not found")
		return
	}

	// Get live app
	var liveApp models.App
	err = c.DB().FindFirst(&liveApp, "account_id = ? AND is_live = ?", account.ID, true)
	if err != nil {
		c.NotFound("Dogfood live app not found")
		return
	}

	// Get secrets
	var secrets []models.Secret
	if err := c.DB().List(&secrets, "app_id = ? AND disabled_at IS NULL", liveApp.ID); err != nil {
		c.ServerError("Failed to list secrets", err)
		return
	}

	if len(secrets) == 0 {
		c.BadRequest("No active secrets found for dogfood app")
		return
	}

	// We can't retrieve the original secret keys (they're hashed)
	// So we just confirm the configuration is correct and update the K8s secret
	// with a placeholder message indicating secrets need to be manually set

	log.Printf("Refreshing K8s secrets for dogfood app %s", liveApp.PublicID)

	// Check if kubectl is available
	if err := checkKubectlAccess(); err != nil {
		c.OK(K8sUpdateResponse{
			Success: false,
			Message: fmt.Sprintf("kubectl not available: %v. Manual update required.", err),
		})
		return
	}

	c.OK(K8sUpdateResponse{
		Success: true,
		Message: "K8s access verified. Secrets exist in database but cannot be retrieved (hashed). Use setup endpoint to generate new secrets if needed.",
	})
}

// UpdateK8sSecretsWithNewKeys updates Kubernetes secrets with newly generated API keys
// This is called during dogfood setup when new secrets are created
func UpdateK8sSecretsWithNewKeys(liveSecret, testSecret string) (bool, string) {
	if liveSecret == "" || testSecret == "" {
		return false, "Missing secrets"
	}

	// Check if kubectl is available
	if err := checkKubectlAccess(); err != nil {
		return false, fmt.Sprintf("kubectl not available: %v", err)
	}

	// Get current secret
	cmd := exec.Command("kubectl", "get", "secret", SecretName, "-n", DefaultNamespace, "-o", "json")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return false, fmt.Sprintf("Failed to get secret: %v - %s", err, stderr.String())
	}

	// Parse the secret
	var secret map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &secret); err != nil {
		return false, fmt.Sprintf("Failed to parse secret: %v", err)
	}

	// Update the data field
	data, ok := secret["data"].(map[string]interface{})
	if !ok {
		data = make(map[string]interface{})
		secret["data"] = data
	}

	// Encode and set the new values
	data["DOGFOOD_LIVE_SECRET"] = base64.StdEncoding.EncodeToString([]byte(liveSecret))
	data["DOGFOOD_TEST_SECRET"] = base64.StdEncoding.EncodeToString([]byte(testSecret))

	// Apply the updated secret
	updatedJSON, err := json.Marshal(secret)
	if err != nil {
		return false, fmt.Sprintf("Failed to marshal updated secret: %v", err)
	}

	cmd = exec.Command("kubectl", "apply", "-f", "-")
	cmd.Stdin = bytes.NewReader(updatedJSON)
	var applyStderr bytes.Buffer
	cmd.Stderr = &applyStderr

	if err := cmd.Run(); err != nil {
		return false, fmt.Sprintf("Failed to apply secret: %v - %s", err, applyStderr.String())
	}

	log.Printf("Successfully updated K8s secret %s in namespace %s", SecretName, DefaultNamespace)

	// Restart dashboard deployment to pick up new secrets
	cmd = exec.Command("kubectl", "rollout", "restart", "deployment/dashboard", "-n", DefaultNamespace)
	if err := cmd.Run(); err != nil {
		log.Printf("Warning: Failed to restart dashboard: %v", err)
		return true, "K8s secrets updated. Failed to auto-restart dashboard - please restart manually."
	}

	return true, "K8s secrets updated and dashboard restarted"
}

func checkKubectlAccess() error {
	// Test kubectl access by trying to get the secret we have permissions for
	cmd := exec.Command("kubectl", "get", "secret", SecretName, "-n", DefaultNamespace, "-o", "name")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := stderr.String()
		if strings.Contains(errMsg, "connection refused") || strings.Contains(errMsg, "Unable to connect") {
			return fmt.Errorf("cannot connect to Kubernetes cluster")
		}
		if strings.Contains(errMsg, "NotFound") || strings.Contains(errMsg, "not found") {
			return fmt.Errorf("secret %s not found in namespace %s", SecretName, DefaultNamespace)
		}
		if strings.Contains(errMsg, "Forbidden") {
			return fmt.Errorf("insufficient permissions to access secrets")
		}
		return fmt.Errorf("kubectl error: %v - %s", err, errMsg)
	}

	return nil
}

// GetK8sSecrets retrieves current dogfood secrets from K8s (for verification)
func GetK8sSecrets() (apiURL, liveSecret, testSecret string, err error) {
	if err := checkKubectlAccess(); err != nil {
		return "", "", "", err
	}

	cmd := exec.Command("kubectl", "get", "secret", SecretName, "-n", DefaultNamespace, "-o", "json")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", "", "", fmt.Errorf("failed to get secret: %v - %s", err, stderr.String())
	}

	var secret struct {
		Data map[string]string `json:"data"`
	}

	if err := json.Unmarshal(stdout.Bytes(), &secret); err != nil {
		return "", "", "", fmt.Errorf("failed to parse secret: %v", err)
	}

	if encoded, ok := secret.Data["DOGFOOD_LIVE_SECRET"]; ok {
		decoded, _ := base64.StdEncoding.DecodeString(encoded)
		liveSecret = string(decoded)
	}

	if encoded, ok := secret.Data["DOGFOOD_TEST_SECRET"]; ok {
		decoded, _ := base64.StdEncoding.DecodeString(encoded)
		testSecret = string(decoded)
	}

	return apiURL, liveSecret, testSecret, nil
}
