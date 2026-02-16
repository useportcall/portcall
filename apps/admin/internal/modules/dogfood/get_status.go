package dogfood

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type DogfoodStatusResponse struct {
	Configured bool          `json:"configured"`
	Account    *AccountInfo  `json:"account,omitempty"`
	LiveApp    *AppInfo      `json:"live_app,omitempty"`
	TestApp    *AppInfo      `json:"test_app,omitempty"`
	HasSecrets bool          `json:"has_secrets"`
	Plan       *PlanInfo     `json:"plan,omitempty"`
	Features   []FeatureInfo `json:"features,omitempty"`
	UserCount  int64         `json:"user_count"`
}

// GetDogfoodStatus returns the current status of the dogfood setup
func GetDogfoodStatus(c *routerx.Context) {
	response := DogfoodStatusResponse{}

	// Check for dogfood account
	var account models.Account
	err := c.DB().FindFirst(&account, "email = ?", DogfoodAccountEmail)
	if err != nil {
		if dbx.IsRecordNotFoundError(err) {
			c.OK(response)
			return
		}
		c.ServerError("Failed to check dogfood status", err)
		return
	}

	response.Configured = true
	response.Account = &AccountInfo{ID: account.ID, Email: account.Email}

	// Check for apps
	var apps []models.App
	if err := c.DB().List(&apps, "account_id = ?", account.ID); err != nil {
		c.ServerError("Failed to list apps", err)
		return
	}

	for _, app := range apps {
		info := &AppInfo{
			ID:       app.ID,
			PublicID: app.PublicID,
			Name:     app.Name,
			IsLive:   app.IsLive,
		}
		if app.IsLive {
			response.LiveApp = info
		} else {
			response.TestApp = info
		}
	}

	// Check for secrets on live app
	if response.LiveApp != nil {
		var secretCount int64
		c.DB().Count(&secretCount, models.Secret{}, "app_id = ? AND disabled_at IS NULL", response.LiveApp.ID)
		response.HasSecrets = secretCount > 0

		// Get all features
		var features []models.Feature
		if err := c.DB().List(&features, "app_id = ?", response.LiveApp.ID); err == nil {
			response.Features = make([]FeatureInfo, len(features))
			for i, f := range features {
				response.Features[i] = FeatureInfo{ID: f.ID, PublicID: f.PublicID}
			}
		}

		// Check for plan
		var plan models.Plan
		if err := c.DB().FindFirst(&plan, "app_id = ? AND status = ?", response.LiveApp.ID, "active"); err == nil {
			response.Plan = &PlanInfo{ID: plan.ID, PublicID: plan.PublicID, Name: plan.Name}
		}

		// Count users (each app using portcall is a "user" in the dogfood account)
		c.DB().Count(&response.UserCount, models.User{}, "app_id = ?", response.LiveApp.ID)
	}

	c.OK(response)
}
