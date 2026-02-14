package app_config

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func GetAppConfig(c *routerx.Context) {
	var appConfig models.AppConfig
	if err := c.DB().FindFirst(&appConfig, "app_id = ?", c.AppID()); err != nil {
		c.NotFound("App config not found")
		return
	}

	if appConfig.DefaultConnectionID != 0 {
		var defaultConnection models.Connection
		if err := c.DB().FindForID(appConfig.DefaultConnectionID, &defaultConnection); err != nil {
			c.ServerError("Failed to get default connection", err)
			return
		}

		appConfig.DefaultConnection = defaultConnection
	}

	c.OK(new(apix.AppConfig).Set(&appConfig))
}
