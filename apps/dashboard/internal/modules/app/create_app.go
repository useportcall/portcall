package app

import (
	"github.com/useportcall/portcall/apps/dashboard/internal/utils"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func CreateApp(c *routerx.Context) {
	// find account for auth email, TODO: abstract FindOrCreate?
	var account models.Account
	if err := c.DB().FindFirst(&account, "email = ?", c.AuthEmail()); err != nil {
		if dbx.IsRecordNotFoundError(err) {
			account = models.Account{}
			account.Email = c.AuthEmail()
			if err := c.DB().Create(&account); err != nil {
				c.ServerError("Failed to create account")
				return
			}
		} else {
			c.ServerError("Failed to find account")
			return
		}
	}

	var body CreateAppRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	if body.Name == "" {
		c.BadRequest("App name is required")
		return
	}

	app := new(models.App)
	app.PublicID = utils.GenPublicID("app")
	app.AccountID = account.ID
	app.Name = body.Name
	app.PublicApiKey = utils.GenPublicID("pk")
	if err := c.DB().Create(app); err != nil {
		c.ServerError("Failed to create app")
		return
	}

	// TODO: explore async job for setup tasks

	// company address
	address := new(models.Address)
	address.AppID = app.ID
	address.Line1 = "123 Main St"
	address.City = "Anytown"
	address.State = "CA"
	address.PostalCode = "12345"
	address.Country = "USA"
	if err := c.DB().Create(address); err != nil {
		c.ServerError("Failed to create address")
		return
	}

	// company
	company := new(models.Company)
	company.AppID = app.ID
	company.BillingAddressID = address.ID
	company.Name = "Default Company"
	company.Email = "default@example.com"
	company.VATNumber = "123456789"
	if err := c.DB().Create(company); err != nil {
		c.ServerError("Failed to create company")
		return
	}

	// connection
	connection := new(models.Connection)
	connection.AppID = app.ID
	connection.Source = "local"
	connection.Name = "Local Payment Provider"
	if err := c.DB().Create(connection); err != nil {
		c.ServerError("Failed to create connection")
		return
	}

	// app config
	appConfig := new(models.AppConfig)
	appConfig.AppID = app.ID
	appConfig.DefaultConnectionID = connection.ID
	if err := c.DB().Create(appConfig); err != nil {
		c.ServerError("Failed to create app config")
		return
	}

	c.OK(new(App).Set(app))
}
