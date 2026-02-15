package main

import (
	"github.com/useportcall/portcall/apps/email/handlers"
	"github.com/useportcall/portcall/libs/go/emailx"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

func runEmailWorker() error {
	s, err := server.NewNoDeps(map[string]int{"email_queue": 10})
	if err != nil {
		return err
	}
	emailClient, err := emailx.New()
	if err != nil {
		return err
	}
	s.SetEmailClient(emailClient)
	s.H("send_invoice_paid_email", handlers.SendInvoicePaidEmail)
	s.H("send_invoice_dunning_email", handlers.SendInvoiceDunningEmail)
	s.H("send_quote_email", handlers.SendQuoteEmail)
	s.H("send_quote_accepted_confirmation_email", handlers.SendQuoteAcceptedConfirmationEmail)
	s.H("send_quote_status_email", handlers.SendQuoteStatusEmail)
	s.H("process_postmark_webhook_event", handlers.ProcessPostmarkWebhookEvent)
	return s.R()
}
