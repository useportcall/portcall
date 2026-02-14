package handlers

import (
	"math"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// makeAddr converts a DB address into a template-friendly map.
func makeAddr(a models.Address, name string) map[string]string {
	street := a.Line1
	if a.Line2 != "" {
		street += " " + a.Line2
	}
	return map[string]string{
		"Name": name, "Street": street,
		"City": a.City, "State": a.State,
		"Zip": a.PostalCode, "Country": a.Country,
	}
}

// buildItems converts invoice items to template-friendly maps.
func buildItems(items []models.InvoiceItem, mult float64) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(items))
	for _, it := range items {
		if (it.Amount == 0 || it.Quantity == 0) && it.PricingModel != "fixed" {
			continue
		}
		qty := int(it.Quantity)
		unit := 0.0
		if qty > 0 {
			unit = float64(it.Amount) / float64(qty) / mult
		}
		out = append(out, map[string]interface{}{
			"Title":       it.Title,
			"Description": it.Description,
			"Quantity":    qty,
			"UnitPrice":   unit,
			"Total":       float64(it.Total) / mult,
		})
	}
	return out
}

// buildInvoiceData assembles the full template data map.
func buildInvoiceData(
	inv models.Invoice,
	companyAddr, billingAddr map[string]string,
	items []map[string]interface{},
) map[string]interface{} {
	m := math.Pow(10, float64(inv.DecimalPlaces))
	return map[string]interface{}{
		"CompanyName":    inv.CompanyName,
		"CompanyAddress": companyAddr,
		"BillingAddress": billingAddr,
		"InvoiceNumber":  inv.InvoiceNumber,
		"InvoiceDate":    inv.CreatedAt.Format("2006-01-02"),
		"CustomerName":   inv.CustomerName,
		"Items":          items,
		"Subtotal":       float64(inv.Total-inv.TaxAmount) / m,
		"Tax":            float64(inv.TaxAmount) / m,
		"Total":          float64(inv.Total) / m,
	}
}
