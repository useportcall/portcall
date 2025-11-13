package app

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ListApps(c *routerx.Context) {
	var account models.Account
	if err := c.DB().FindFirst(&account, "email = ?", c.AuthEmail()); err != nil {
		if dbx.IsRecordNotFoundError(err) {
			c.OK([]App{})
			return
		}

		c.ServerError("Failed to list apps", err)
		return
	}

	apps := []models.App{}
	if err := c.DB().List(&apps, "account_id = ?", account.ID); err != nil {
		c.ServerError("Failed to list apps", err)
		return
	}

	response := make([]App, len(apps))
	for i, app := range apps {
		response[i].Set(&app)
	}

	c.OK(response)
}
