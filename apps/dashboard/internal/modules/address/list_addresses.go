package address

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ListAddresses(c *routerx.Context) {
	addresses := []models.Address{}
	if err := c.DB().List(&addresses, "app_id = ?", c.AppID()); err != nil {
		c.NotFound("Addresses not found")
		return
	}

	response := make([]apix.Address, len(addresses))
	for i, addr := range addresses {
		response[i].Set(&addr)
	}

	c.OK(response)
}
