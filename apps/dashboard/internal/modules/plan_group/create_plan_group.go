package plan_group

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type CreatePlanGroupRequest struct {
	Name   string `json:"name"`
	PlanID string `json:"plan_id"`
}

func CreatePlanGroup(c *routerx.Context) {
	var body CreatePlanGroupRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	planGroup := &models.PlanGroup{
		PublicID: dbx.GenPublicID("pg"),
		AppID:    c.AppID(),
		Name:     body.Name}
	if err := c.DB().Create(planGroup); err != nil {
		c.NotFound("Failed to create plan group")
		return
	}

	if body.PlanID != "" {
		var plan models.Plan
		if err := c.DB().GetForPublicID(c.AppID(), body.PlanID, &plan); err != nil {
			c.NotFound("Plan not found")
			return
		}

		plan.PlanGroupID = &planGroup.ID
		if err := c.DB().Save(&plan); err != nil {
			c.NotFound("Failed to update plan with plan group")
			return
		}
	}

	c.OK(new(apix.PlanGroup).Set(planGroup))
}
