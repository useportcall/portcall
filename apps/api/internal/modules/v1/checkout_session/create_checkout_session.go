package checkout_session

import (
	"errors"

	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/routerx"
	cs "github.com/useportcall/portcall/libs/go/services/checkout_session"
)

// CreateCheckoutSessionRequest is the HTTP request body.
type CreateCheckoutSessionRequest struct {
	PlanID      string `json:"plan_id"`
	UserID      string `json:"user_id"`
	CancelURL   string `json:"cancel_url"`
	RedirectURL string `json:"redirect_url"`
}

// CreateCheckoutSession handles POST /v1/checkout-sessions.
func CreateCheckoutSession(c *routerx.Context) {
	var body CreateCheckoutSessionRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("invalid request body")
		return
	}

	svc := cs.NewService(c.DB(), c.Crypto())
	result, err := svc.Create(&cs.CreateInput{
		AppID:       c.AppID(),
		PlanID:      body.PlanID,
		UserID:      body.UserID,
		CancelURL:   body.CancelURL,
		RedirectURL: body.RedirectURL,
	})
	if err != nil {
		var ve *cs.ValidationError
		if errors.As(err, &ve) {
			c.BadRequest(ve.Message)
			return
		}
		c.ServerError("error creating checkout session", err)
		return
	}

	response := new(apix.CheckoutSession).Set(result.Session)
	if result.CheckoutURL != "" {
		response.URL = result.CheckoutURL
	}
	c.OK(response)
}
