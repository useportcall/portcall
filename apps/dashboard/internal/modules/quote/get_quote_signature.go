package quote

import (
	"fmt"
	"net/http"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func GetQuoteSignature(c *routerx.Context) {
	id := c.Param("id")

	var quote models.Quote
	if err := c.DB().GetForPublicID(c.AppID(), id, &quote); err != nil {
		c.NotFound("Quote not found")
		return
	}
	if quote.Status != "accepted" {
		c.BadRequest("Quote is not signed yet")
		return
	}

	data, err := c.Store().GetFromSignatureBucket(quote.PublicID, c)
	if err != nil {
		c.NotFound("Signed quote not found")
		return
	}

	c.Header("Content-Type", "image/png")
	c.Header("Content-Disposition",
		fmt.Sprintf("attachment; filename=\"%s-signature.png\"", quote.PublicID))
	c.Data(http.StatusOK, "image/png", data)
}
