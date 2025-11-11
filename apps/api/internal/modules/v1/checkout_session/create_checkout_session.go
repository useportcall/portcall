package checkout_session

import (
	"fmt"
	"time"

	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/paymentx"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type CreateCheckoutSessionRequest struct {
	PlanID      string `json:"plan_id"`
	UserID      string `json:"user_id"`
	CancelURL   string `json:"cancel_url"`
	RedirectURL string `json:"redirect_url"`
}

func CreateCheckoutSession(c *routerx.Context) {
	var body CreateCheckoutSessionRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("invalid request body")
		return
	}

	var company models.Company
	if err := c.DB().FindFirst(&company, "app_id = ?", c.AppID); err != nil {
		c.ServerError("error retrieving company for app")
		return
	}

	if company.BillingAddressID == 0 {
		c.BadRequest("company does not have a billing address")
		return
	}

	var plan models.Plan
	if err := c.DB().GetForPublicID(c.AppID(), body.PlanID, &plan); err != nil {
		c.ServerError("error retrieving plan")
		return
	}

	if plan.Status != "published" {
		c.BadRequest(fmt.Sprintf("plan with id '%s' is not yet published", plan.PublicID))
		return
	}

	var config models.AppConfig
	if err := c.DB().FindFirst(&config, "app_id = ?", c.AppID()); err != nil {
		c.ServerError("error retrieving app config for app")
	}

	var connection models.Connection
	if err := c.DB().FindForID(config.DefaultConnectionID, &connection); err != nil {
		c.ServerError("error retrieving payment connection for app")
		return
	}

	payment, err := paymentx.New(&connection, c.Crypto())
	if err != nil {
		c.ServerError("error creating payment client")
		return
	}

	var user models.User
	if err := c.DB().GetForPublicID(c.AppID(), body.UserID, &user); err != nil {
		c.ServerError(fmt.Sprintf("error retrieving user %s for app %d: %v", body.UserID, c.AppID(), err))
		return
	}

	if user.PaymentCustomerID == "" {
		customerID, err := payment.CreateCustomer(user.Email, user.Name)
		if err != nil {
			c.ServerError("error creating payment customer")
			return
		}

		user.PaymentCustomerID = customerID
		if err := c.DB().Save(&user); err != nil {
			c.ServerError("error saving user with payment customer ID")
			return
		}
	}

	sessionID, clientSecret, err := payment.CreateCheckoutSession(user.PaymentCustomerID)
	if err != nil {
		c.ServerError("error creating payment checkout session")
		return
	}

	checkoutSession := &models.CheckoutSession{
		PublicID:             dbx.GenPublicID("cs"),
		ExpiresAt:            time.Now().Add(time.Hour * 24 * 2),
		AppID:                c.AppID(),
		UserID:               user.ID,
		PlanID:               plan.ID,
		ExternalClientSecret: clientSecret,
		ExternalPublicKey:    connection.PublicKey,
		ExternalSessionID:    sessionID,
		ExternalProvider:     connection.Source,
		CancelURL:            &body.CancelURL,
		RedirectURL:          &body.RedirectURL,
		CompanyAddressID:     company.BillingAddressID,
	}
	if err := c.DB().Create(checkoutSession); err != nil {
		c.ServerError("error creating checkout session")
		return
	}

	response := new(apix.CheckoutSession).Set(checkoutSession)

	c.OK(response)
}
