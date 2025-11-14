package handlers

import (
	"github.com/useportcall/portcall/libs/go/emailx"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

func SendInvoicePaidEmail(c server.IContext) error {
	return emailx.SendTemplateEmail(
		c.Payload(),
		"templates/invoice_paid_receipt.html",
		"Invoice Paid",
		"Dev Bot <dev@example.test>",
		[]string{"you@example.test"},
	)
}
