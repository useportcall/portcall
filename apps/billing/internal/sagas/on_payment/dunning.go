package on_payment

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hibiken/asynq"
	"github.com/useportcall/portcall/libs/go/qx/server"
	"github.com/useportcall/portcall/libs/go/services/payment"
)

const defaultDunningMaxAttempts = 4

func processPayFailure(c server.IContext, svc payment.Service, input *payment.PayInput, payErr error) error {
	retryCount, _ := asynq.GetRetryCount(c)
	attempt := retryCount + 1
	maxAttempts := dunningMaxAttempts(c)

	result, err := svc.ProcessDunning(&payment.DunningInput{
		InvoiceID:     input.InvoiceID,
		Attempt:       attempt,
		MaxAttempts:   maxAttempts,
		FailureReason: payErr.Error(),
	})
	if err != nil {
		return err
	}
	if result.EmailPayload != nil {
		if err := sendInvoiceDunningEmail.Enqueue(c.Queue(), result.EmailPayload); err != nil {
			return err
		}
	}
	if result.FinalAttempt {
		return fmt.Errorf("%w: %v", asynq.SkipRetry, payErr)
	}
	return payErr
}

func dunningMaxAttempts(c server.IContext) int {
	maxAttempts := defaultDunningMaxAttempts
	if raw := os.Getenv("BILLING_DUNNING_MAX_ATTEMPTS"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			maxAttempts = parsed
		}
	}
	maxRetry, ok := asynq.GetMaxRetry(c)
	if !ok {
		return maxAttempts
	}
	maxRetry += 1
	if maxRetry > 0 && maxAttempts > maxRetry {
		return maxRetry
	}
	return maxAttempts
}
