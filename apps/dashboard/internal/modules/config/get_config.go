package config

import (
	"os"

	"github.com/useportcall/portcall/libs/go/routerx"
)

// GetConfig returns runtime configuration for the frontend
func GetConfig(c *routerx.Context) {
	keycloakURL := os.Getenv("KEYCLOAK_API_URL")
	if keycloakURL == "" {
		keycloakURL = "http://localhost:8090" // fallback for local dev
	}

	c.OK(map[string]any{
		"keycloak_url": keycloakURL,
	})
}
