package dogfood

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

// SetupPaymentConnection creates a local mock payment connection for the dogfood live app
func SetupPaymentConnection(c *routerx.Context) {
	// Get the dogfood account
	var account models.Account
	if err := c.DB().FindFirst(&account, "email = ?", DogfoodAccountEmail); err != nil {
		c.ServerError("Failed to find dogfood account", err)
		return
	}

	// Get the live app
	var liveApp models.App
	if err := c.DB().FindFirst(&liveApp, "account_id = ? AND is_live = ?", account.ID, true); err != nil {
		c.ServerError("Failed to find dogfood live app", err)
		return
	}

	// Check if a connection already exists
	var existingConnection models.Connection
	if err := c.DB().FindFirst(&existingConnection, "app_id = ?", liveApp.ID); err == nil {
		// Connection exists, check if it's set as default
		var appConfig models.AppConfig
		if err := c.DB().FindFirst(&appConfig, "app_id = ?", liveApp.ID); err == nil {
			if appConfig.DefaultConnectionID == existingConnection.ID {
				c.OK(map[string]interface{}{
					"message":         "Payment connection already configured",
					"connection_id":   existingConnection.ID,
					"connection_name": existingConnection.Name,
					"is_default":      true,
				})
				return
			} else {
				// Update to set as default
				appConfig.DefaultConnectionID = existingConnection.ID
				if err := c.DB().Save(&appConfig); err != nil {
					c.ServerError("Failed to set connection as default", err)
					return
				}
				c.OK(map[string]interface{}{
					"message":         "Existing connection set as default",
					"connection_id":   existingConnection.ID,
					"connection_name": existingConnection.Name,
					"is_default":      true,
				})
				return
			}
		}
	}

	// Create a new local mock connection
	connection := &models.Connection{
		PublicID:     dbx.GenPublicID("conn"),
		AppID:        liveApp.ID,
		Name:         "Local Mock Payment",
		Source:       "local",
		PublicKey:    "pk_local",
		EncryptedKey: "", // Local connections don't need a key
	}

	if err := c.DB().Create(connection); err != nil {
		c.ServerError("Failed to create payment connection", err)
		return
	}

	// Set as default connection in app config
	var appConfig models.AppConfig
	if err := c.DB().FindFirst(&appConfig, "app_id = ?", liveApp.ID); err != nil {
		// Create app config if it doesn't exist
		appConfig = models.AppConfig{
			AppID:               liveApp.ID,
			DefaultConnectionID: connection.ID,
		}
		if err := c.DB().Create(&appConfig); err != nil {
			c.ServerError("Failed to create app config", err)
			return
		}
	} else {
		// Update existing config
		appConfig.DefaultConnectionID = connection.ID
		if err := c.DB().Save(&appConfig); err != nil {
			c.ServerError("Failed to update app config", err)
			return
		}
	}

	c.OK(map[string]interface{}{
		"message":         "Payment connection created and set as default",
		"connection_id":   connection.ID,
		"connection_name": connection.Name,
		"public_id":       connection.PublicID,
		"is_default":      true,
	})
}
