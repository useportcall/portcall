package plan

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func DeletePlan(c *routerx.Context) {
	id := c.Param("id")
	plan := &models.Plan{}
	if err := c.DB().GetForPublicID(c.AppID(), id, plan); err != nil {
		c.NotFound("Plan not found")
		return
	}

	if plan.Status != "draft" {
		c.BadRequest("Cannot delete plan with status '" + plan.Status + "'")
		return
	}

	if err := c.DB().DeleteForID(plan); err != nil {
		c.ServerError("Failed to delete plan")
		return
	}

	c.OK(map[string]any{"deleted": true, "id": id})
}
