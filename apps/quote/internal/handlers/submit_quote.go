package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/paymentx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func SubmitQuote(c *routerx.Context) {
	quoteAppURL := os.Getenv("QUOTE_APP_URL")

	var quote models.Quote
	if err := c.DB().FindFirst(&quote, "public_id = ?", c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "quote not found")
		return
	}

	var company models.Company
	if err := c.DB().FindFirst(&company, "app_id = ?", quote.AppID); err != nil {
		c.ServerError("error retrieving company for app", err)
		return
	}

	var plan models.Plan
	if err := c.DB().FindForID(quote.PlanID, &plan); err != nil {
		c.String(http.StatusNotFound, "quote not found")
		return
	}

	var user models.User
	if err := c.DB().FindForID(*quote.UserID, &user); err != nil {
		c.String(http.StatusNotFound, "quote not found")
		return
	}

	signatureData := c.PostForm("signatureData")
	if signatureData == "" {
		c.String(http.StatusBadRequest, "No signature data provided")
		return
	}

	// Remove data URL prefix if present
	base64Data := signatureData
	if strings.HasPrefix(signatureData, "data:image/png;base64,") {
		base64Data = strings.TrimPrefix(signatureData, "data:image/png;base64,")
	}

	imgBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid base64 data: %v", err)
		return
	}

	// Save the image to the signature bucket
	if err := c.Store().PutInSignatureBucket(quote.PublicID, imgBytes, c); err != nil {
		c.ServerError("Failed to save signature image", err)
		return
	}

	if !quote.DirectCheckout {
		quote.Status = "accepted"

		if err := c.DB().Save(&quote); err != nil {
			c.ServerError("error saving quote with signature", err)
			return
		}

		// send quote accepted confirmation email
		if err := c.Queue().Enqueue("send_quote_accepted_confirmation_email", map[string]any{
			"RecipientName": user.Name,
			"QuoteLink":     fmt.Sprintf("%s/quotes/%s", quoteAppURL, quote.PublicID),
			"Year":          time.Now().Year(),
		}, "email_queue"); err != nil {
			c.ServerError("error enqueueing quote accepted confirmation email", err)
			return
		}

		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/quotes/%s/success", quote.PublicID))
		return
	}

	var config models.AppConfig
	if err := c.DB().FindFirst(&config, "app_id = ?", quote.AppID); err != nil {
		c.ServerError("error retrieving app config for app", err)
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

	if user.PaymentCustomerID == "" {
		customerID, err := payment.CreateCustomer(user.Email, user.Name)
		if err != nil {
			c.ServerError("error creating payment customer", err)
			return
		}

		user.PaymentCustomerID = customerID
		if err := c.DB().Save(&user); err != nil {
			c.ServerError("error saving user with payment customer ID", err)
			return
		}
	}

	sessionID, clientSecret, err := payment.CreateCheckoutSession(user.PaymentCustomerID)
	if err != nil {
		c.ServerError("error creating payment checkout session", err)
		return
	}

	cancelURL := fmt.Sprintf("%s/quotes/%s/cancel", quoteAppURL, quote.PublicID)
	redirectURL := fmt.Sprintf("%s/quotes/%s/success", quoteAppURL, quote.PublicID)

	checkoutSession := &models.CheckoutSession{
		PublicID:             dbx.GenPublicID("cs"),
		ExpiresAt:            time.Now().Add(time.Hour * 24 * 2),
		AppID:                quote.AppID,
		UserID:               user.ID,
		PlanID:               plan.ID,
		ExternalClientSecret: clientSecret,
		ExternalPublicKey:    connection.PublicKey,
		ExternalSessionID:    sessionID,
		ExternalProvider:     connection.Source,
		CancelURL:            &cancelURL,   // TODO: add cancel url to body
		RedirectURL:          &redirectURL, // TODO: add redirect url to body
		CompanyAddressID:     company.BillingAddressID,
	}
	if err := c.DB().Create(checkoutSession); err != nil {
		c.ServerError("error creating checkout session", err)
		return
	}

	// redirect to checkout view
	checkoutAppURL := os.Getenv("CHECKOUT_APP_URL")
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("%s?id=%s", checkoutAppURL, checkoutSession.PublicID)) // TODO: use real frontend url
}
