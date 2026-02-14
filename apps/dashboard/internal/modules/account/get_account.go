package account

import (
	"strings"

	"github.com/useportcall/portcall/libs/go/apix"
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

		account.FirstName = claimValue(claims.GivenName)
		account.LastName = claimValue(claims.FamilyName)

		if err := c.DB().Create(account); err != nil {
			c.ServerError("Internal server error", err)
			return
		}
	}

	givenName := claimValue(claims.GivenName)
	familyName := claimValue(claims.FamilyName)

	shouldUpdate := false
	if account.FirstName == "" && givenName != "" {
		account.FirstName = givenName
		shouldUpdate = true
	}

	if account.LastName == "" && familyName != "" {
		account.LastName = familyName
		shouldUpdate = true
	}

	if shouldUpdate {
		if err := c.DB().Save(account); err != nil {
			c.ServerError("Internal server error", err)
			return
		}
	}

	c.OK(new(apix.Account).Set(account))
}

func claimValue(value *string) string {
	if value == nil {
		return ""
	}

	return strings.TrimSpace(*value)
}
