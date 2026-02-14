package handlers

import (
	"net/http"

	"github.com/useportcall/portcall/libs/go/routerx"
)

// QuoteSuccess renders the success page after quote submission (non-direct checkout)
func QuoteSuccess(c *routerx.Context) {
	id := c.Param("id")
	state := c.DefaultQuery("state", "accepted")
	c.HTML(http.StatusOK, "success.html", map[string]interface{}{
		"ID":    id,
		"State": state,
	})
}
