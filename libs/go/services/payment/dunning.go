package payment

import (
	"fmt"

	"github.com/useportcall/portcall/libs/go/services/followup"
)

// ProcessDunning updates billing state for a failed payment attempt
// and prepares an email payload for customer notifications.
func (s *service) ProcessDunning(input *DunningInput) (*DunningResult, error) {
	invoice, err := findInvoice(s.db, input.InvoiceID)
	if err != nil {
		return nil, err
	}
	if invoice.Status == "paid" {
		return &DunningResult{}, nil
	}
	alreadyPastDue := invoice.Status == "past_due"

	decision := followup.DecideInvoiceFailure(followup.InvoiceFailureInput{
		Attempt:     input.Attempt,
		MaxAttempts: input.MaxAttempts,
		NoRetry:     input.NoRetry,
	})
	finalAttempt := decision.FinalAttempt
	if err := s.updateDunningInvoiceStatus(invoice, finalAttempt); err != nil {
		return nil, err
	}
	if err := s.updateDunningSubscriptionStatus(invoice, finalAttempt); err != nil {
		return nil, err
	}

	if invoice.Total < 1 {
		return &DunningResult{FinalAttempt: finalAttempt}, nil
	}

	user, company, err := lookupResolveDeps(s.db, invoice.UserID, invoice.AppID)
	if err != nil {
		return nil, err
	}
	if alreadyPastDue && !finalAttempt && decision.Attempt <= 1 {
		return &DunningResult{FinalAttempt: false}, nil
	}

	amountDue := float64(invoice.Total) / 100.0
	return &DunningResult{
		FinalAttempt: finalAttempt,
		EmailPayload: &DunningEmailPayload{
			InvoiceNumber:  invoice.InvoiceNumber,
			AmountDue:      fmt.Sprintf("$%.2f", amountDue),
			DueDate:        invoice.DueBy.Format("January 2, 2006"),
			CompanyName:    company.Name,
			RecipientEmail: user.Email,
			Attempt:        decision.Attempt,
			MaxAttempts:    decision.MaxAttempts,
			FinalAttempt:   finalAttempt,
			FollowUpStage:  decision.Stage,
			FailureReason:  input.FailureReason,
			LogoURL:        emailLogoURL(company),
			PaymentStatus:  decision.PaymentStatus,
		},
	}, nil
}
