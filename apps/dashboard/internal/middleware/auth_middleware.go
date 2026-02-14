package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/useportcall/portcall/libs/go/authx"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func Auth(db dbx.IORM) routerx.HandlerFunc {
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
		// Check for Admin API key authentication first
		// This allows admin service to make calls on behalf of apps
		apiKey := c.GetHeader("X-Admin-API-Key")
		targetAppID := c.GetHeader("X-Target-App-ID")

		if apiKey != "" && targetAppID != "" {
			if adminAPIKey == "" {
				log.Println("Auth Middleware - ADMIN_API_KEY not configured")
				c.JSON(http.StatusUnauthorized, gin.H{"access": "unauthorized", "error": "Admin API key authentication not configured"})
				c.Abort()
				return
			}

			if !strings.HasPrefix(apiKey, "admin_") || apiKey != adminAPIKey {
				log.Println("Auth Middleware - Invalid admin API key")
				c.JSON(http.StatusUnauthorized, gin.H{"access": "unauthorized", "error": "Invalid admin API key"})
				c.Abort()
				return
			}

			// Find the target app
			var app models.App
			if err := db.FindFirst(&app, "public_id = ? OR id = ?", targetAppID, targetAppID); err != nil {
				log.Println("Auth Middleware - Target app not found:", err)
				c.JSON(http.StatusNotFound, gin.H{"error": "Target app not found"})
				c.Abort()
				return
			}

			// Set admin context with app info
			c.Set("auth_email", "admin@portcall.internal")
			c.Set("auth_type", "admin_api_key")
			c.Set("app_id", app.ID)
			c.Set("public_app_id", app.PublicID)
			c.Set("is_live", app.IsLive)
			c.Set("is_billing_exempt", app.BillingExempt)
			c.Set("target_app", app)
			log.Printf("Auth Middleware - Authenticated via admin API key for app: %s", app.PublicID)
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

		if c.Param("app_id") == "" {
			c.Next()
			return
		}

		var app models.App
		if err := db.FindFirst(&app, "public_id = ?", c.Param("app_id")); err != nil {
			log.Println("Auth Middleware - App not found:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "App not found"})
			c.Abort()
			return
		}

		c.Set("app_id", app.ID)
		c.Set("public_app_id", app.PublicID)
		c.Set("is_live", app.IsLive)
		c.Set("is_billing_exempt", app.BillingExempt)

		c.Next()
	}
}
