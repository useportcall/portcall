package address

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type UpdateCheckoutSessionAddressRequest struct {
	Line1      string `json:"line1"`
	Line2      string `json:"line2"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

func UpdateCheckoutSessionAddress(c *routerx.Context) {
	sessionID := c.Param("id")

	var body UpdateCheckoutSessionAddressRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("invalid request payload")
		return
	}

	session := &models.CheckoutSession{}
	if err := c.DB().FindFirst(session, "public_id = ?", sessionID); err != nil {
		c.NotFound("checkout session not found")
		return
	}

	var user models.User
	if err := c.DB().FindForID(session.UserID, &user); err != nil {
		c.ServerError("internal server error", err)
		return
	}

	address := &models.Address{
		PublicID:   dbx.GenPublicID("addr"),
		AppID:      session.AppID,
		Line1:      body.Line1,
		Line2:      body.Line2,
		City:       body.City,
		State:      body.State,
		PostalCode: body.PostalCode,
		Country:    body.Country,
	}
	if err := c.DB().Create(address); err != nil {
		c.ServerError("failed to create address", err)
		return
	}

	user.BillingAddressID = &address.ID
	if err := c.DB().Save(&user); err != nil {
		c.ServerError("internal server error", err)
		return
	}

	// TODO: should only apply to mock payment processor
	payload := map[string]any{
		"external_session_id":        session.ExternalSessionID,
		"external_payment_method_id": "pm_test_123",
	}
	if err := c.Queue().Enqueue("create_payment_method", payload, "billing_queue"); err != nil {
		c.ServerError("internal server error", err)
		return
	}

	c.OK(nil)
}
