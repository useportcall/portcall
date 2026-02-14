package entitlement

import (
	"time"

	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type ToggleEntitlementRequest struct {
	FeatureID string `json:"feature_id"`
	Enabled   bool   `json:"enabled"`
}

func ToggleEntitlement(c *routerx.Context) {
	userID := c.Param("user_id")

	var body ToggleEntitlementRequest
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
	var entitlement models.Entitlement
	if err := c.DB().FindFirstOrNil(&entitlement, "app_id = ? AND user_id = ? AND feature_public_id = ?", c.AppID(), user.ID, body.FeatureID); err != nil {
		c.ServerError("Failed to check entitlement", err)
		return
	}
	entitlementExists := entitlement.ID != 0

	if body.Enabled {
		// Create or update entitlement
		if !entitlementExists {
			// Create new entitlement
			now := time.Now()
			entitlement = models.Entitlement{
				AppID:           c.AppID(),
				UserID:          user.ID,
				FeaturePublicID: body.FeatureID,
				Interval:        "month",
				Quota:           -1, // -1 means unlimited/enabled
				Usage:           0,
				IsMetered:       feature.IsMetered,
				LastResetAt:     &now,
				NextResetAt:     &now,
			}
		} else {
			// Already exists, ensure it's enabled
			entitlement.Quota = -1
		}

		if err := c.DB().Save(&entitlement); err != nil {
			c.ServerError("Failed to save entitlement", err)
			return
		}

		c.OK(new(apix.Entitlement).Set(&entitlement))
	} else {
		// Delete entitlement if it exists
		if entitlementExists {
			if err := c.DB().Delete(&entitlement, "id = ?", entitlement.ID); err != nil {
				c.ServerError("Failed to delete entitlement", err)
				return
			}
		}

		c.OK(map[string]bool{"deleted": true})
	}
}
