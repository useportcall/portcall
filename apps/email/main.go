package main

import (
	"log"
	"net/http"
	"os"

	"github.com/useportcall/portcall/apps/email/handlers"
	"github.com/useportcall/portcall/libs/go/emailx"
	"github.com/useportcall/portcall/libs/go/envx"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

func main() {
	envx.Load()

	server, err := server.NewNoDeps(map[string]int{"email_queue": 10})
	if err != nil {
		log.Fatalf("failed to init worker server: %v", err)
	}
	emailClient, err := emailx.New()
	if err != nil {
		log.Fatalf("failed to init email client: %v", err)
	}
	server.SetEmailClient(emailClient)

	server.H("send_invoice_paid_email", handlers.SendInvoicePaidEmail)
	server.H("send_invoice_dunning_email", handlers.SendInvoiceDunningEmail)
	server.H("send_quote_email", handlers.SendQuoteEmail)
	server.H("send_quote_accepted_confirmation_email", handlers.SendQuoteAcceptedConfirmationEmail)
	server.H("send_quote_status_email", handlers.SendQuoteStatusEmail)
	server.H("process_postmark_webhook_event", handlers.ProcessPostmarkWebhookEvent)
	if err := server.R(); err != nil {
		log.Fatalf("email worker failed: %v", err)
	}
}

// start a small health endpoint for orchestration
func init() {
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "9091"
		}
		mux := http.NewServeMux()
		mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})
		mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})
		log.Printf("email health server listening on :%s", port)
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			log.Printf("email health server error: %v", err)
		}
	}()
}
