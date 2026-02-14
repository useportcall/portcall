package payment

import (
	"fmt"
	"strconv"

	"github.com/stripe/stripe-go"
)

var hardDeclineCodes = map[string]bool{
	"fraudulent":                       true,
	"lost_card":                        true,
	"pickup_card":                      true,
	"restricted_card":                  true,
	"revocation_of_all_authorizations": true,
	"security_violation":               true,
	"stolen_card":                      true,
}

func parseStripeInvoiceID(metadata map[string]string) uint {
	raw := metadata[stripeInvoiceMetadataKey]
	parsed, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || parsed == 0 {
		return 0
	}
	return uint(parsed)
}

func formatPaymentIntentFailureReason(data stripe.PaymentIntent, eventType string) string {
	if data.LastPaymentError == nil {
		return fmt.Sprintf("stripe %s", eventType)
	}
	decline := string(data.LastPaymentError.DeclineCode)
	code := string(data.LastPaymentError.Code)
	msg := data.LastPaymentError.Msg
	return joinFailureParts(eventType, decline, code, msg)
}

func formatChargeFailureReason(data stripe.Charge) string {
	return joinFailureParts("charge.failed", data.FailureCode, "", data.FailureMessage)
}

func joinFailureParts(eventType, code, secondaryCode, msg string) string {
	if msg != "" && code != "" {
		return fmt.Sprintf("stripe %s: %s (%s)", eventType, msg, code)
	}
	if msg != "" {
		return fmt.Sprintf("stripe %s: %s", eventType, msg)
	}
	if code != "" {
		return fmt.Sprintf("stripe %s: %s", eventType, code)
	}
	if secondaryCode != "" {
		return fmt.Sprintf("stripe %s: %s", eventType, secondaryCode)
	}
	return fmt.Sprintf("stripe %s", eventType)
}

func isHardDeclineCode(code string) bool {
	return hardDeclineCodes[code]
}
