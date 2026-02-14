package entitlement

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ListEntitlements(c *routerx.Context) {
	userID := c.Query("user_id")
	isMetered := c.QueryBool("is_metered", false)

	var user models.User
	if err := c.DB().GetForPublicID(c.AppID(), userID, &user); err != nil {
		if dbx.IsRecordNotFoundError(err) {
			c.NotFound("User not found")
			return
		}
		c.ServerError("Failed to find user for entitlements", err)
		return
	}

	var query []any
	switch isMetered {
	case false:
		query = []any{"app_id = ? AND user_id = ? AND is_metered = ?", c.AppID(), user.ID, false}
	case true:
		query = []any{"app_id = ? AND user_id = ? AND is_metered = ?", c.AppID(), user.ID, true}
	}

	var entitlements []models.Entitlement
	if err := c.DB().List(&entitlements, query...); err != nil {
		c.ServerError("Failed to list entitlements", err)
		return
	}

	response := make([]apix.Entitlement, len(entitlements))
	for i, entitlement := range entitlements {
		response[i].Set(&entitlement)

		var feature models.Feature
		if err := c.DB().GetForPublicID(c.AppID(), entitlement.FeaturePublicID, &feature); err != nil {
			c.ServerError("Failed to find feature for entitlement", err)
			return
		}
	}

	c.OK(response)
}
