package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

type CreateInvoicePayload struct {
	SubscriptionID uint `json:"subscription_id"`
}

func CreateInvoice(c server.IContext) error {
	var p CreateInvoicePayload
	if err := json.Unmarshal(c.Payload(), &p); err != nil {
		return err
	}

	var subscription models.Subscription
	if err := c.DB().FindForID(p.SubscriptionID, &subscription); err != nil {
		return err
	}

	var company models.Company
	if err := c.DB().FindFirst(&company, "app_id = ?", subscription.AppID); err != nil {
		return err
	}

	var count int64
	if err := c.DB().Count(&count, &models.Invoice{}, "app_id = ?", subscription.AppID); err != nil {
		return err
	}

	publicID := dbx.GenPublicID("invoice")

	invoice := &models.Invoice{
		AppID:             subscription.AppID,
		SubscriptionID:    &subscription.ID,
		UserID:            subscription.UserID,
		PublicID:          publicID,
		Status:            "pending",
		Currency:          subscription.Currency,
		PDFURL:            fmt.Sprintf("http://localhost:8085/view/%s", publicID),          // TODO: fix
		EmailURL:          fmt.Sprintf("http://localhost:8085/invoice-email/%s", publicID), // TODO: fix
		DueBy:             time.Now().AddDate(0, 0, subscription.InvoiceDueByDays),
		InvoiceNumber:     fmt.Sprintf("INV-%07d", count+1), // invoice number should be INV-0000001 format
		CompanyAddressID:  company.BillingAddressID,
		BillingAddressID:  subscription.BillingAddressID,
		ShippingAddressID: &subscription.BillingAddressID,
	}
	if err := c.DB().Create(invoice); err != nil {
		return err
	}

	if err := c.Queue().Enqueue("create_invoice_items", map[string]any{"invoice_id": invoice.ID}, "billing_queue"); err != nil {
		return fmt.Errorf("failed to enqueue calculate_invoice_totals task: %w", err)
	}

	return nil
}
