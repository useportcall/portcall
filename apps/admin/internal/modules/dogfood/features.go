package dogfood

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

// FeatureResponse represents a feature in the response
type FeatureResponse struct {
	ID        uint   `json:"id"`
	PublicID  string `json:"public_id"`
	IsMetered bool   `json:"is_metered"`
}

// ListFeatures lists all features for a dogfood app
func ListFeatures(c *routerx.Context) {
	appIDStr := c.Param("app_id")
	if appIDStr == "" {
		c.BadRequest("Missing app_id parameter")
		return
	}

	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		c.BadRequest("Invalid app_id parameter")
		return
	}

	// Verify the app belongs to dogfood account
	var app models.App
	if err := c.DB().FindForID(uint(appID), &app); err != nil {
		c.NotFound("App not found")
		return
	}

	var account models.Account
	if err := c.DB().FindForID(app.AccountID, &account); err != nil || account.Email != DogfoodAccountEmail {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied - not a dogfood app"})
		return
	}

	// Get all features for this app
	var features []models.Feature
	if err := c.DB().List(&features, "app_id = ?", uint(appID)); err != nil {
		c.ServerError("Failed to list features", err)
		return
	}

	result := make([]FeatureResponse, len(features))
	for i, f := range features {
		result[i] = FeatureResponse{
			ID:        f.ID,
			PublicID:  f.PublicID,
			IsMetered: f.IsMetered,
		}
	}

	c.OK(result)
}

// CreateFeatureRequest represents a request to create a feature
type CreateFeatureRequest struct {
	PublicID  string `json:"public_id" binding:"required"`
	IsMetered bool   `json:"is_metered"`
}

// CreateFeature creates a new feature for a dogfood app
func CreateFeature(c *routerx.Context) {
	appIDStr := c.Param("app_id")
	if appIDStr == "" {
		c.BadRequest("Missing app_id parameter")
		return
	}

	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		c.BadRequest("Invalid app_id parameter")
		return
	}

	// Verify the app belongs to dogfood account
	var app models.App
	if err := c.DB().FindForID(uint(appID), &app); err != nil {
		c.NotFound("App not found")
		return
	}

	var account models.Account
	if err := c.DB().FindForID(app.AccountID, &account); err != nil || account.Email != DogfoodAccountEmail {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied - not a dogfood app"})
		return
	}

	var body CreateFeatureRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	if body.PublicID == "" {
		c.BadRequest("Missing public_id")
		return
	}

	// Check if feature already exists
	var existing models.Feature
	if err := c.DB().FindFirst(&existing, "app_id = ? AND public_id = ?", uint(appID), body.PublicID); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Feature already exists"})
		return
	}

	// Create feature
	feature := models.Feature{
		PublicID:  body.PublicID,
		AppID:     uint(appID),
		IsMetered: body.IsMetered,
	}

	if err := c.DB().Create(&feature); err != nil {
		c.ServerError("Failed to create feature", err)
		return
	}

	c.OK(FeatureResponse{
		ID:        feature.ID,
		PublicID:  feature.PublicID,
		IsMetered: feature.IsMetered,
	})
}
