package payment

import (
	"log"
	"strconv"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/paymentx"
)

// Pay processes payment for an invoice.
// This calls the external payment provider - not a database mutation.
func (s *service) Pay(input *PayInput) (*PayResult, error) {
	log.Printf("Processing Pay for invoice ID %d", input.InvoiceID)

	invoice, err := findInvoice(s.db, input.InvoiceID)
	if err != nil {
		return nil, err
	}

	// Skip payment if zero or negative amount
	if invoice.Total <= 0 {
		log.Printf("Invoice %d has zero/negative total, skipping payment", input.InvoiceID)
		return &PayResult{Invoice: invoice}, nil
	}

	user, pm, conn, err := lookupPaymentDeps(s.db, invoice.UserID, invoice.AppID)
	if err != nil {
		return nil, err
	}

	client, err := paymentx.New(conn, s.crypto)
	if err != nil {
		return nil, err
	}

	if err := client.CreateCharge(
		user.PaymentCustomerID,
		invoice.Total,
		invoice.Currency,
		pm.ExternalID,
		buildChargeMetadata(invoice),
	); err != nil {
		return nil, err
	}

	log.Printf("Payment processed for invoice %d", input.InvoiceID)
	return &PayResult{Invoice: invoice}, nil
}

func findInvoice(db dbx.IORM, id uint) (*models.Invoice, error) {
	var invoice models.Invoice
	if err := db.FindForID(id, &invoice); err != nil {
		return nil, err
	}
	return &invoice, nil
}

func lookupPaymentDeps(db dbx.IORM, userID, appID uint) (*models.User, *models.PaymentMethod, *models.Connection, error) {
	var user models.User
	if err := db.FindForID(userID, &user); err != nil {
		return nil, nil, nil, err
	}

	pm, err := findLatestPaymentMethod(db, userID, appID)
	if err != nil {
		return nil, nil, nil, err
	}

	var conn models.Connection
	if err := db.FindFirst(&conn, "app_id = ?", appID); err != nil {
		return nil, nil, nil, err
	}

	return &user, pm, &conn, nil
}

func buildChargeMetadata(invoice *models.Invoice) map[string]string {
	metadata := map[string]string{
		"portcall_invoice_id": strconv.FormatUint(uint64(invoice.ID), 10),
		"portcall_app_id":     strconv.FormatUint(uint64(invoice.AppID), 10),
	}
	if invoice.SubscriptionID != nil {
		metadata["portcall_subscription_id"] = strconv.FormatUint(uint64(*invoice.SubscriptionID), 10)
	}
	return metadata
}
