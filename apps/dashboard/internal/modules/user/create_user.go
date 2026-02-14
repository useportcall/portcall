package user

import (
	"log"

	modulebilling "github.com/useportcall/portcall/apps/dashboard/internal/modules/billing"
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type CreateUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name"`
}

func CreateUser(c *routerx.Context) {
	// check if df number of users threshold reached
	isBillingExempt := c.GetBool("is_billing_exempt")
	if !isBillingExempt {
		entitlement, err := modulebilling.CheckFeatureEntitlement(c, modulebilling.DFFeatures.NumberOfUsers)
		if err != nil {
			if !c.IsLive() {
				log.Printf("Skipping number_of_users entitlement check in test mode for app %s: %v", c.PublicAppID(), err)
			} else {
				log.Printf("Error checking number of users entitlement: %v", err)
				c.ServerError("error checking number of users entitlement", err)
				return
			}
		} else if entitlement.Enabled == false {
			log.Printf("Max number of users limit reached for app %s", c.PublicAppID())
			c.BadRequest("max number of users limit reached")
			return
		}
	}

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

	// increment df max subscriptions by 1
	if !isBillingExempt {
		if err := c.Queue().Enqueue("df_increment", map[string]any{"user_id": c.PublicAppID(), "feature": modulebilling.DFFeatures.NumberOfUsers, "is_test": !c.IsLive()}, "billing_queue"); err != nil {
			log.Printf("Error enqueueing df_increment: %v", err)
			c.ServerError("error updating feature usage", err)
			return
		}
	}

	c.OK(new(apix.User).Set(&user))
}
