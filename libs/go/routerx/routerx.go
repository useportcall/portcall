package routerx

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/useportcall/portcall/libs/go/authx"
	"github.com/useportcall/portcall/libs/go/cryptox"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/logx"
	"github.com/useportcall/portcall/libs/go/qx"
	"github.com/useportcall/portcall/libs/go/storex"
)

type IRouter interface {
	GET(relativePath string, handler HandlerFunc) gin.IRoutes
	POST(relativePath string, handler HandlerFunc) gin.IRoutes
	DELETE(relativePath string, handler HandlerFunc) gin.IRoutes
	Group(relativePath string) *gin.RouterGroup
	Use(middleware HandlerFunc) gin.IRoutes
	Run(addr ...string) error
	NoRoute(handler HandlerFunc)
	LoadHTMLGlob(pattern string)
	SetStore(store storex.IStore)
}

type HandlerFunc func(*Context)

type Context struct {
	*router
	*gin.Context
}

func (c *Context) Body() map[string]any {
	var body map[string]any
	if err := c.BindJSON(&body); err != nil {
		return nil
	}
	return body
}

func (c *Context) QueryBool(key string, defaultValue bool) bool {
	value, exists := c.GetQuery(key)
	if !exists {
		return defaultValue
	}
	return value == "true" || value == "1"
}

func (c *Context) OK(obj any) {
	c.JSON(http.StatusOK, gin.H{"data": obj})
}

func (c *Context) DB() dbx.IORM {
	return c.db
}

func (c *Context) ServerError(message string, err error) {
	log.Println("server error:", message, err.Error())
	c.JSON(http.StatusInternalServerError, gin.H{"error": message})
}

func (c *Context) BadRequest(message string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": message})
}

func (c *Context) NotFound(message string) {
	c.JSON(http.StatusNotFound, gin.H{"error": message})
}

func (c *Context) Unauthorized(message string) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": message})
}

func (c *Context) AuthEmail() string {
	return c.MustGet("auth_email").(string)
}

func (c *Context) AuthClaims() *authx.Claims {
	return c.MustGet("auth_claims").(*authx.Claims)
}

func (c *Context) AppID() uint {
	return c.MustGet("app_id").(uint)
}

func (c *Context) Store() storex.IStore {
	return c.store
}

type router struct {
	instance *gin.Engine
	db       dbx.IORM
	crypto   cryptox.ICrypto
	queue    qx.IQueue
	store    storex.IStore
}

func New(db dbx.IORM, crypto cryptox.ICrypto, q qx.IQueue) IRouter {
	logx.Init()

	r := gin.Default()

	r.Use(gin.Logger(), gin.Recovery())

	r.SetTrustedProxies([]string{"127.0.0.1", "172.17.0.0/16"}) // TODO: review for production

	return &router{instance: r, db: db, crypto: crypto, queue: q}
}

func (r *router) SetStore(store storex.IStore) {
	r.store = store
}

func (r *router) NoRoute(handler HandlerFunc) {
	r.instance.NoRoute(func(c *gin.Context) {
		ctx := &Context{r, c}

		handler(ctx)
	})
}

func (r *router) GET(relativePath string, handler HandlerFunc) gin.IRoutes {
	return r.instance.GET(relativePath, func(c *gin.Context) {
		ctx := &Context{r, c}

		handler(ctx)
	})
}

func (r *router) POST(relativePath string, handler HandlerFunc) gin.IRoutes {
	return r.instance.POST(relativePath, func(c *gin.Context) {
		ctx := &Context{r, c}

		handler(ctx)
	})
}

func (r *router) DELETE(relativePath string, handler HandlerFunc) gin.IRoutes {
	return r.instance.DELETE(relativePath, func(c *gin.Context) {
		ctx := &Context{r, c}

		handler(ctx)
	})
}

func (r *router) Group(relativePath string) *gin.RouterGroup {
	return r.instance.Group(relativePath)
}

func (r *router) Use(middleware HandlerFunc) gin.IRoutes {
	return r.instance.Use(func(c *gin.Context) {
		ctx := &Context{r, c}

		middleware(ctx)
	})
}

func (r *router) LoadHTMLGlob(pattern string) {
	r.instance.LoadHTMLGlob(pattern)
}

func (c *Context) Crypto() cryptox.ICrypto {
	return c.crypto
}

func (c *Context) Queue() qx.IQueue {
	return c.queue
}

func (r *router) Run(addr ...string) error {
	return r.instance.Run(addr...)
}
