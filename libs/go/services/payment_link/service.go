package payment_link

import (
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
)

type Service interface {
	Create(input *CreateInput) (*CreateResult, error)
	Redeem(input *RedeemInput) (*RedeemResult, error)
}

type service struct {
	db     dbx.IORM
	crypto cryptox.ICrypto
}

func NewService(db dbx.IORM, crypto cryptox.ICrypto) Service {
	return &service{db: db, crypto: crypto}
}
