package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

type CalculateInvoiceTotalsPayload struct {
	InvoiceID uint `json:"invoice_id"`
}

func CalculateInvoiceTotals(c server.IContext) error {
	var p CalculateInvoiceTotalsPayload
	if err := json.Unmarshal(c.Payload(), &p); err != nil {
		return err
	}

	var invoice models.Invoice
	if err := c.DB().FindForID(p.InvoiceID, &invoice); err != nil {
		return err
	}

	items := []models.InvoiceItem{}
	if err := c.DB().List(&items, "invoice_id = ?", p.InvoiceID); err != nil {
		return err
	}

	var totalAmount int64
	for _, item := range items {
		totalAmount += item.Total
	}

	invoice.SubTotal = totalAmount
	invoice.TaxAmount = 0
	invoice.DiscountAmount = 0
	invoice.Total = invoice.SubTotal + invoice.TaxAmount - invoice.DiscountAmount
	invoice.Status = "issued"
	if err := c.DB().Save(&invoice); err != nil {
		return err
	}

	if invoice.Total > 0 {
		if err := c.Queue().Enqueue("pay_invoice", map[string]any{"invoice_id": invoice.ID}, "billing_queue"); err != nil {
			return fmt.Errorf("failed to enqueue invoice payment: %w", err)
		}
	} else {
		log.Println("Invoice total is zero or negative, skipping payment:", invoice.ID)

		if err := c.Queue().Enqueue("resolve_invoice", map[string]any{"invoice_id": invoice.ID}, "billing_queue"); err != nil {
			return err
		}
	}

	if err := c.Queue().Enqueue("send_invoice_issued_email", map[string]any{"public_id": invoice.PublicID}, "email_queue"); err != nil {
		return fmt.Errorf("failed to enqueue invoice issued email: %w", err)
	}

	return nil
}
