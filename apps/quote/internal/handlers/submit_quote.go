package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

// SubmitQuote handles accepting or declining a quote via form submission.
func SubmitQuote(c *routerx.Context) {
	quoteAppURL := os.Getenv("QUOTE_APP_URL")
	checkoutAppURL := os.Getenv("CHECKOUT_APP_URL")
	if quoteAppURL == "" || checkoutAppURL == "" {
		c.ServerError("QUOTE_APP_URL or CHECKOUT_APP_URL is not set", nil)
		return
	}

	quote, plan, user, company, ok := loadSubmitDeps(c)
	if !ok {
		return
	}

	imgBytes, ok := validateSignature(c)
	if !ok {
		return
	}
	if err := c.Store().PutInSignatureBucket(quote.PublicID, imgBytes, c); err != nil {
		c.ServerError("Failed to save signature image", err)
		return
	}

	now := time.Now().UTC()
	recipientEmail := strings.TrimSpace(quote.RecipientEmail)
	if recipientEmail == "" {
		recipientEmail = user.Email
	}

	quoteLink, err := buildQuoteLink(c.Crypto(), quoteAppURL, quote, now)
	if err != nil {
		c.ServerError("Failed to create quote access token", err)
		return
	}

	if !quote.DirectCheckout {
		handleNonCheckoutAccept(c, quote, user, company, quoteLink, recipientEmail, now)
		return
	}
	handleDirectCheckout(c, quote, plan, user, company, quoteAppURL, checkoutAppURL, quoteLink, recipientEmail, now)
}

func loadSubmitDeps(c *routerx.Context) (*models.Quote, *models.Plan, *models.User, *models.Company, bool) {
	var quote models.Quote
	if err := c.DB().FindFirst(&quote, "public_id = ?", c.Param("id")); err != nil {
		c.String(http.StatusNotFound, "quote not found")
		return nil, nil, nil, nil, false
	}
	if !verifyQuoteAccess(c, &quote) {
		return nil, nil, nil, nil, false
	}
	if !quoteCanBeAccepted(&quote, time.Now().UTC()) {
		c.String(http.StatusBadRequest, "quote is no longer accept-able")
		return nil, nil, nil, nil, false
	}
	if quote.UserID == nil {
		c.String(http.StatusBadRequest, "quote has no associated user")
		return nil, nil, nil, nil, false
	}
	var company models.Company
	if err := c.DB().FindFirst(&company, "app_id = ?", quote.AppID); err != nil {
		c.ServerError("error retrieving company for app", err)
		return nil, nil, nil, nil, false
	}
	var plan models.Plan
	if err := c.DB().FindForID(quote.PlanID, &plan); err != nil {
		c.String(http.StatusNotFound, "quote not found")
		return nil, nil, nil, nil, false
	}
	var user models.User
	if err := c.DB().FindForID(*quote.UserID, &user); err != nil {
		c.String(http.StatusNotFound, "quote not found")
		return nil, nil, nil, nil, false
	}
	return &quote, &plan, &user, &company, true
}

func handleNonCheckoutAccept(c *routerx.Context, quote *models.Quote, user *models.User, company *models.Company, quoteLink, recipientEmail string, now time.Time) {
	markQuoteAccepted(quote, now)
	if err := c.DB().Save(quote); err != nil {
		c.ServerError("error saving quote with signature", err)
		return
	}
	enqueueAcceptanceEmails(c, user, company, quote, quoteLink, recipientEmail, now)
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/quotes/%s/success", quote.PublicID))
}
