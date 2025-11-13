package plan_item

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func DeletePlanItem(c *routerx.Context) {
	id := c.Param("id")
	planItem := &models.PlanItem{}
	if err := c.DB().GetForPublicID(c.AppID(), id, planItem); err != nil {
		c.NotFound("Plan item not found")
		return
	}

	if err := c.DB().DeleteForID(planItem); err != nil {
		c.ServerError("Failed to delete plan item", err)
		return
	}

	c.OK(map[string]any{"deleted": true, "id": id})
}
