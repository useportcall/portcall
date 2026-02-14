package address

import (
	"strings"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type UpdateCheckoutSessionAddressRequest struct {
	Line1      string `json:"line1" binding:"required,max=200"`
	Line2      string `json:"line2" binding:"max=200"`
	City       string `json:"city" binding:"required,max=120"`
	State      string `json:"state" binding:"max=120"`
	PostalCode string `json:"postal_code" binding:"required,max=32"`
	Country    string `json:"country" binding:"required,len=2,uppercase"`
}

func UpdateCheckoutSessionAddress(c *routerx.Context, session *models.CheckoutSession) {

	var body UpdateCheckoutSessionAddressRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("invalid request payload")
		return
	}

	var user models.User
	if err := c.DB().FindForID(session.UserID, &user); err != nil {
		c.ServerError("internal server error", err)
		return
	}

	country := strings.ToUpper(strings.TrimSpace(body.Country))
	address := &models.Address{
		PublicID:   dbx.GenPublicID("addr"),
		AppID:      session.AppID,
		Line1:      strings.TrimSpace(body.Line1),
		Line2:      strings.TrimSpace(body.Line2),
		City:       strings.TrimSpace(body.City),
		State:      strings.TrimSpace(body.State),
		PostalCode: strings.TrimSpace(body.PostalCode),
		Country:    country,
	}

	if err := c.DB().Txn(func(tx dbx.IORM) error {
		if err := tx.Create(address); err != nil {
			return err
		}
		user.BillingAddressID = &address.ID
		if err := tx.Save(&user); err != nil {
			return err
		}
		session.BillingAddressID = &address.ID
		return tx.Save(session)
	}); err != nil {
		c.ServerError("failed to update address", err)
		return
	}

	c.OK(nil)
}
