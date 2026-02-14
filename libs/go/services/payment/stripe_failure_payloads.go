package payment

import (
	"fmt"

	"github.com/stripe/stripe-go"
)

func parsePaymentIntentFailure(event stripe.Event) (*StripeFailurePayload, error) {
	var data stripe.PaymentIntent
	if err := data.UnmarshalJSON(event.Data.Raw); err != nil {
		return nil, err
	}
	invoiceID := parseStripeInvoiceID(data.Metadata)
	if invoiceID == 0 {
		return nil, nil
	}
	declineCode := ""
	if data.LastPaymentError != nil {
		declineCode = string(data.LastPaymentError.DeclineCode)
	}
	return &StripeFailurePayload{
		InvoiceID:     invoiceID,
		Attempt:       1,
		NoRetry:       isHardDeclineCode(declineCode),
		EventType:     string(event.Type),
		FailureReason: formatPaymentIntentFailureReason(data, string(event.Type)),
	}, nil
}

func parseChargeFailure(event stripe.Event) (*StripeFailurePayload, error) {
	var data stripe.Charge
	if err := data.UnmarshalJSON(event.Data.Raw); err != nil {
		return nil, err
	}
	invoiceID := parseStripeInvoiceID(data.Metadata)
	if invoiceID == 0 {
		return nil, nil
	}
	return &StripeFailurePayload{
		InvoiceID:     invoiceID,
		Attempt:       1,
		NoRetry:       isHardDeclineCode(data.FailureCode),
		EventType:     string(event.Type),
		FailureReason: formatChargeFailureReason(data),
	}, nil
}

func parseInvoiceFailure(event stripe.Event) (*StripeFailurePayload, error) {
	var data stripe.Invoice
	if err := data.UnmarshalJSON(event.Data.Raw); err != nil {
		return nil, err
	}
	invoiceID := parseStripeInvoiceID(data.Metadata)
	if invoiceID == 0 {
		return nil, nil
	}
	attempt := int(data.AttemptCount)
	if attempt < 1 {
		attempt = 1
	}
	return &StripeFailurePayload{
		InvoiceID: invoiceID,
		Attempt:   attempt,
		EventType: string(event.Type),
		FailureReason: fmt.Sprintf(
			"stripe %s (attempt %d)",
			event.Type,
			attempt,
		),
	}, nil
}
