package on_payment

import (
	"encoding/json"

	"github.com/useportcall/portcall/libs/go/qx/server"
	"github.com/useportcall/portcall/libs/go/services/payment"
)

func payHandler(c server.IContext) error {
	var input payment.PayInput
	if err := json.Unmarshal(c.Payload(), &input); err != nil {
		return err
	}

	svc := payment.NewService(c.DB(), c.Crypto())
	result, err := svc.Pay(&input)
	if err != nil {
		return processPayFailure(c, svc, &input, err)
	}
	return ResolveInvoice.Enqueue(c.Queue(), map[string]any{
		"invoice_id": result.Invoice.ID,
	})
}

func resolveHandler(c server.IContext) error {
	var input payment.ResolveInput
	if err := json.Unmarshal(c.Payload(), &input); err != nil {
		return err
	}

	svc := payment.NewService(c.DB(), c.Crypto())
	result, err := svc.Resolve(&input)
	if err != nil {
		return err
	}
	if result.EmailPayload != nil {
		return sendInvoiceEmail.Enqueue(c.Queue(), result.EmailPayload)
	}
	return nil
}
