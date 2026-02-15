package control

import (
	"encoding/json"
	"net/http"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func (s *Server) handleSeedInvoice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	invoice, err := seedTestInvoice(s.H.DB, s.H.AppID)
	if err != nil {
		http.Error(w, "seed failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"invoice_id": invoice.PublicID})
}

func seedTestInvoice(db dbx.IORM, appID uint) (*models.Invoice, error) {
	companyAddr := &models.Address{
		PublicID: dbx.GenPublicID("addr"), AppID: appID,
		Line1: "100 Main St", City: "San Francisco",
		State: "CA", PostalCode: "94105", Country: "US",
	}
	if err := db.Create(companyAddr); err != nil {
		return nil, err
	}

	billingAddr := &models.Address{
		PublicID: dbx.GenPublicID("addr"), AppID: appID,
		Line1: "200 Elm St", City: "New York",
		State: "NY", PostalCode: "10001", Country: "US",
	}
	if err := db.Create(billingAddr); err != nil {
		return nil, err
	}

	inv := &models.Invoice{
		PublicID: dbx.GenPublicID("invoice"), AppID: appID,
		UserID: 1, Currency: "usd", Total: 12900, SubTotal: 12900,
		TaxAmount: 0, DecimalPlaces: 2, Status: "paid",
		InvoiceNumber: "INV-0000001",
		CompanyName: "Acme Corp", CustomerName: "Jane Doe",
		CompanyAddressID: companyAddr.ID, BillingAddressID: billingAddr.ID,
	}
	if err := db.Create(inv); err != nil {
		return nil, err
	}

	items := []models.InvoiceItem{
		{PublicID: dbx.GenPublicID("ii"), AppID: appID,
			InvoiceID: inv.ID, Title: "Pro Plan",
			Description: "Monthly subscription",
			Quantity: 1, Amount: 9900, Total: 9900, PricingModel: "fixed"},
		{PublicID: dbx.GenPublicID("ii"), AppID: appID,
			InvoiceID: inv.ID, Title: "API Calls",
			Description: "1,000 requests",
			Quantity: 1000, Amount: 3000, Total: 3000, PricingModel: "metered"},
	}
	for i := range items {
		if err := db.Create(&items[i]); err != nil {
			return nil, err
		}
	}

	return inv, nil
}
