package connection

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ListConnections(c *routerx.Context) {
	var connections []models.Connection
	if err := c.DB().ListForAppID(c.AppID(), &connections, nil); err != nil {
		c.ServerError("Failed to list connections", err)
		return
	}

	response := make([]apix.Connection, len(connections))
	for i, connection := range connections {
		response[i].Set(&connection)
	}

	c.OK(response)
}
