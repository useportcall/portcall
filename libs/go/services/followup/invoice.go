package followup

const defaultInvoiceMaxAttempts = 4

const (
	StageInvoiceFirst = "invoice_first_reminder"
	StageInvoiceRetry = "invoice_retry_reminder"
	StageInvoiceFinal = "invoice_final_notice"
)

type InvoiceFailureInput struct {
	Attempt     int
	MaxAttempts int
	NoRetry     bool
}

type InvoiceFailureDecision struct {
	Attempt       int
	MaxAttempts   int
	FinalAttempt  bool
	Stage         string
	PaymentStatus string
}

func DecideInvoiceFailure(input InvoiceFailureInput) InvoiceFailureDecision {
	maxAttempts := normalizePositive(input.MaxAttempts, defaultInvoiceMaxAttempts)
	attempt := normalizePositive(input.Attempt, 1)
	if input.NoRetry && attempt < maxAttempts {
		attempt = maxAttempts
	}
	if attempt > maxAttempts {
		attempt = maxAttempts
	}

	decision := InvoiceFailureDecision{
		Attempt:       attempt,
		MaxAttempts:   maxAttempts,
		FinalAttempt:  attempt >= maxAttempts,
		Stage:         StageInvoiceRetry,
		PaymentStatus: "past_due",
	}
	if decision.FinalAttempt {
		decision.Stage = StageInvoiceFinal
		decision.PaymentStatus = "uncollectible"
		return decision
	}
	if attempt == 1 {
		decision.Stage = StageInvoiceFirst
	}
	return decision
}

func normalizePositive(value, fallback int) int {
	if value > 0 {
		return value
	}
	return fallback
}
