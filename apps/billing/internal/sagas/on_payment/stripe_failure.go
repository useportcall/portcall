package on_payment

import (
	"encoding/json"

	"github.com/useportcall/portcall/libs/go/qx/server"
	"github.com/useportcall/portcall/libs/go/services/payment"
)

func stripePaymentFailureHandler(c server.IContext) error {
	var input payment.StripeFailurePayload
	if err := json.Unmarshal(c.Payload(), &input); err != nil {
		return err
	}

	svc := payment.NewService(c.DB(), c.Crypto())
	maxAttempts := dunningMaxAttempts(c)
	result, err := svc.ProcessDunning(&payment.DunningInput{
		InvoiceID:     input.InvoiceID,
		Attempt:       input.Attempt,
		MaxAttempts:   maxAttempts,
		NoRetry:       input.NoRetry,
		FailureReason: input.FailureReason,
	})
	if err != nil || result.EmailPayload == nil {
		return err
	}
	return sendInvoiceDunningEmail.Enqueue(c.Queue(), result.EmailPayload)
}
