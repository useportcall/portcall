package entitlement

import (
	"github.com/useportcall/portcall/libs/go/dbx"
)

// Service is the entitlement service interface.
// Mock this interface in saga handler tests.
type Service interface {
	ResetAll(input *ResetAllInput) (*ResetAllResult, error)
	StartUpsert(input *StartUpsertInput) (*StartUpsertResult, error)
	Upsert(input *UpsertInput) (*UpsertResult, error)
	IncrementUsage(input *IncrementUsageInput) (*IncrementUsageResult, error)
}

type service struct {
	db dbx.IORM
}

// NewService creates a new entitlement service.
func NewService(db dbx.IORM) Service {
	return &service{db: db}
}
