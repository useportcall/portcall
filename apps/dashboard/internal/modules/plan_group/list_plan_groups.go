package plan_group

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ListPlanGroups(c *routerx.Context) {
	planGroups := []models.PlanGroup{}
	if err := c.DB().List(&planGroups, "app_id = ?", c.AppID()); err != nil {
		c.NotFound("Plan groups not found")
		return
	}

	response := make([]PlanGroup, len(planGroups))
	for i, pg := range planGroups {
		response[i].Set(&pg)
	}

	c.OK(response)
}
