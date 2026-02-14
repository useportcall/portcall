package payment

import (
	"log"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// CreatePaymentMethod creates or updates a payment method.
// Single mutation: upserts the payment method.
func (s *service) CreatePaymentMethod(input *CreateMethodInput) (*CreateMethodResult, error) {
	log.Printf("Processing CreatePaymentMethod for user ID %d", input.UserID)

	pm := &models.PaymentMethod{
		PublicID:     dbx.GenPublicID("pm"),
		AppID:        input.AppID,
		UserID:       input.UserID,
		ExternalID:   input.ExternalPaymentMethodID,
		ExternalType: "card",
	}

	var existing models.PaymentMethod
	if err := s.db.FindFirst(&existing, "user_id = ? AND external_id = ?", pm.UserID, pm.ExternalID); err != nil {
		if !dbx.IsRecordNotFoundError(err) {
			return nil, err
		}
		// Not found — create new
		if err := s.db.Create(pm); err != nil {
			return nil, err
		}
		log.Printf("Payment method created with ID %d for user ID %d", pm.ID, pm.UserID)
	}

	return &CreateMethodResult{
		PaymentMethod: pm,
		AppID:         input.AppID,
		UserID:        input.UserID,
		PlanID:        input.PlanID,
	}, nil
}
