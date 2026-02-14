package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"html/template"

	"github.com/useportcall/portcall/libs/go/qx/server"
)

type invoiceDunningPayload struct {
	InvoiceNumber  string `json:"invoice_number"`
	AmountDue      string `json:"amount_due"`
	DueDate        string `json:"due_date"`
	CompanyName    string `json:"company_name"`
	RecipientEmail string `json:"recipient_email"`
	Attempt        int    `json:"attempt"`
	MaxAttempts    int    `json:"max_attempts"`
	FinalAttempt   bool   `json:"final_attempt"`
	FailureReason  string `json:"failure_reason"`
	LogoURL        string `json:"logo_url"`
	PaymentStatus  string `json:"payment_status"`
}

func SendInvoiceDunningEmail(c server.IContext) error {
	var payload invoiceDunningPayload
	if err := json.Unmarshal(c.Payload(), &payload); err != nil {
		return err
	}
	if payload.RecipientEmail == "" {
		return fmt.Errorf("recipient_email not found in payload")
	}
	if payload.PaymentStatus == "" {
		payload.PaymentStatus = "past_due"
	}

	tmpl, err := template.ParseFiles("templates/invoice_dunning_notice.html")
	if err != nil {
		return err
	}
	var htmlContentBuf bytes.Buffer
	if err := tmpl.Execute(&htmlContentBuf, payload); err != nil {
		return err
	}

	subject := fmt.Sprintf("Payment failed (attempt %d of %d)", payload.Attempt, payload.MaxAttempts)
	if payload.FinalAttempt {
		subject = "Final payment attempt failed"
	}
	from := os.Getenv("EMAIL_FROM")
	if from == "" {
		return fmt.Errorf("EMAIL_FROM environment variable is not set")
	}
	return c.EmailClient().Send(
		htmlContentBuf.String(),
		subject,
		from,
		[]string{payload.RecipientEmail},
	)
}
