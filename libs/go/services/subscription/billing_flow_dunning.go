package subscription

import (
	"os"
	"strconv"

	"github.com/useportcall/portcall/libs/go/services/payment"
)

const defaultDunningMaxAttempts = 4

func processPayFailure(input *BillingFlowInput, svc payment.Service, invoiceID uint, payErr error) error {
	result, err := svc.ProcessDunning(&payment.DunningInput{
		InvoiceID:     invoiceID,
		Attempt:       1,
		MaxAttempts:   dunningMaxAttempts(),
		FailureReason: payErr.Error(),
	})
	if err != nil {
		return err
	}
	if input.Queue == nil || result.EmailPayload == nil {
		return payErr
	}
	if err := input.Queue.Enqueue("send_invoice_dunning_email", result.EmailPayload, "email_queue"); err != nil {
		return err
	}
	return payErr
}

func dunningMaxAttempts() int {
	maxAttempts := defaultDunningMaxAttempts
	raw := os.Getenv("BILLING_DUNNING_MAX_ATTEMPTS")
	if raw == "" {
		return maxAttempts
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed <= 0 {
		return maxAttempts
	}
	return parsed
}
