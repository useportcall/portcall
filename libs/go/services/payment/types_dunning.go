package payment

// DunningInput is the input for processing a failed payment attempt.
type DunningInput struct {
	InvoiceID     uint
	Attempt       int
	MaxAttempts   int
	NoRetry       bool
	FailureReason string
}

// DunningEmailPayload is the payload for payment-failure notification emails.
type DunningEmailPayload struct {
	InvoiceNumber  string `json:"invoice_number"`
	AmountDue      string `json:"amount_due"`
	DueDate        string `json:"due_date"`
	CompanyName    string `json:"company_name"`
	RecipientEmail string `json:"recipient_email"`
	Attempt        int    `json:"attempt"`
	MaxAttempts    int    `json:"max_attempts"`
	FinalAttempt   bool   `json:"final_attempt"`
	FollowUpStage  string `json:"followup_stage"`
	FailureReason  string `json:"failure_reason"`
	LogoURL        string `json:"logo_url"`
	PaymentStatus  string `json:"payment_status"`
}

// DunningResult is the result of processing a failed payment attempt.
type DunningResult struct {
	EmailPayload *DunningEmailPayload
	FinalAttempt bool
}
