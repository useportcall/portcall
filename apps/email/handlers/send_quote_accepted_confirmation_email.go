package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"text/template"

	"github.com/useportcall/portcall/libs/go/qx/server"
)

func SendQuoteAcceptedConfirmationEmail(c server.IContext) error {
	var body any
	if err := json.Unmarshal(c.Payload(), &body); err != nil {
		return err
	}

	tmpl, err := template.ParseFiles("templates/quote_accepted_confirmation.html")
	if err != nil {
		log.Fatal("failed to parse template:", err)
	}

	var htmlContentBuf bytes.Buffer
	if err := tmpl.Execute(&htmlContentBuf, body); err != nil {
		log.Fatal("failed to execute template:", err)
	}
	content := htmlContentBuf.String()

	return c.EmailClient().Send(
		content,
		"Quote Issued",
		"dev@example.test",
		[]string{"you@example.test"},
	)
}
