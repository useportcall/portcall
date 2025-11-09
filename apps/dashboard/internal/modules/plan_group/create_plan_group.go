package plan_group

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

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

	c.OK(new(PlanGroup).Set(planGroup))
}
