package entitlement

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type UpdateEntitlementRequest struct {
	Quota *int64 `json:"quota"`
	Usage *int64 `json:"usage"`
}

func UpdateEntitlement(c *routerx.Context) {
	userID := c.Param("user_id")
	featureID := c.Param("feature_id")

	var body UpdateEntitlementRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	var user models.User
	if err := c.DB().GetForPublicID(c.AppID(), userID, &user); err != nil {
		c.NotFound("User not found")
		return
	}

	var entitlement models.Entitlement
	if err := c.DB().FindFirst(&entitlement, "app_id = ? AND user_id = ? AND feature_public_id = ?", c.AppID(), user.ID, featureID); err != nil {
		c.NotFound("Entitlement not found")
		return
	}

	if body.Quota != nil {
		entitlement.Quota = *body.Quota
	}

	if body.Usage != nil {
		entitlement.Usage = *body.Usage
	}

	if err := c.DB().Save(&entitlement); err != nil {
		c.ServerError("Failed to save entitlement", err)
		return
	}

	c.OK(new(apix.Entitlement).Set(&entitlement))
}
