package checkout_session

import (
	"log"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// Resolve handles checkout session resolution.
// It finds the session by external ID, checks it is still pending resolution,
// and updates the status to "resolved".
func (s *service) Resolve(payload *ResolvePayload) (*ResolveResult, error) {
	var session models.CheckoutSession
	if err := s.db.FindFirst(&session, "external_session_id = ?", payload.ExternalSessionID); err != nil {
		return nil, err
	}

	if session.Status != "active" && session.Status != "pending" {
		log.Printf("[Resolve] Checkout session %s is not active/pending (current status: %s), skipping",
			payload.ExternalSessionID, session.Status)
		return &ResolveResult{Skipped: true}, nil
	}

	session.Status = "resolved"
	if err := s.db.Save(&session); err != nil {
		return nil, err
	}

	return &ResolveResult{
		Session:                 &session,
		ExternalPaymentMethodID: payload.ExternalPaymentMethodID,
	}, nil
}
