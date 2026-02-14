package payment_link

import (
	"errors"
	"time"

	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/routerx"
	pl "github.com/useportcall/portcall/libs/go/services/payment_link"
)

type CreatePaymentLinkRequest struct {
	PlanID      string     `json:"plan_id"`
	UserID      string     `json:"user_id"`
	UserEmail   string     `json:"user_email"`
	UserName    string     `json:"user_name"`
	CancelURL   string     `json:"cancel_url"`
	RedirectURL string     `json:"redirect_url"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

func CreatePaymentLink(c *routerx.Context) {
	var body CreatePaymentLinkRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("invalid request body")
		return
	}
	result, err := pl.NewService(c.DB(), c.Crypto()).Create(&pl.CreateInput{
		AppID:                 c.AppID(),
		PlanID:                body.PlanID,
		UserID:                body.UserID,
		UserEmail:             body.UserEmail,
		UserName:              body.UserName,
		CancelURL:             body.CancelURL,
		RedirectURL:           body.RedirectURL,
		ExpiresAt:             body.ExpiresAt,
		RequireBillingAddress: true,
	})
	if err != nil {
		var ve *pl.ValidationError
		if errors.As(err, &ve) {
			c.BadRequest(ve.Message)
			return
		}
		c.ServerError("error creating payment link", err)
		return
	}
	response := new(apix.PaymentLink).Set(result.PaymentLink)
	if result.URL != "" {
		response.URL = result.URL
	}
	c.OK(response)
}
