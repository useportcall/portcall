package payment

import (
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
)

// Service is the payment service interface.
// Mock this interface in saga handler tests.
type Service interface {
	Pay(input *PayInput) (*PayResult, error)
	Resolve(input *ResolveInput) (*ResolveResult, error)
	ProcessStripeWebhook(input *StripeWebhookInput) (*StripeResult, error)
	ProcessBraintreeWebhook(input *BraintreeWebhookInput) (*BraintreeResult, error)
	ProcessDunning(input *DunningInput) (*DunningResult, error)
	CreatePaymentMethod(input *CreateMethodInput) (*CreateMethodResult, error)
}

type service struct {
	db     dbx.IORM
	crypto cryptox.ICrypto
}

// NewService creates a new payment service.
func NewService(db dbx.IORM, crypto cryptox.ICrypto) Service {
	return &service{db: db, crypto: crypto}
}
