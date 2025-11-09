package plan_group

import "github.com/useportcall/portcall/libs/go/dbx/models"

type PlanGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (pg *PlanGroup) Set(planGroup *models.PlanGroup) *PlanGroup {
	pg.ID = planGroup.PublicID
	pg.Name = planGroup.Name
	return pg
}

type CreatePlanGroupRequest struct {
	Name string `json:"name"`
}
