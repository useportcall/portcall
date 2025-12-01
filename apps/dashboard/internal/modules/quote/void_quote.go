package quote

import (
	"time"

	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func VoidQuote(c *routerx.Context) {
	id := c.Param("id")

	var quote models.Quote
	if err := c.DB().GetForPublicID(c.AppID(), id, &quote); err != nil {
		c.NotFound("Quote not found")
		return
	}

	now := time.Now()
	quote.VoidedAt = &now
	quote.Status = "voided"

	if err := c.DB().Save(&quote); err != nil {
		c.ServerError("Failed to update quote", err)
		return
	}

	response := new(apix.Quote)

	c.OK(response.Set(&quote))
}
