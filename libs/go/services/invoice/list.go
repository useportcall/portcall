package invoice

import (
	"fmt"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// List returns invoices for an app, optionally filtered by subscription or user.
func (s *service) List(input *ListInput) (*ListResult, error) {
	conds, err := s.buildListConditions(input)
	if err != nil {
		return nil, err
	}

	var invoices []models.Invoice
	if err := s.db.ListWithOrder(&invoices, "created_at DESC", conds...); err != nil {
		return nil, fmt.Errorf("listing invoices: %w", err)
	}

	return &ListResult{Invoices: invoices}, nil
}

func (s *service) buildListConditions(input *ListInput) ([]any, error) {
	if input.SubscriptionID != "" {
		var sub models.Subscription
		if err := s.db.GetForPublicID(input.AppID, input.SubscriptionID, &sub); err != nil {
			return nil, fmt.Errorf("subscription not found")
		}
		return []any{"app_id = ? AND subscription_id = ?", input.AppID, sub.ID}, nil
	}

	if input.UserID != "" {
		var user models.User
		if err := s.db.GetForPublicID(input.AppID, input.UserID, &user); err != nil {
			return nil, fmt.Errorf("user not found")
		}
		return []any{"app_id = ? AND user_id = ?", input.AppID, user.ID}, nil
	}

	return []any{"app_id = ?", input.AppID}, nil
}
