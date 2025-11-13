package entitlement

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func GetEntitlement(c *routerx.Context) {
	userID := c.Param("user_id")
	entitlementID := c.Param("id")

	query := "app_id = ? AND user_id = ? AND feature_public_id = ?"

	var user models.User
	if err := c.DB().GetForPublicID(c.AppID(), userID, &user); err != nil {
		c.NotFound("User not found")
		return
	}

	var entitlement models.Entitlement
	if err := c.DB().FindFirst(&entitlement, query, c.AppID(), user.ID, entitlementID); err != nil {
		if !dbx.IsRecordNotFoundError(err) {
			c.ServerError("Internal server error", err)
			return
		}
	}

	c.OK(new(apix.Entitlement).Set(&entitlement))
}
