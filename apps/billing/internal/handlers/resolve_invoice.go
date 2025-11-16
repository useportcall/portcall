package handlers

import (
	"encoding/json"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

type ResolveInvoicePayload struct {
	InvoiceID uint `json:"invoice_id"`
}

func ResolveInvoice(c server.IContext) error {
	var p ResolveInvoicePayload
	if err := json.Unmarshal(c.Payload(), &p); err != nil {
		return err
	}

	var invoice models.Invoice
	if err := c.DB().FindForID(p.InvoiceID, &invoice); err != nil {
		return err
	}

	invoice.Status = "paid"
	if err := c.DB().Save(&invoice); err != nil {
		return err
	}

	return nil
}
