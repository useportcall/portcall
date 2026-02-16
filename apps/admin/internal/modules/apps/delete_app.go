package apps

import (
	"strconv"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

// DeleteAppResponse represents the result of deleting an app
type DeleteAppResponse struct {
	Success      bool   `json:"success"`
	AppID        uint   `json:"app_id"`
	AppPublicID  string `json:"app_public_id"`
	AppName      string `json:"app_name"`
	Message      string `json:"message"`
	DeletedItems struct {
		Users         int `json:"users"`
		Subscriptions int `json:"subscriptions"`
		Plans         int `json:"plans"`
		Connections   int `json:"connections"`
		Secrets       int `json:"secrets"`
	} `json:"deleted_items"`
}

// DeleteApp deletes an app and all its associated data
// CAUTION: This is a destructive operation
func DeleteApp(c *routerx.Context) {
	appIDStr := c.Param("app_id")
	appID, err := strconv.ParseUint(appIDStr, 10, 64)
	if err != nil {
		c.BadRequest("Invalid app ID")
		return
	}

	var app models.App
	if err := c.DB().FindForID(uint(appID), &app); err != nil {
		c.NotFound("App not found")
		return
	}

	response := DeleteAppResponse{
		AppID:       app.ID,
		AppPublicID: app.PublicID,
		AppName:     app.Name,
	}

	// Delete in order of dependencies
	// 1. Delete subscriptions first (they reference users and plans)
	var subscriptions []models.Subscription
	c.DB().List(&subscriptions, "app_id = ?", app.ID)
	for _, sub := range subscriptions {
		c.DB().DeleteForID(&sub)
		response.DeletedItems.Subscriptions++
	}

	// 2. Delete users
	var users []models.User
	c.DB().List(&users, "app_id = ?", app.ID)
	for _, user := range users {
		c.DB().DeleteForID(&user)
		response.DeletedItems.Users++
	}

	// 3. Delete plan items and features
	var plans []models.Plan
	c.DB().List(&plans, "app_id = ?", app.ID)
	for _, plan := range plans {
		// Delete plan items
		var planItems []models.PlanItem
		c.DB().List(&planItems, "plan_id = ?", plan.ID)
		for _, item := range planItems {
			c.DB().DeleteForID(&item)
		}
		// Delete plan features
		var planFeatures []models.PlanFeature
		c.DB().List(&planFeatures, "plan_id = ?", plan.ID)
		for _, feature := range planFeatures {
			c.DB().DeleteForID(&feature)
		}
		c.DB().DeleteForID(&plan)
		response.DeletedItems.Plans++
	}

	// 4. Delete connections
	var connections []models.Connection
	c.DB().List(&connections, "app_id = ?", app.ID)
	for _, conn := range connections {
		c.DB().DeleteForID(&conn)
		response.DeletedItems.Connections++
	}

	// 5. Delete secrets
	var secrets []models.Secret
	c.DB().List(&secrets, "app_id = ?", app.ID)
	for _, secret := range secrets {
		c.DB().DeleteForID(&secret)
		response.DeletedItems.Secrets++
	}

	// 6. Delete app config
	var appConfig models.AppConfig
	if err := c.DB().FindFirst(&appConfig, "app_id = ?", app.ID); err == nil {
		c.DB().DeleteForID(&appConfig)
	}

	// 7. Delete company and address
	var company models.Company
	if err := c.DB().FindFirst(&company, "app_id = ?", app.ID); err == nil {
		c.DB().DeleteForID(&company)
	}

	// 8. Finally delete the app itself
	if err := c.DB().DeleteForID(&app); err != nil {
		c.ServerError("Failed to delete app", err)
		return
	}

	response.Success = true
	response.Message = "App and all associated data deleted successfully"

	c.OK(response)
}
