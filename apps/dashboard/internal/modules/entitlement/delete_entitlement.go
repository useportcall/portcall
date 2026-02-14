package entitlement

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func DeleteEntitlement(c *routerx.Context) {
	userID := c.Param("user_id")
	featureID := c.Param("feature_id")

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

	if err := c.DB().Delete(&entitlement, "id = ?", entitlement.ID); err != nil {
		c.ServerError("Failed to delete entitlement", err)
		return
	}

	c.OK(map[string]bool{"deleted": true})
}
