package subscription

import (
	"fmt"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/services/invoice"
	"github.com/useportcall/portcall/libs/go/services/payment"
)

func createStandardInvoice(db dbx.IORM, subscriptionID uint, itemCount int) (uint, error) {
	if itemCount == 0 {
		return 0, nil
	}
	svc := invoice.NewService(db)
	result, err := svc.Create(&invoice.CreateInput{SubscriptionID: subscriptionID})
	if err != nil || result.Skipped {
		return 0, err
	}
	return result.Invoice.ID, nil
}

func payAndResolveInvoice(input *BillingFlowInput, invoiceID uint) error {
	if input.AsyncPayment {
		if input.Queue == nil {
			return fmt.Errorf("async payment requires queue")
		}
		return input.Queue.Enqueue(
			"pay_invoice",
			map[string]any{"invoice_id": invoiceID},
			"billing_queue",
		)
	}

	svc := payment.NewService(input.DB, input.Crypto)
	if _, err := svc.Pay(&payment.PayInput{InvoiceID: invoiceID}); err != nil {
		return processPayFailure(input, svc, invoiceID, err)
	}

	result, err := svc.Resolve(&payment.ResolveInput{InvoiceID: invoiceID})
	if err != nil || result.EmailPayload == nil || input.Queue == nil {
		return err
	}
	return input.Queue.Enqueue("send_invoice_paid_email", result.EmailPayload, "email_queue")
}
