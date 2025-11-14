package main

import (
	"github.com/useportcall/portcall/apps/email/handlers"
	"github.com/useportcall/portcall/libs/go/envx"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

func main() {
	envx.Load()

	server := server.NewNoDeps(map[string]int{"email_queue": 10})

	server.H("send_invoice_paid_email", handlers.SendInvoicePaidEmail)
}
