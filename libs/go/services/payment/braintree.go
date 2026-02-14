package payment

import (
	"fmt"
	"log"
	"strings"
)

const (
	braintreeTxnSettled            = "transaction_settled"
	braintreeTxnSettlementDeclined = "transaction_settlement_declined"
	braintreeSubChargedSuccess     = "subscription_charged_successfully"
	braintreeSubChargedFailed      = "subscription_charged_unsuccessfully"
)

func (s *service) ProcessBraintreeWebhook(input *BraintreeWebhookInput) (*BraintreeResult, error) {
	switch strings.TrimSpace(input.Kind) {
	case braintreeTxnSettled, braintreeSubChargedSuccess:
		sessionID := parseOrderMetadataValue(input.OrderID, "portcall_checkout_session_id")
		if sessionID == "" || input.PaymentMethodToken == "" {
			return &BraintreeResult{Handled: true}, nil
		}
		return &BraintreeResult{
			Action:          "resolve_checkout_session",
			SessionID:       sessionID,
			PaymentMethodID: input.PaymentMethodToken,
			Handled:         true,
		}, nil
	case braintreeTxnSettlementDeclined, braintreeSubChargedFailed:
		invoiceID := parseOrderMetadataUint(input.OrderID, "portcall_invoice_id")
		if invoiceID == 0 {
			return &BraintreeResult{Handled: true}, nil
		}
		attempt := input.FailureCount
		if attempt < 1 {
			attempt = 1
		}
		return &BraintreeResult{
			Action: "process_stripe_payment_failure",
			Failure: &StripeFailurePayload{
				InvoiceID:     invoiceID,
				Attempt:       attempt,
				EventType:     input.Kind,
				FailureReason: formatBraintreeFailureReason(input),
			},
			Handled: true,
		}, nil
	default:
		log.Println("UNHANDLED_BRAINTREE_WEBHOOK_KIND:", input.Kind)
		return &BraintreeResult{Handled: false}, nil
	}
}

func formatBraintreeFailureReason(input *BraintreeWebhookInput) string {
	reason := strings.TrimSpace(input.FailureReason)
	if reason == "" {
		return fmt.Sprintf("braintree %s", input.Kind)
	}
	return fmt.Sprintf("braintree %s: %s", input.Kind, reason)
}
