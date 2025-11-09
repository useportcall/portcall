package app

import (
	"github.com/useportcall/portcall/apps/dashboard/internal/routerx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func GetApp(c *routerx.Context) {

	var app models.App
	if err := c.DB().FindFirst(&app, "id = ?", c.AppID()); err != nil {
		c.NotFound("App not found")
		return
	}

	c.OK(new(App).Set(&app))
}
