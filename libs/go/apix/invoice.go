package apix

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type Invoice struct {
	ID             string     `json:"id"`
	SubscriptionID *string    `json:"subscription_id"`
	InvoiceNumber  string     `json:"invoice_number"`
	Currency       string     `json:"currency"`
	DueBy          *time.Time `json:"due_date"`
	Status         string     `json:"status"`
	Total          *int64     `json:"total"`
	SubTotal       *int64     `json:"subtotal"`
	PdfURL         *string    `json:"pdf_url"`
	EmailURL       *string    `json:"email_url"`
	RecipientEmail string     `json:"recipient_email"`
	RecipientID    string     `json:"recipient_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	Items          []any      `json:"items"`
}

func (i *Invoice) Set(invoice *models.Invoice) *Invoice {
	i.ID = invoice.PublicID
	i.Status = invoice.Status
	i.CreatedAt = invoice.CreatedAt
	i.UpdatedAt = invoice.UpdatedAt
	i.Currency = invoice.Currency
	i.InvoiceNumber = invoice.InvoiceNumber
	i.DueBy = &invoice.DueBy
	i.PdfURL = &invoice.PDFURL
	i.EmailURL = &invoice.EmailURL
	i.Total = &invoice.Total
	i.SubTotal = &invoice.SubTotal
	return i
}
