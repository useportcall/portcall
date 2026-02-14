package plan_item

import (
	quotemodule "github.com/useportcall/portcall/apps/dashboard/internal/modules/quote"
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
	locked, err := quotemodule.HasLockedQuoteForPlan(c, planItem.PlanID)
	if err != nil {
		c.ServerError("Failed to validate quote state", err)
		return
	}
	if locked {
		c.BadRequest("Plan item cannot be deleted after quote is issued")
		return
	}

	if err := c.DB().DeleteForID(planItem); err != nil {
		c.ServerError("Failed to delete plan item", err)
		return
	}

	c.OK(map[string]any{"deleted": true, "id": id})
}
