package main

import (
	"fmt"
	"log"

	"github.com/useportcall/portcall/apps/email/handlers"
	"github.com/useportcall/portcall/libs/go/emailx"
	"github.com/useportcall/portcall/libs/go/qx"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

func startEmailWorker() (server.IServer, error) {
	s, err := server.NewNoDeps(map[string]int{"email_queue": 1})
	if err != nil {
		return nil, fmt.Errorf("new worker server: %w", err)
	}
	emailClient, err := emailx.New()
	if err != nil {
		return nil, fmt.Errorf("new email client: %w", err)
	}
	s.SetEmailClient(emailClient)
	s.H("send_invoice_paid_email", handlers.SendInvoicePaidEmail)
	s.H("send_invoice_dunning_email", handlers.SendInvoiceDunningEmail)
	s.H("send_quote_email", handlers.SendQuoteEmail)
	s.H("send_quote_accepted_confirmation_email", handlers.SendQuoteAcceptedConfirmationEmail)
	s.H("send_quote_status_email", handlers.SendQuoteStatusEmail)
	s.H("process_postmark_webhook_event", handlers.ProcessPostmarkWebhookEvent)
	go func() {
		if err := s.R(); err != nil {
			log.Printf("email e2e worker stopped: %v", err)
		}
	}()
	return s, nil
}

func enqueueInvoicePaid(recipient string) error {
	q, err := qx.New()
	if err != nil {
		return err
	}
	defer q.Close()
	return q.Enqueue("send_invoice_paid_email", map[string]any{
		"recipient_email": recipient,
		"invoice_number":  "INV-E2E-001",
		"company_name":    "Portcall",
		"amount_paid":     "$25.00",
		"date_paid":       "February 10, 2026",
	}, "email_queue")
}

func enqueueStatus(recipient, subject string) error {
	q, err := qx.New()
	if err != nil {
		return err
	}
	defer q.Close()
	return q.Enqueue("send_quote_status_email", map[string]any{
		"recipient_email": recipient,
		"subject":         subject,
		"title":           "Local E2E",
		"message":         "Email worker e2e test via Resend",
		"year":            "2026",
	}, "email_queue")
}
