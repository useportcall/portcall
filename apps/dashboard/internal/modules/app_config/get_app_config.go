package app_config

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func GetAppConfig(c *routerx.Context) {
	var appConfig models.AppConfig
	if err := c.DB().FindFirst(&appConfig, "app_id = ?", c.AppID()); err != nil {
		c.NotFound("App config not found")
		return
	}

	c.OK(new(AppConfig).Set(&appConfig))
}
