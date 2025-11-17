package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/paymentx"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

type PayInvoicePayload struct {
	InvoiceID uint `json:"invoice_id"`
}

func PayInvoice(c server.IContext) error {
	var p PayInvoicePayload
	if err := json.Unmarshal(c.Payload(), &p); err != nil {
		return err
	}

	var invoice models.Invoice
	if err := c.DB().FindForID(p.InvoiceID, &invoice); err != nil {
		return err
	}

	var user models.User
	if err := c.DB().FindForID(invoice.UserID, &user); err != nil {
		return err
	}

	var company models.Company
	if err := c.DB().FindFirst(&company, "app_id = ?", invoice.AppID); err != nil {
		return err
	}

	var paymentMethod models.PaymentMethod
	if err := c.DB().FindFirst(&paymentMethod, "user_id = ?", invoice.UserID); err != nil {
		return err
	}

	var connection models.Connection
	if err := c.DB().FindFirst(&connection, "app_id = ?", invoice.AppID); err != nil {
		return err
	}

	payment, err := paymentx.New(&connection, c.Crypto())
	if err != nil {
		return err
	}

	if err := payment.CreateCharge(
		user.PaymentCustomerID,
		invoice.Total,
		invoice.Currency,
		paymentMethod.ExternalID,
	); err != nil {
		return err
	}

	if err := c.Queue().Enqueue("resolve_invoice", map[string]any{"invoice_id": invoice.ID}, "billing_queue"); err != nil {
		return err
	}

	// should be in decimals not cents
	amountPaid := float64(invoice.Total) / 100.0

	payload := map[string]any{
		"invoice_number": invoice.InvoiceNumber,
		"amount_paid":    fmt.Sprintf("$%.2f", amountPaid),
		"date_paid":      time.Now().Format("January 2, 2006"),
		"company_name":   company.Name,
	}
	if err := c.Queue().Enqueue("send_invoice_paid_email", payload, "email_queue"); err != nil {
		return fmt.Errorf("failed to enqueue send_invoice_paid_email task: %w", err)
	}

	return nil
}
