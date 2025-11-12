package plan_feature

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func DeletePlanFeature(c *routerx.Context) {
	id := c.Param("id")
	planFeature := &models.PlanFeature{}
	if err := c.DB().GetForPublicID(c.AppID(), id, planFeature); err != nil {
		c.NotFound("Plan feature not found")
		return
	}

	if err := c.DB().DeleteForID(planFeature); err != nil {
		c.ServerError("Failed to delete plan feature")
		return
	}

	c.OK(map[string]any{"deleted": true, "id": id})
}
