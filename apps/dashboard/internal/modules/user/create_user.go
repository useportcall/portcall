package user

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func CreateUser(c *routerx.Context) {
	var body CreateUserRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	user := models.User{
		PublicID: dbx.GenPublicID("user"),
		AppID:    c.AppID(),
	}

	if body.Email == "" {
		c.BadRequest("Missing or invalid email for user")
		return
	}
	user.Email = body.Email
	user.Name = body.Name

	if err := c.DB().Create(&user); err != nil {
		c.ServerError("Failed to create user", err)
		return
	}

	c.OK(new(User).Set(&user))
}
