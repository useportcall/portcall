package connections

import (
	"fmt"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

// MockConnectionInfo represents info about an app's mock connections
type MockConnectionInfo struct {
	AppID            uint   `json:"app_id"`
	AppPublicID      string `json:"app_public_id"`
	AppName          string `json:"app_name"`
	IsLive           bool   `json:"is_live"`
	MockConnectionID uint   `json:"mock_connection_id,omitempty"`
	ConnectionName   string `json:"connection_name,omitempty"`
}

// CheckMockConnectionsResponse shows which live apps have mock connections
type CheckMockConnectionsResponse struct {
	HasMockConnections bool                 `json:"has_mock_connections"`
	LiveAppsWithMock   []MockConnectionInfo `json:"live_apps_with_mock"`
	TotalLiveApps      int                  `json:"total_live_apps"`
	TotalWithMock      int                  `json:"total_with_mock"`
}

// CheckMockConnections checks if any live apps have mock/local payment connections
func CheckMockConnections(c *routerx.Context) {
	// Find all live apps
	var liveApps []models.App
	if err := c.DB().List(&liveApps, "is_live = ?", true); err != nil {
		c.ServerError("Failed to list apps", err)
		return
	}

	response := CheckMockConnectionsResponse{
		LiveAppsWithMock: make([]MockConnectionInfo, 0),
		TotalLiveApps:    len(liveApps),
	}

	for _, app := range liveApps {
		// Find mock connections for this app
		var connections []models.Connection
		if err := c.DB().List(&connections, "app_id = ? AND source = ?", app.ID, "local"); err != nil {
			continue
		}

		for _, conn := range connections {
			response.LiveAppsWithMock = append(response.LiveAppsWithMock, MockConnectionInfo{
				AppID:            app.ID,
				AppPublicID:      app.PublicID,
				AppName:          app.Name,
				IsLive:           app.IsLive,
				MockConnectionID: conn.ID,
				ConnectionName:   conn.Name,
			})
		}
	}

	response.TotalWithMock = len(response.LiveAppsWithMock)
	response.HasMockConnections = response.TotalWithMock > 0

	c.OK(response)
}

// ClearMockConnectionsResponse shows the result of clearing mock connections
type ClearMockConnectionsResponse struct {
	Cleared  int      `json:"cleared"`
	Errors   []string `json:"errors,omitempty"`
	AppIDs   []uint   `json:"app_ids"`
	Messages []string `json:"messages"`
}

// ClearMockConnections removes all mock/local payment connections from live apps
func ClearMockConnections(c *routerx.Context) {
	// Find all live apps
	var liveApps []models.App
	if err := c.DB().List(&liveApps, "is_live = ?", true); err != nil {
		c.ServerError("Failed to list apps", err)
		return
	}

	response := ClearMockConnectionsResponse{
		AppIDs:   make([]uint, 0),
		Messages: make([]string, 0),
		Errors:   make([]string, 0),
	}

	for _, app := range liveApps {
		// Find mock connections for this app
		var connections []models.Connection
		if err := c.DB().List(&connections, "app_id = ? AND source = ?", app.ID, "local"); err != nil {
			response.Errors = append(response.Errors, fmt.Sprintf("Failed to list connections for app %d: %v", app.ID, err))
			continue
		}

		for _, conn := range connections {
			// Check if this is the default connection in app config
			var appConfig models.AppConfig
			if err := c.DB().FindFirst(&appConfig, "app_id = ?", app.ID); err == nil {
				if appConfig.DefaultConnectionID == conn.ID {
					// Need to set default_connection_id to NULL using raw SQL
					// because the field is uint and can't be set to nil in Go
					if err := c.DB().Exec("UPDATE app_configs SET default_connection_id = NULL WHERE app_id = ?", app.ID); err != nil {
						response.Errors = append(response.Errors, fmt.Sprintf("Failed to clear default connection for app %d: %v", app.ID, err))
						continue
					}
				}
			}

			// Delete the mock connection
			if err := c.DB().DeleteForID(&conn); err != nil {
				response.Errors = append(response.Errors, fmt.Sprintf("Failed to delete connection %d: %v", conn.ID, err))
				continue
			}

			response.Cleared++
			response.AppIDs = append(response.AppIDs, app.ID)
			response.Messages = append(response.Messages, fmt.Sprintf("Cleared mock connection '%s' from app '%s' (ID: %d)", conn.Name, app.Name, app.ID))
		}
	}

	c.OK(response)
}

// ListConnectionsForApp lists all connections for a specific app
func ListConnectionsForApp(c *routerx.Context) {
	appID := c.Param("app_id")
	if appID == "" {
		c.BadRequest("App ID is required")
		return
	}

	var app models.App
	if err := c.DB().FindFirst(&app, "id = ? OR public_id = ?", appID, appID); err != nil {
		c.NotFound("App not found")
		return
	}

	var connections []models.Connection
	if err := c.DB().List(&connections, "app_id = ?", app.ID); err != nil {
		c.ServerError("Failed to list connections", err)
		return
	}

	type ConnectionInfo struct {
		ID        uint   `json:"id"`
		PublicID  string `json:"public_id"`
		Name      string `json:"name"`
		Source    string `json:"source"`
		PublicKey string `json:"public_key"`
		IsDefault bool   `json:"is_default"`
	}

	// Get app config to determine default
	var appConfig models.AppConfig
	_ = c.DB().FindFirst(&appConfig, "app_id = ?", app.ID)

	result := make([]ConnectionInfo, len(connections))
	for i, conn := range connections {
		result[i] = ConnectionInfo{
			ID:        conn.ID,
			PublicID:  conn.PublicID,
			Name:      conn.Name,
			Source:    conn.Source,
			PublicKey: conn.PublicKey,
			IsDefault: appConfig.DefaultConnectionID == conn.ID,
		}
	}

	c.OK(result)
}
