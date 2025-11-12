package user

import (
	"strings"

	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ListUsers(c *routerx.Context) {
	appID := c.AppID()

	where := []string{"app_id = ?"}
	args := []any{appID}

	if email := c.Query("email"); email != "" {
		where = append(where, "email = ?")
		args = append(args, email)
	}

	query := strings.Join(where, " AND ")

	conds := []any{query}
	conds = append(conds, args...)

	var users []models.User
	if err := c.DB().List(&users, conds...); err != nil {
		c.ServerError("internal server error")
		return
	}

	response := make([]apix.User, len(users))
	for i, user := range users {
		response[i] = *new(apix.User).Set(&user)
	}

	c.OK(response)
}
