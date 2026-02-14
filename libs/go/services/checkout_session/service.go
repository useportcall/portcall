package checkout_session

import (
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
)

// Service is the checkout session service interface.
// Mock this interface in handler tests to avoid real DB/payment calls.
type Service interface {
	Create(input *CreateInput) (*CreateResult, error)
	Resolve(payload *ResolvePayload) (*ResolveResult, error)
}

type service struct {
	db     dbx.IORM
	crypto cryptox.ICrypto
}

// NewService creates a new checkout session service.
func NewService(db dbx.IORM, crypto cryptox.ICrypto) Service {
	return &service{db: db, crypto: crypto}
}
