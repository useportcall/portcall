package subscription

import "github.com/useportcall/portcall/libs/go/dbx"

// Service is the subscription service interface.
type Service interface {
	Find(input *FindInput) (*FindResult, error)
	Create(input *CreateInput) (*CreateResult, error)
	Update(input *UpdateInput) (*UpdateResult, error)
	PlanSwitch(input *PlanSwitchInput) (*PlanSwitchResult, error)
	CreateItems(input *CreateItemsInput) (*CreateItemsResult, error)
	FindResets() (*FindResetsResult, error)
	StartReset(input *StartResetInput) (*StartResetResult, error)
	EndReset(input *EndResetInput) (*EndResetResult, error)
}

type service struct {
	db dbx.IORM
}

// NewService creates a new subscription service.
func NewService(db dbx.IORM) Service {
	return &service{db: db}
}