package invoice

import (
	"fmt"
	"os"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// CreateUpgrade creates an invoice for the prorated price difference
// when upgrading plans. Creates invoice + item + totals in a single transaction.
func (s *service) CreateUpgrade(input *CreateUpgradeInput) (*CreateUpgradeResult, error) {
	invoiceAppURL := os.Getenv("INVOICE_APP_URL")
	if invoiceAppURL == "" {
		return nil, fmt.Errorf("INVOICE_APP_URL environment variable is not set")
	}

	sub, err := findSubscription(s.db, input.SubscriptionID)
	if err != nil {
		return nil, err
	}

	var company models.Company
	if err := s.db.FindFirst(&company, "app_id = ?", sub.AppID); err != nil {
		return nil, err
	}

	var count int64
	if err := s.db.Count(&count, &models.Invoice{}, "app_id = ?", sub.AppID); err != nil {
		return nil, err
	}

	var oldPlan, newPlan models.Plan
	if err := s.db.FindForID(input.OldPlanID, &oldPlan); err != nil {
		return nil, err
	}
	if err := s.db.FindForID(input.NewPlanID, &newPlan); err != nil {
		return nil, err
	}

	publicID := dbx.GenPublicID("invoice")
	invoice := buildUpgradeInvoice(sub, &company, count, publicID, invoiceAppURL)
	item := buildUpgradeItem(sub.AppID, &oldPlan, &newPlan, input.PriceDifference)

	invoice.SubTotal = input.PriceDifference
	invoice.Total = input.PriceDifference
	invoice.Status = "issued"

	if err := s.db.Txn(func(tx dbx.IORM) error {
		if err := tx.Create(invoice); err != nil {
			return err
		}
		item.InvoiceID = invoice.ID
		return tx.Create(item)
	}); err != nil {
		return nil, err
	}

	return &CreateUpgradeResult{
		Invoice:   invoice,
		ShouldPay: invoice.Total > 0,
	}, nil
}

func buildUpgradeInvoice(
	sub *models.Subscription, company *models.Company,
	count int64, publicID, invoiceAppURL string,
) *models.Invoice {
	return &models.Invoice{
		AppID: sub.AppID, SubscriptionID: &sub.ID, UserID: sub.UserID,
		PublicID: publicID, Status: "pending", Currency: sub.Currency,
		PDFURL:             fmt.Sprintf("%s/invoices/%s/view", invoiceAppURL, publicID),
		EmailURL:           fmt.Sprintf("%s/invoice-email/%s", invoiceAppURL, publicID),
		DueBy:              time.Now().AddDate(0, 0, sub.InvoiceDueByDays),
		InvoiceNumber:      fmt.Sprintf("INV-%07d", count+1),
		InvoiceNumberCount: count + 1,
		CompanyAddressID:   company.BillingAddressID,
		BillingAddressID:   *sub.BillingAddressID,
		ShippingAddressID:  sub.BillingAddressID,
	}
}

func buildUpgradeItem(
	appID uint, oldPlan, newPlan *models.Plan, priceDiff int64,
) *models.InvoiceItem {
	return &models.InvoiceItem{
		PublicID:     dbx.GenPublicID("ii"),
		AppID:        appID,
		Quantity:     1,
		Title:        fmt.Sprintf("Plan Upgrade: %s \u2192 %s", oldPlan.Name, newPlan.Name),
		Description:  fmt.Sprintf("Pro-rated upgrade from %s to %s", oldPlan.Name, newPlan.Name),
		PricingModel: "fixed",
		Total:        priceDiff,
		Amount:       priceDiff,
	}
}
