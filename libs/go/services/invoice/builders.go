package invoice

import (
	"fmt"
	"os"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func buildInvoice(sub *models.Subscription, company *models.Company, count int64) *models.Invoice {
	invoiceAppURL := os.Getenv("INVOICE_APP_URL")
	discountPct := 0
	if sub.DiscountPct > 0 && count+1 <= int64(sub.DiscountQty) {
		discountPct = sub.DiscountPct
	}
	publicID := dbx.GenPublicID("invoice")
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
		DiscountPct:        discountPct,
	}
}
