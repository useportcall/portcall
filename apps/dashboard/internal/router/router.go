package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/useportcall/portcall/apps/dashboard/internal/middleware"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/qx"
)

func Init(db dbx.IORM, crypto cryptox.ICrypto, q qx.IQueue) *gin.Engine {
	r := gin.Default() // TODO: replace defaults

	r.SetTrustedProxies([]string{"127.0.0.1", "172.17.0.0/16"}) // TODO: review for production

	r.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "pong"}) })

	api := r.Group("/api")

	api.Use(middleware.Auth(db))

	return r
}
