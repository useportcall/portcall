package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/useportcall/portcall/libs/go/authx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func Auth() routerx.HandlerFunc {
	client, err := authx.New()
	adminAPIKey := os.Getenv("ADMIN_API_KEY")
	if err != nil {
		log.Printf("Auth Middleware - failed to init auth client: %v", err)
		return func(c *routerx.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication unavailable"})
			c.Abort()
		}
	}

	return func(c *routerx.Context) {
		// Check for API key authentication first
		apiKey := c.GetHeader("X-Admin-API-Key")
		if apiKey != "" {
			if adminAPIKey == "" {
				log.Println("Auth Middleware - ADMIN_API_KEY not configured")
				c.JSON(http.StatusUnauthorized, gin.H{"access": "unauthorized", "error": "API key authentication not configured"})
				c.Abort()
				return
			}

			if !strings.HasPrefix(apiKey, "admin_") || apiKey != adminAPIKey {
				log.Println("Auth Middleware - Invalid API key")
				c.JSON(http.StatusUnauthorized, gin.H{"access": "unauthorized", "error": "Invalid API key"})
				c.Abort()
				return
			}

			// API key is valid - set admin context
			c.Set("auth_email", "admin@portcall.internal")
			c.Set("auth_type", "api_key")
			log.Println("Auth Middleware - Authenticated via API key")
			c.Next()
			return
		}

		// Fall back to Keycloak authentication
		claims, err := client.Validate(c.Request.Context(), c.Request.Header)
		if err != nil {
			log.Println("Auth Middleware - Unauthorized access attempt:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"access": "unauthorized"})
			c.Abort()
			return
		}

		c.Set("auth_email", claims.Email)
		c.Set("auth_claims", claims)
		c.Set("auth_type", "keycloak")

		c.Next()
	}
}
