package handlers

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/paymentx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ensurePaymentCustomer(c *routerx.Context, user *models.User, payment paymentx.IPaymentClient) error {
	if user.PaymentCustomerID != "" {
		return nil
	}
	customerID, err := payment.CreateCustomer(user.Email, user.Name)
	if err != nil {
		c.ServerError("error creating payment customer", err)
		return err
	}
	user.PaymentCustomerID = customerID
	if err := c.DB().Save(user); err != nil {
		c.ServerError("error saving user with payment customer ID", err)
		return err
	}
	return nil
}

func buildCheckoutSession(
	quote *models.Quote, plan *models.Plan, user *models.User, company *models.Company,
	conn *models.Connection, quoteAppURL, sessionID, clientSecret string,
) *models.CheckoutSession {
	cancelURL := fmt.Sprintf("%s/quotes/%s/cancel", quoteAppURL, quote.PublicID)
	redirectURL := fmt.Sprintf("%s/quotes/%s/success", quoteAppURL, quote.PublicID)
	return &models.CheckoutSession{
		PublicID:             dbx.GenPublicID("cs"),
		ExpiresAt:            time.Now().Add(48 * time.Hour),
		AppID:                quote.AppID,
		UserID:               user.ID,
		PlanID:               plan.ID,
		ExternalClientSecret: clientSecret,
		ExternalPublicKey:    conn.PublicKey,
		ExternalSessionID:    sessionID,
		ExternalProvider:     conn.Source,
		CancelURL:            &cancelURL,
		RedirectURL:          &redirectURL,
		CompanyAddressID:     &company.BillingAddressID,
	}
}

func buildCheckoutRedirectURL(baseURL string, session *models.CheckoutSession, crypto cryptox.ICrypto) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil {
		return "", err
	}
	query := parsed.Query()
	query.Set("id", session.PublicID)
	if crypto != nil {
		token, err := cryptox.CreateCheckoutSessionToken(crypto, session.PublicID, session.ExpiresAt)
		if err != nil {
			return "", err
		}
		query.Set("st", token)
	}
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}
