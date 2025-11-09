package entitlement

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ListEntitlements(c *routerx.Context) {
	userID := c.Query("user_id")
	filter := c.Query("filter")

	var user models.User
	if err := c.DB().GetForPublicID(c.AppID(), userID, &user); err != nil {
		c.ServerError("Failed to find user for entitlements")
		return
	}

	var query []any
	switch filter {
	case "exclude_metered":
		query = []any{"app_id = ? AND user_id = ? AND is_metered = ?", c.AppID, user.ID, false}
	case "include_metered":
		query = []any{"app_id = ? AND user_id = ? AND is_metered = ?", c.AppID, user.ID, true}
	default:
		query = []any{"app_id = ? AND user_id = ?", c.AppID, user.ID}
	}

	var entitlements []models.Entitlement
	if err := c.DB().List(&entitlements, query...); err != nil {
		c.ServerError("Failed to list entitlements")
		return
	}

	response := make([]Entitlement, len(entitlements))
	for i, entitlement := range entitlements {
		response[i].Set(&entitlement)

		var feature models.Feature
		if err := c.DB().FindForID(entitlement.FeatureID, &feature); err != nil {
			c.ServerError("Failed to find feature for entitlement")
			return
		}

		response[i].Feature = feature.PublicID
	}

	c.OK(response)
}
