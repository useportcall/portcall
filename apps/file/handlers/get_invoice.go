package handlers

import (
	"log"
	"math"
	"net/http"
	"text/template"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func GetInvoice(c *routerx.Context) {
	invoiceID := c.Param("invoice_id")

	var invoice models.Invoice
	if err := c.DB().FindFirst(&invoice, "public_id = ?", invoiceID); err != nil {
		log.Printf("Invoice not found: %v", err)
		c.String(http.StatusNotFound, "Invoice not found")
		return
	}

	var billingAddress models.Address
	if err := c.DB().FindForID(invoice.BillingAddressID, &billingAddress); err != nil {
		log.Printf("Billing address not found: %v", err)
		c.String(http.StatusInternalServerError, "Billing address not found")
		return
	}

	var companyAddress models.Address
	if err := c.DB().FindForID(invoice.CompanyAddressID, &companyAddress); err != nil {
		log.Printf("Company address not found: %v", err)
		c.String(http.StatusInternalServerError, "Company address not found")
		return
	}

	var invoiceItems []models.InvoiceItem
	if err := c.DB().List(&invoiceItems, "invoice_id = ?", invoice.ID); err != nil {
		log.Printf("Error loading invoice items: %v", err)
		c.String(http.StatusInternalServerError, "Error loading invoice items")
		return
	}

	log.Printf("Loaded invoice: %+v", invoice)
	log.Printf("Loaded invoice items: %+v", invoiceItems)

	companyName := invoice.CompanyName
	invoiceNumber := invoice.InvoiceNumber
	invoiceDate := invoice.CreatedAt.Format("2006-01-02")
	customerName := invoice.CustomerName

	// parse template
	tmpl, err := template.ParseFiles("templates/invoice.html")
	if err != nil {
		log.Printf("Error parsing templates/invoice.html: %v", err)
		c.String(http.StatusInternalServerError, "Template error")
		return
	}

	// map addresses to template-friendly shape
	makeAddr := func(a models.Address, nameFallback string) map[string]string {
		street := a.Line1
		if a.Line2 != "" {
			street = street + " " + a.Line2
		}
		name := nameFallback
		return map[string]string{
			"Name":    name,
			"Street":  street,
			"City":    a.City,
			"State":   a.State,
			"Zip":     a.PostalCode,
			"Country": a.Country,
		}
	}

	companyAddr := makeAddr(companyAddress, companyName)
	billingAddr := makeAddr(billingAddress, customerName)

	multiplier := math.Pow(10, float64(invoice.DecimalPlaces))

	tmplItems := make([]map[string]interface{}, 0, len(invoiceItems))
	var computedSubtotal float64
	for _, it := range invoiceItems {
		if (it.Amount == 0 || it.Quantity == 0) && it.PricingModel != "fixed" {
			continue
		}

		qty := int(it.Quantity)
		unit := 0.0
		if qty > 0 {
			unit = float64(it.Amount) / float64(qty) / multiplier
		}
		total := float64(it.Total) / multiplier
		computedSubtotal += total
		tmplItems = append(tmplItems, map[string]interface{}{
			"Title":       it.Title,
			"Description": it.Description,
			"Quantity":    qty,
			"UnitPrice":   unit,
			"Total":       total,
		})
	}

	subtotal := float64(invoice.Total-invoice.TaxAmount) / multiplier
	tax := float64(invoice.TaxAmount) / multiplier
	total := float64(invoice.Total) / multiplier

	data := map[string]interface{}{
		"CompanyName":    companyName,
		"CompanyAddress": companyAddr,
		"BillingAddress": billingAddr,
		"InvoiceNumber":  invoiceNumber,
		"InvoiceDate":    invoiceDate,
		"CustomerName":   customerName,
		"Items":          tmplItems,
		"Subtotal":       subtotal,
		"Tax":            tax,
		"Total":          total,
	}

	c.Header("Content-Type", "text/html")
	c.Status(http.StatusOK)
	if err := tmpl.Execute(c.Writer, data); err != nil {
		log.Printf("Template execution error: %v", err)
		c.String(http.StatusInternalServerError, "Template execution error")
		return
	}
}
