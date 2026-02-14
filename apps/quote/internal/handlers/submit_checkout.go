package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/paymentx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

// handleDirectCheckout processes a quote that requires immediate payment via checkout.
func handleDirectCheckout(
	c *routerx.Context,
	quote *models.Quote, plan *models.Plan,
	user *models.User, company *models.Company,
	quoteAppURL, checkoutAppURL, quoteLink, recipientEmail string,
	now time.Time,
) {
	var config models.AppConfig
	if err := c.DB().FindFirst(&config, "app_id = ?", quote.AppID); err != nil {
		c.ServerError("error retrieving app config for app", err)
		return
	}

	var connection models.Connection
	if err := c.DB().FindForID(config.DefaultConnectionID, &connection); err != nil {
		c.ServerError("error retrieving payment connection for app", err)
		return
	}

	payment, err := paymentx.New(&connection, c.Crypto())
	if err != nil {
		c.ServerError("error creating payment client", err)
		return
	}

	if err := ensurePaymentCustomer(c, user, payment); err != nil {
		return
	}

	sessionID, clientSecret, err := payment.CreateCheckoutSession(user.PaymentCustomerID)
	if err != nil {
		c.ServerError("error creating payment checkout session", err)
		return
	}

	checkoutSession := buildCheckoutSession(quote, plan, user, company, &connection, quoteAppURL, sessionID, clientSecret)
	if err := c.DB().Create(checkoutSession); err != nil {
		c.ServerError("error creating checkout session", err)
		return
	}

	markQuoteAccepted(quote, now)
	if err := c.DB().Save(quote); err != nil {
		c.ServerError("error saving accepted quote", err)
		return
	}

	enqueueAcceptanceEmails(c, user, company, quote, quoteLink, recipientEmail, now)

	checkoutURL, err := buildCheckoutRedirectURL(checkoutAppURL, checkoutSession, c.Crypto())
	if err != nil {
		c.ServerError("error building checkout URL", err)
		return
	}
	c.Redirect(303, checkoutURL)
}

func enqueueAcceptanceEmails(c *routerx.Context, user *models.User, company *models.Company, quote *models.Quote, quoteLink, recipientEmail string, now time.Time) {
	if err := c.Queue().Enqueue("send_quote_accepted_confirmation_email", map[string]any{
		"RecipientName": user.Name, "QuoteLink": quoteLink,
		"Year": now.Year(), "recipient_email": recipientEmail,
	}, "email_queue"); err != nil {
		log.Printf("error enqueueing acceptance email: %v", err)
	}
	if company.Email != "" {
		if err := c.Queue().Enqueue("send_quote_status_email", map[string]any{
			"recipient_email": company.Email,
			"subject":         fmt.Sprintf("Quote accepted by %s", recipientEmail),
			"title":           "Quote accepted",
			"message":         fmt.Sprintf("Quote %s was accepted.", quote.PublicID),
			"action_text":     "View quote",
			"action_url":      quoteLink,
			"year":            now.Year(),
		}, "email_queue"); err != nil {
			log.Printf("error enqueueing status email: %v", err)
		}
	}
}
