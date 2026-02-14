package app_config

import (
	"log"

	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type UpdateAppConfigRequest struct {
	ConnectionID string `json:"connection_id"`
}

func UpdateAppConfig(c *routerx.Context) {
	var body UpdateAppConfigRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Printf("[UpdateAppConfig] Invalid request body: %v", err)
		c.BadRequest("Invalid request body")
		return
	}

	log.Printf("[UpdateAppConfig] AppID=%d, ConnectionID=%s", c.AppID(), body.ConnectionID)

	// Find or create the app config
	var appConfig models.AppConfig
	if err := c.DB().FindFirst(&appConfig, "app_id = ?", c.AppID()); err != nil {
		if !dbx.IsRecordNotFoundError(err) {
			log.Printf("[UpdateAppConfig] Error finding app config: %v", err)
			c.ServerError("Error finding app config", err)
			return
		}

		// Create a new app config
		log.Printf("[UpdateAppConfig] Creating new app config for app %d", c.AppID())
		appConfig = models.AppConfig{
			AppID: c.AppID(),
		}
	}

	// Find the connection by public ID
	var connection models.Connection
	if err := c.DB().GetForPublicID(c.AppID(), body.ConnectionID, &connection); err != nil {
		log.Printf("[UpdateAppConfig] Connection not found: %s, error: %v", body.ConnectionID, err)
		c.NotFound("Connection not found")
		return
	}

	log.Printf("[UpdateAppConfig] Found connection ID=%d, Source=%s, Name=%s", connection.ID, connection.Source, connection.Name)

	// Update the default connection
	appConfig.DefaultConnectionID = connection.ID

	// Save the app config
	if appConfig.ID == 0 {
		if err := c.DB().Create(&appConfig); err != nil {
			log.Printf("[UpdateAppConfig] Error creating app config: %v", err)
			c.ServerError("Error creating app config", err)
			return
		}
		log.Printf("[UpdateAppConfig] Created app config ID=%d", appConfig.ID)
	} else {
		if err := c.DB().Save(&appConfig); err != nil {
			log.Printf("[UpdateAppConfig] Error saving app config: %v", err)
			c.ServerError("Error saving app config", err)
			return
		}
		log.Printf("[UpdateAppConfig] Updated app config ID=%d", appConfig.ID)
	}

	// Reload with the connection association
	appConfig.DefaultConnection = connection

	log.Printf("[UpdateAppConfig] Successfully updated default connection to %s", connection.PublicID)
	c.OK(new(apix.AppConfig).Set(&appConfig))
}
