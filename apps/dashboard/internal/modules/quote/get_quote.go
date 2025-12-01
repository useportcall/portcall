package quote

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func GetQuote(c *routerx.Context) {
	id := c.Param("id")
	response := new(apix.Quote)

	var quote models.Quote
	if err := c.DB().GetForPublicID(c.AppID(), id, &quote); err != nil {
		c.NotFound("Quote not found")
		return
	}

	var plan models.Plan
	if err := c.DB().FindForID(quote.PlanID, &plan); err != nil {
		c.NotFound("Plan not found")
		return
	}

	quote.Plan = plan

	if quote.UserID != nil {
		var user models.User
		if err := c.DB().FindForID(*quote.UserID, &user); err != nil {
			c.NotFound("User not found")
			return
		}
		quote.User = &user
	}

	c.OK(response.Set(&quote))
}
