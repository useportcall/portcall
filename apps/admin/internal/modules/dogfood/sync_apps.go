package dogfood

import (
	"log"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/dogfoodx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type SyncAppsRequest struct {
	DryRun bool `json:"dry_run"`
}

type SyncResult struct {
	AppID     string `json:"app_id"`
	AppName   string `json:"app_name"`
	IsLive    bool   `json:"is_live"`
	Success   bool   `json:"success"`
	Error     string `json:"error,omitempty"`
	UserCount int64  `json:"user_count"`
	SubCount  int64  `json:"subscription_count"`
}

type SyncAppsResponse struct {
	TotalApps int          `json:"total_apps"`
	Synced    int          `json:"synced"`
	Skipped   int          `json:"skipped"`
	Failed    int          `json:"failed"`
	DryRun    bool         `json:"dry_run"`
	PlanID    string       `json:"plan_id"`
	Results   []SyncResult `json:"results"`
}

// SyncApps syncs all non-billing-exempt apps to dogfood
// This creates users, subscriptions, and updates usage for all existing apps
func SyncApps(c *routerx.Context) {
	var body SyncAppsRequest
	_ = c.ShouldBindJSON(&body)

	// Find the dogfood account and live app
	var account models.Account
	if err := c.DB().FindFirst(&account, "email = ?", DogfoodAccountEmail); err != nil {
		c.ServerError("Dogfood account not found. Run /api/dogfood/setup first.", err)
		return
	}

	// Find the live app (we use live for production apps)
	var liveApp models.App
	if err := c.DB().FindFirst(&liveApp, "account_id = ? AND name = ?", account.ID, DogfoodLiveAppName); err != nil {
		c.ServerError("Dogfood live app not found. Run /api/dogfood/setup first.", err)
		return
	}

	// Find the free tier plan
	var plan models.Plan
	if err := c.DB().FindFirst(&plan, "app_id = ? AND is_free = ?", liveApp.ID, true); err != nil {
		c.ServerError("Dogfood free tier plan not found. Run /api/dogfood/setup first.", err)
		return
	}

	// Get all non-billing-exempt apps
	var apps []models.App
	if err := c.DB().List(&apps, "billing_exempt = ?", false); err != nil {
		c.ServerError("Failed to list apps", err)
		return
	}

	response := SyncAppsResponse{
		TotalApps: len(apps),
		DryRun:    body.DryRun,
		PlanID:    plan.PublicID,
		Results:   make([]SyncResult, 0, len(apps)),
	}

	for _, app := range apps {
		result := SyncResult{
			AppID:   app.PublicID,
			AppName: app.Name,
			IsLive:  app.IsLive,
		}

		// Count users and subscriptions
		var userCount, subCount int64
		c.DB().Count(&userCount, &models.User{}, "app_id = ?", app.ID)
		c.DB().Count(&subCount, &models.Subscription{}, "app_id = ?", app.ID)
		result.UserCount = userCount
		result.SubCount = subCount

		if body.DryRun {
			result.Success = true
			response.Synced++
		} else {
			if err := dogfoodx.SyncAppToDogfood(c.DB(), &app, plan.PublicID); err != nil {
				log.Printf("[dogfood/sync] Failed to sync app %s: %v", app.PublicID, err)
				result.Success = false
				result.Error = err.Error()
				response.Failed++
			} else {
				log.Printf("[dogfood/sync] Successfully synced app %s (users: %d, subs: %d)", app.PublicID, userCount, subCount)
				result.Success = true
				response.Synced++
			}
		}

		response.Results = append(response.Results, result)
	}

	c.OK(response)
}

// SyncSingleApp syncs a single app to dogfood
func SyncSingleApp(c *routerx.Context) {
	appID := c.Param("app_id")
	if appID == "" {
		c.BadRequest("app_id is required")
		return
	}

	// Find the app
	var app models.App
	if err := c.DB().FindFirst(&app, "public_id = ?", appID); err != nil {
		if dbx.IsRecordNotFoundError(err) {
			c.NotFound("App not found")
			return
		}
		c.ServerError("Failed to find app", err)
		return
	}

	if app.BillingExempt {
		c.BadRequest("Cannot sync billing-exempt apps")
		return
	}

	// Find the dogfood account and live app
	var account models.Account
	if err := c.DB().FindFirst(&account, "email = ?", DogfoodAccountEmail); err != nil {
		c.ServerError("Dogfood account not found", err)
		return
	}

	var liveApp models.App
	if err := c.DB().FindFirst(&liveApp, "account_id = ? AND name = ?", account.ID, DogfoodLiveAppName); err != nil {
		c.ServerError("Dogfood live app not found", err)
		return
	}

	var plan models.Plan
	if err := c.DB().FindFirst(&plan, "app_id = ? AND is_free = ?", liveApp.ID, true); err != nil {
		c.ServerError("Dogfood free tier plan not found", err)
		return
	}

	// Sync the app
	if err := dogfoodx.SyncAppToDogfood(c.DB(), &app, plan.PublicID); err != nil {
		c.ServerError("Failed to sync app", err)
		return
	}

	// Count for response
	var userCount, subCount int64
	c.DB().Count(&userCount, &models.User{}, "app_id = ?", app.ID)
	c.DB().Count(&subCount, &models.Subscription{}, "app_id = ?", app.ID)

	c.OK(map[string]any{
		"message":            "App synced successfully",
		"app_id":             app.PublicID,
		"app_name":           app.Name,
		"plan_id":            plan.PublicID,
		"user_count":         userCount,
		"subscription_count": subCount,
	})
}
