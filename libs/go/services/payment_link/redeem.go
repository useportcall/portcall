package payment_link

import (
	"errors"
	"strings"
	"time"

	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	checkout_session "github.com/useportcall/portcall/libs/go/services/checkout_session"
)

func (s *service) Redeem(input *RedeemInput) (*RedeemResult, error) {
	if input == nil {
		return nil, NewValidationError("invalid or expired payment link")
	}
	id := strings.TrimSpace(input.ID)
	now := time.Now().UTC()
	if !cryptox.IsValidPaymentLinkID(id) {
		return nil, NewValidationError("invalid or expired payment link")
	}
	if token := strings.TrimSpace(input.Token); token != "" &&
		cryptox.VerifyPaymentLinkToken(s.crypto, token, id, now) != nil {
		return nil, NewValidationError("invalid or expired payment link")
	}
	link, err := s.loadActiveLink(id, now)
	if err != nil {
		return nil, err
	}
	var plan models.Plan
	if err := s.db.FindFirst(&plan, "app_id = ? AND id = ?", link.AppID, link.PlanID); err != nil {
		return nil, err
	}
	var user models.User
	if err := s.db.FindFirst(&user, "app_id = ? AND id = ?", link.AppID, link.UserID); err != nil {
		return nil, err
	}
	result, err := checkout_session.NewService(s.db, s.crypto).Create(&checkout_session.CreateInput{
		AppID:                 link.AppID,
		PlanID:                plan.PublicID,
		UserID:                user.PublicID,
		CancelURL:             valueOrEmpty(link.CancelURL),
		RedirectURL:           valueOrEmpty(link.RedirectURL),
		RequireBillingAddress: link.RequireBillingAddress,
	})
	if err != nil {
		var ve *checkout_session.ValidationError
		if errors.As(err, &ve) {
			return nil, NewValidationError("%s", ve.Message)
		}
		return nil, err
	}
	link.Plan = plan
	link.User = user
	return &RedeemResult{PaymentLink: link, Session: result.Session, CheckoutURL: result.CheckoutURL}, nil
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
