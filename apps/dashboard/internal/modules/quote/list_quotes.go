package quote

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ListQuotes(c *routerx.Context) {
	var q []any
	{
		q = []any{"app_id = ?", c.AppID()}
	}

	quotes := []models.Quote{}
	if err := c.DB().ListWithOrder(&quotes, "created_at DESC", q...); err != nil {
		c.ServerError("Failed to list quotes", err)
		return
	}

	response := make([]apix.Quote, len(quotes))
	for i, quote := range quotes {
		response[i].Set(&quote)
	}

	c.OK(response)
}
