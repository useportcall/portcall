package payment_link

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func (s *service) loadPublishedPlan(appID uint, planID string) (*models.Plan, error) {
	var plan models.Plan
	if err := s.db.GetForPublicID(appID, planID, &plan); err != nil {
		return nil, err
	}
	if plan.Status != "published" {
		return nil, NewValidationError("plan with id '%s' is not yet published", plan.PublicID)
	}
	return &plan, nil
}

func (s *service) loadActiveLink(id string, now time.Time) (*models.PaymentLink, error) {
	var link models.PaymentLink
	if err := s.db.FindFirst(&link, "public_id = ?", id); err != nil {
		return nil, NewValidationError("invalid or expired payment link")
	}
	if link.Status != "active" || now.After(link.ExpiresAt) {
		return nil, NewValidationError("invalid or expired payment link")
	}
	return &link, nil
}
