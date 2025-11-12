package account

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func GetAccount(c *routerx.Context) {
	email := c.AuthEmail()

	account := new(models.Account)
	if err := c.DB().FindFirst(account, "email = ?", email); err != nil {
		if !dbx.IsRecordNotFoundError(err) {
			c.ServerError("Internal server error")
			return
		}

		account.Email = email
		if err := c.DB().Create(account); err != nil {
			c.ServerError("Internal server error")
			return
		}
	}

	c.OK(new(Account).Set(account))
}
