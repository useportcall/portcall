package payment

import (
	"fmt"
	"log"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// Resolve marks an invoice as paid and prepares email data.
// Single mutation: updates invoice status to paid.
func (s *service) Resolve(input *ResolveInput) (*ResolveResult, error) {
	log.Printf("Processing Resolve for invoice ID %d", input.InvoiceID)

	invoice, err := findInvoice(s.db, input.InvoiceID)
	if err != nil {
		return nil, err
	}

	user, company, err := lookupResolveDeps(s.db, invoice.UserID, invoice.AppID)
	if err != nil {
		return nil, err
	}

	invoice.Status = "paid"
	if err := s.db.Save(invoice); err != nil {
		return nil, err
	}
	if err := recoverSubscriptionStatus(s.db, invoice.SubscriptionID); err != nil {
		return nil, err
	}

	log.Printf("Invoice %d marked as paid", input.InvoiceID)

	// Skip email for zero-amount invoices
	if invoice.Total < 1 {
		return &ResolveResult{Invoice: invoice}, nil
	}

	amountPaid := float64(invoice.Total) / 100.0
	return &ResolveResult{
		Invoice: invoice,
		EmailPayload: &InvoicePaidEmailPayload{
			InvoiceNumber:  invoice.InvoiceNumber,
			AmountPaid:     fmt.Sprintf("$%.2f", amountPaid),
			DatePaid:       time.Now().Format("January 2, 2006"),
			CompanyName:    company.Name,
			RecipientEmail: user.Email,
			LogoURL:        emailLogoURL(company),
			PaymentStatus:  "paid",
		},
	}, nil
}

func lookupResolveDeps(db dbx.IORM, userID, appID uint) (*models.User, *models.Company, error) {
	var user models.User
	if err := db.FindForID(userID, &user); err != nil {
		return nil, nil, err
	}

	var company models.Company
	if err := db.FindFirst(&company, "app_id = ?", appID); err != nil {
		return nil, nil, err
	}

	return &user, &company, nil
}
