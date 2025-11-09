package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/useportcall/portcall/libs/go/authx"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func Auth(db dbx.IORM) routerx.HandlerFunc {
	client := authx.New()

	return func(c *routerx.Context) {
		email, err := client.Validate(c.Request.Context(), c.Request.Header)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"access": "unauthorized"})
			c.Abort()
			return
		}

		c.Set("auth_email", email)

		if c.Param("app_id") == "" {
			c.Next()
			return
		}

		var app models.App
		if err := db.FindFirst(&app, "public_id = ?", c.Param("app_id")); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "App not found"})
			c.Abort()
			return
		}

		c.Set("app_id", app.ID)

		c.Next()
	}
}
