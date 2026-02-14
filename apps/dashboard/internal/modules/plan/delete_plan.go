package plan

import (
	quotemodule "github.com/useportcall/portcall/apps/dashboard/internal/modules/quote"
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
	locked, err := quotemodule.HasLockedQuoteForPlan(c, plan.ID)
	if err != nil {
		c.ServerError("Failed to validate quote state", err)
		return
	}
	if locked {
		c.BadRequest("Plan cannot be deleted after quote is issued")
		return
	}

	if plan.Status != "draft" {
		c.BadRequest("Cannot delete plan with status '" + plan.Status + "'")
		return
	}

	if err := c.DB().DeleteForID(plan); err != nil {
		c.ServerError("Failed to delete plan", err)
		return
	}

	c.OK(map[string]any{"deleted": true, "id": id})
}
