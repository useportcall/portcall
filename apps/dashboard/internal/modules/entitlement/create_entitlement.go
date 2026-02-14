package entitlement

import (
	"time"

	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type CreateEntitlementRequest struct {
	FeatureID string `json:"feature_id" binding:"required"`
	Quota     int64  `json:"quota"`
	Interval  string `json:"interval"`
}

func CreateEntitlement(c *routerx.Context) {
	userID := c.Param("user_id")

	var body CreateEntitlementRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	var user models.User
	if err := c.DB().GetForPublicID(c.AppID(), userID, &user); err != nil {
		c.NotFound("User not found")
		return
	}

	// Check if feature exists
	var feature models.Feature
	if err := c.DB().FindFirst(&feature, "app_id = ? AND public_id = ?", c.AppID(), body.FeatureID); err != nil {
		c.NotFound("Feature not found")
		return
	}

	// Check if entitlement already exists
	var existingEntitlement models.Entitlement
	if err := c.DB().FindFirstOrNil(&existingEntitlement, "app_id = ? AND user_id = ? AND feature_public_id = ?", c.AppID(), user.ID, body.FeatureID); err != nil {
		c.ServerError("Failed to check entitlement", err)
		return
	}

	if existingEntitlement.ID != 0 {
		c.BadRequest("Entitlement already exists for this feature")
		return
	}

	// Set defaults
	interval := body.Interval
	if interval == "" {
		interval = "month"
	}

	quota := body.Quota
	if quota == 0 && !feature.IsMetered {
		quota = -1 // -1 means unlimited/enabled for basic features
	}

	now := time.Now()
	entitlement := models.Entitlement{
		AppID:           c.AppID(),
		UserID:          user.ID,
		FeaturePublicID: body.FeatureID,
		Interval:        interval,
		Quota:           quota,
		Usage:           0,
		IsMetered:       feature.IsMetered,
		LastResetAt:     &now,
		NextResetAt:     &now,
	}

	if err := c.DB().Create(&entitlement); err != nil {
		c.ServerError("Failed to create entitlement", err)
		return
	}

	c.OK(new(apix.Entitlement).Set(&entitlement))
}
