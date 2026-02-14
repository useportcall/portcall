package app

import (
	"log"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func CreateApp(c *routerx.Context) {
	var account models.Account
	if err := c.DB().FindFirst(&account, "email = ?", c.AuthEmail()); err != nil {
		if !dbx.IsRecordNotFoundError(err) {
			c.ServerError("Failed to find account", err)
			return
		}
		account.Email = c.AuthEmail()
		if err := c.DB().Create(&account); err != nil {
			c.ServerError("Failed to create account", err)
			return
		}
		log.Printf("[ACCOUNT_CREATED] account_id=%d source=dashboard_create_app", account.ID)
	}

	var count int64
	if err := c.DB().Count(&count, &models.App{}, "account_id = ?", account.ID); err != nil {
		c.ServerError("Failed to check existing apps", err)
		return
	}
	if count > 0 {
		c.BadRequest("You can only create one project")
		return
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

	var result App
	var appsToRegister []*models.App
	if err := c.DB().Txn(func(txn dbx.IORM) error {
		testApp, err := createAppInstance(txn, account.ID, body.Name, false)
		if err != nil {
			return err
		}
		liveApp, err := createAppInstance(txn, account.ID, body.Name, true)
		if err != nil {
			return err
		}
		appsToRegister = append(appsToRegister, testApp, liveApp)
		result = *new(App).Set(testApp)
		return nil
	}); err != nil {
		c.ServerError("Failed to create app", err)
		return
	}

	sendAccountSignupNotification(account.Email, body.Name, appsToRegister)
	for _, app := range appsToRegister {
		log.Printf("Enqueueing df registration for app %s", app.PublicID)
		if err := c.Queue().Enqueue("df_create_user", map[string]any{"app_id": app.ID}, "billing_queue"); err != nil {
			log.Printf("Error enqueueing df_create_user for app %s: %v", app.PublicID, err)
		}
	}

	c.OK(result)
}
