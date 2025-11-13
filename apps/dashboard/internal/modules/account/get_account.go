package account

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func GetAccount(c *routerx.Context) {
	claims := c.AuthClaims()

	account := new(models.Account)
	if err := c.DB().FindFirst(account, "email = ?", claims.Email); err != nil {
		if !dbx.IsRecordNotFoundError(err) {
			c.ServerError("Internal server error", err)
			return
		}

		account.Email = claims.Email

		if claims.GivenName != nil {
			account.FirstName = *claims.GivenName
		}

		if claims.FamilyName != nil {
			account.LastName = *claims.FamilyName
		}

		if err := c.DB().Create(account); err != nil {
			c.ServerError("Internal server error", err)
			return
		}
	}

	c.OK(new(Account).Set(account))
}
