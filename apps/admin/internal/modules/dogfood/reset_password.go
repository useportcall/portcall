package dogfood

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/useportcall/portcall/libs/go/routerx"
)

// DogfoodKeycloakUser represents the dogfood user in Keycloak
// This user is used to log in to the dashboard as the dogfood account owner
const (
	DogfoodKeycloakUsername = "dogfood"
	DogfoodKeycloakEmail    = "dogfood@useportcall.com"
	DefaultDogfoodRealm     = "dev" // The realm where dashboard users live (can be overridden via DOGFOOD_KEYCLOAK_REALM env var)
)

// getDogfoodRealm returns the realm to use for the dogfood user
// Can be configured via DOGFOOD_KEYCLOAK_REALM environment variable
func getDogfoodRealm() string {
	realm := os.Getenv("DOGFOOD_KEYCLOAK_REALM")
	if realm == "" {
		realm = os.Getenv("KEYCLOAK_REALM")
	}
	if realm == "" {
		realm = DefaultDogfoodRealm
	}
	// The dogfood user lives in the dashboard/dev realm, not the admin realm
	// If we're configured with the admin realm, use "dev" instead
	if realm == "admin" {
		realm = DefaultDogfoodRealm
	}
	return realm
}

type ResetPasswordResponse struct {
	Success  bool   `json:"success"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Realm    string `json:"realm"`
	Message  string `json:"message"`
}

type KeycloakTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type KeycloakUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Enabled  bool   `json:"enabled"`
}

// ResetDogfoodPassword resets the password for the dogfood user in Keycloak
// and returns the new password so the admin can log in to the dashboard
func ResetDogfoodPassword(c *routerx.Context) {
	// Get Keycloak configuration from environment
	keycloakURL := os.Getenv("KEYCLOAK_INTERNAL_URL")
	if keycloakURL == "" {
		keycloakURL = os.Getenv("KEYCLOAK_API_URL")
	}
	if keycloakURL == "" {
		c.BadRequest("KEYCLOAK_API_URL not configured")
		return
	}

	// Get admin credentials from environment or K8s secrets
	adminUsername := os.Getenv("KC_BOOTSTRAP_ADMIN_USERNAME")
	if adminUsername == "" {
		adminUsername = "admin"
	}
	adminPassword := os.Getenv("KC_BOOTSTRAP_ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = os.Getenv("KEYCLOAK_ADMIN_PASSWORD")
	}
	if adminPassword == "" {
		c.BadRequest("Keycloak admin password not configured (set KC_BOOTSTRAP_ADMIN_PASSWORD or KEYCLOAK_ADMIN_PASSWORD)")
		return
	}

	// Normalize the URL
	keycloakURL = strings.TrimSuffix(keycloakURL, "/")

	// Get the realm for dogfood user
	realm := getDogfoodRealm()

	log.Printf("[dogfood] Resetting password for dogfood user in Keycloak realm: %s", realm)

	// 1. Get admin access token
	token, err := getKeycloakAdminToken(keycloakURL, adminUsername, adminPassword)
	if err != nil {
		log.Printf("[dogfood] Failed to get Keycloak admin token: %v", err)
		c.ServerError("Failed to authenticate with Keycloak", err)
		return
	}

	// 2. Find or create the dogfood user
	user, err := findOrCreateKeycloakUser(keycloakURL, token, realm, DogfoodKeycloakUsername, DogfoodKeycloakEmail)
	if err != nil {
		log.Printf("[dogfood] Failed to find/create dogfood user: %v", err)
		c.ServerError("Failed to find or create dogfood user in Keycloak", err)
		return
	}

	// 3. Generate a new secure password
	newPassword, err := generateSecurePassword(24)
	if err != nil {
		log.Printf("[dogfood] Failed to generate password: %v", err)
		c.ServerError("Failed to generate secure password", err)
		return
	}

	// 4. Set the new password in Keycloak
	err = setKeycloakUserPassword(keycloakURL, token, realm, user.ID, newPassword)
	if err != nil {
		log.Printf("[dogfood] Failed to set password: %v", err)
		c.ServerError("Failed to set password in Keycloak", err)
		return
	}

	log.Printf("[dogfood] Successfully reset password for user %s in realm %s", user.Username, realm)

	c.OK(ResetPasswordResponse{
		Success:  true,
		Username: user.Username,
		Email:    user.Email,
		Password: newPassword,
		Realm:    realm,
		Message:  "Password has been reset. Use these credentials to log in to the dashboard.",
	})
}

// getKeycloakAdminToken gets an admin access token from Keycloak
func getKeycloakAdminToken(keycloakURL, username, password string) (string, error) {
	tokenURL := fmt.Sprintf("%s/realms/master/protocol/openid-connect/token", keycloakURL)

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", "admin-cli")
	data.Set("username", username)
	data.Set("password", password)

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to request token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp KeycloakTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("received empty access token")
	}

	return tokenResp.AccessToken, nil
}

// findOrCreateKeycloakUser finds an existing user by username or creates a new one
func findOrCreateKeycloakUser(keycloakURL, token, realm, username, email string) (*KeycloakUser, error) {
	// Try to find existing user
	searchURL := fmt.Sprintf("%s/admin/realms/%s/users?username=%s&exact=true", keycloakURL, realm, url.QueryEscape(username))

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user search failed (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var users []KeycloakUser
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, fmt.Errorf("failed to parse users: %w", err)
	}

	// If user exists, return it
	if len(users) > 0 {
		return &users[0], nil
	}

	// Create new user
	log.Printf("[dogfood] Creating new Keycloak user: %s", username)

	newUser := map[string]interface{}{
		"username":      username,
		"email":         email,
		"enabled":       true,
		"emailVerified": true,
		"firstName":     "Dogfood",
		"lastName":      "Account",
	}

	userJSON, err := json.Marshal(newUser)
	if err != nil {
		return nil, err
	}

	createURL := fmt.Sprintf("%s/admin/realms/%s/users", keycloakURL, realm)
	req, err = http.NewRequest("POST", createURL, bytes.NewReader(userJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user creation failed (HTTP %d): %s", resp.StatusCode, string(body))
	}

	// Fetch the created user to get the ID
	req, err = http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch created user: %w", err)
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, fmt.Errorf("failed to parse created user: %w", err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user was created but could not be retrieved")
	}

	return &users[0], nil
}

// setKeycloakUserPassword sets a new password for a user in Keycloak
func setKeycloakUserPassword(keycloakURL, token, realm, userID, password string) error {
	passwordURL := fmt.Sprintf("%s/admin/realms/%s/users/%s/reset-password", keycloakURL, realm, userID)

	credentialData := map[string]interface{}{
		"type":      "password",
		"value":     password,
		"temporary": false,
	}

	credJSON, err := json.Marshal(credentialData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", passwordURL, bytes.NewReader(credJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("password reset failed (HTTP %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// generateSecurePassword generates a cryptographically secure password
func generateSecurePassword(length int) (string, error) {
	// Use a mix of characters for a strong password
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Base64 encode and trim to desired length
	password := base64.URLEncoding.EncodeToString(bytes)
	if len(password) > length {
		password = password[:length]
	}

	// Replace some characters to make it more readable
	password = strings.ReplaceAll(password, "-", "!")
	password = strings.ReplaceAll(password, "_", "@")

	return password, nil
}
