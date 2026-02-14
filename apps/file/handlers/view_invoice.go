package handlers

import (
	"fmt"
	"log"
	"net/http"
	"html/template"

	"github.com/useportcall/portcall/libs/go/routerx"
)

func ViewInvoice(c *routerx.Context) {
	invoiceId := c.Param("invoice_id")
	iframeUrl := fmt.Sprintf("/invoice/%s", invoiceId)
	pdfUrl := fmt.Sprintf("/pdf/%s", invoiceId)
	c.Header("Content-Type", "text/html")

	tmpl, err := template.ParseFiles(tmplPaths("view.html", "view-styles.html")...)
	if err != nil {
		log.Printf("Error loading view template: %v", err)
		c.String(http.StatusInternalServerError, "Template error")
		return
	}

	data := map[string]string{
		"Src":    iframeUrl,
		"PdfUrl": pdfUrl,
	}

	if err := tmpl.Execute(c.Writer, data); err != nil {
		log.Printf("Template execution error: %v", err)
		c.String(http.StatusInternalServerError, "Template error")
	}
}
