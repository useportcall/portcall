package payment_link

import (
	"errors"
	"net/http"
	"strings"

	"github.com/useportcall/portcall/libs/go/routerx"
	pl "github.com/useportcall/portcall/libs/go/services/payment_link"
)

const paymentLinkTokenHeader = "X-Payment-Link-Token"

func RedeemPaymentLink(c *routerx.Context) {
	result, err := pl.NewService(c.DB(), c.Crypto()).Redeem(&pl.RedeemInput{
		ID:    strings.TrimSpace(c.Param("id")),
		Token: strings.TrimSpace(c.GetHeader(paymentLinkTokenHeader)),
	})
	if err != nil {
		var ve *pl.ValidationError
		if errors.As(err, &ve) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]any{"error": ve.Message})
			return
		}
		c.ServerError("error redeeming payment link", err)
		return
	}
	c.OK(map[string]any{"checkout_url": result.CheckoutURL})
}
