package account

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func GetAccount(c *routerx.Context) {
	email := c.AuthEmail()

	var account models.Account
	if err := c.DB().FindFirst(&account, "email = ?", email); err != nil {
		c.NotFound("Account not found")
		return
	}

	c.OK(new(Account).Set(&account))
}
