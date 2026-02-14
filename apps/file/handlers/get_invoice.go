package handlers

import (
	"log"
	"math"
	"net/http"
	"html/template"

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
		c.String(http.StatusInternalServerError, "Failed to load invoice data")
		return
	}

	var companyAddress models.Address
	if err := c.DB().FindForID(invoice.CompanyAddressID, &companyAddress); err != nil {
		log.Printf("Company address not found: %v", err)
		c.String(http.StatusInternalServerError, "Failed to load invoice data")
		return
	}

	var invoiceItems []models.InvoiceItem
	if err := c.DB().List(&invoiceItems, "invoice_id = ?", invoice.ID); err != nil {
		log.Printf("Error loading invoice items: %v", err)
		c.String(http.StatusInternalServerError, "Failed to load invoice data")
		return
	}

	tmpl, err := template.ParseFiles(tmplPaths("invoice.html", "invoice-styles.html")...)
	if err != nil {
		log.Printf("Error parsing invoice template: %v", err)
		c.String(http.StatusInternalServerError, "Template error")
		return
	}

	mult := math.Pow(10, float64(invoice.DecimalPlaces))
	companyAddr := makeAddr(companyAddress, invoice.CompanyName)
	billingAddr := makeAddr(billingAddress, invoice.CustomerName)
	items := buildItems(invoiceItems, mult)
	data := buildInvoiceData(invoice, companyAddr, billingAddr, items)

	c.Header("Content-Type", "text/html")
	c.Status(http.StatusOK)
	if err := tmpl.Execute(c.Writer, data); err != nil {
		log.Printf("Template execution error: %v", err)
		c.String(http.StatusInternalServerError, "Template execution error")
	}
}
