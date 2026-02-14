package payment

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func (s *service) updateDunningInvoiceStatus(invoice *models.Invoice, finalAttempt bool) error {
	status := "past_due"
	if finalAttempt {
		status = "uncollectible"
	}
	if invoice.Status == status {
		return nil
	}
	invoice.Status = status
	return s.db.Save(invoice)
}

func (s *service) updateDunningSubscriptionStatus(invoice *models.Invoice, finalAttempt bool) error {
	if invoice.SubscriptionID == nil {
		return nil
	}

	var sub models.Subscription
	if err := s.db.FindForID(*invoice.SubscriptionID, &sub); err != nil {
		return err
	}

	nextStatus := "past_due"
	if finalAttempt {
		nextStatus = "canceled"
	}
	if sub.Status == nextStatus {
		return nil
	}
	sub.Status = nextStatus
	return s.db.Save(&sub)
}
