package address

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func GetAddress(c *routerx.Context) {
	id := c.Param("id")

	address := &models.Address{}
	if err := c.DB().GetForPublicID(c.AppID(), id, address); err != nil {
		c.NotFound("Address not found")
		return
	}

	c.OK(new(Address).Set(address))
}
