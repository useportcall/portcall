package payment_link

import (
	"os"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func (s *service) Create(input *CreateInput) (*CreateResult, error) {
	if input == nil {
		return nil, NewValidationError("missing create input")
	}
	now := time.Now().UTC()
	expiresAt, err := resolveLinkExpiry(input.ExpiresAt, now)
	if err != nil {
		return nil, err
	}
	redirectURL, err := validateReturnURL(input.RedirectURL, "redirect_url")
	if err != nil {
		return nil, err
	}
	cancelURL, err := validateReturnURL(input.CancelURL, "cancel_url")
	if err != nil {
		return nil, err
	}
	plan, err := s.loadPublishedPlan(input.AppID, input.PlanID)
	if err != nil {
		return nil, err
	}
	user, err := s.resolveUser(input)
	if err != nil {
		return nil, err
	}
	if err := s.ensurePaymentConnection(input.AppID); err != nil {
		return nil, err
	}
	link := &models.PaymentLink{
		PublicID:              dbx.GenPublicID("pl"),
		AppID:                 input.AppID,
		PlanID:                plan.ID,
		UserID:                user.ID,
		ExpiresAt:             expiresAt,
		Status:                "active",
		RedirectURL:           stringPointer(redirectURL),
		CancelURL:             stringPointer(cancelURL),
		RequireBillingAddress: input.RequireBillingAddress,
	}
	if err := s.db.Create(link); err != nil {
		return nil, err
	}
	link.Plan = *plan
	link.User = *user
	linkURL, err := buildPaymentLinkURL(os.Getenv("CHECKOUT_URL"), link, s.crypto)
	if err != nil {
		return nil, err
	}
	return &CreateResult{PaymentLink: link, URL: linkURL}, nil
}

func stringPointer(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
