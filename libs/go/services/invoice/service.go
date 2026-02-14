package invoice

import (
	"github.com/useportcall/portcall/libs/go/dbx"
)

// Service is the invoice service interface.
// Mock this interface in saga handler tests.
type Service interface {
	List(input *ListInput) (*ListResult, error)
	Create(input *CreateInput) (*CreateResult, error)
	CreateUpgrade(input *CreateUpgradeInput) (*CreateUpgradeResult, error)
}

type service struct {
	db dbx.IORM
}

// NewService creates a new invoice service.
func NewService(db dbx.IORM) Service {
	return &service{db: db}
}
