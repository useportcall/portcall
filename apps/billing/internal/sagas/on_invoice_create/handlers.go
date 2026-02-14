package on_invoice_create

import (
	"encoding/json"

	"github.com/useportcall/portcall/apps/billing/internal/sagas/on_payment"
	"github.com/useportcall/portcall/libs/go/qx/server"
	"github.com/useportcall/portcall/libs/go/services/invoice"
)

func createInvoiceHandler(c server.IContext) error {
	var payload struct {
		SubscriptionID uint `json:"subscription_id"`
	}
	if err := json.Unmarshal(c.Payload(), &payload); err != nil {
		return err
	}

	svc := invoice.NewService(c.DB())
	result, err := svc.Create(&invoice.CreateInput{
		SubscriptionID: payload.SubscriptionID,
	})
	if err != nil || result.Skipped {
		return err
	}

	return on_payment.PayInvoice.Enqueue(c.Queue(), map[string]any{
		"invoice_id": result.Invoice.ID,
	})
}
