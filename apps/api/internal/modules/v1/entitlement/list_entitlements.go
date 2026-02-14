package entitlement

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ListEntitlements(c *routerx.Context) {
	userID := c.Query("user_id")

	if userID == "" {
		c.BadRequest("user_id query parameter is required")
		return
	}

	var user models.User
	if err := c.DB().GetForPublicID(c.AppID(), userID, &user); err != nil {
		c.NotFound("User not found")
		return
	}

	query := []any{"app_id = ? AND user_id = ?", c.AppID(), user.ID}

	// Filter by is_metered if provided
	if isMeteredParam := c.Query("is_metered"); isMeteredParam != "" {
		isMetered := isMeteredParam == "true"
		query = []any{"app_id = ? AND user_id = ? AND is_metered = ?", c.AppID(), user.ID, isMetered}
	}

	var entitlements []models.Entitlement
	if err := c.DB().List(&entitlements, query...); err != nil {
		c.ServerError("Failed to list entitlements", err)
		return
	}

	response := make([]apix.Entitlement, len(entitlements))
	for i, entitlement := range entitlements {
		response[i].Set(&entitlement)
	}

	c.OK(response)
}
