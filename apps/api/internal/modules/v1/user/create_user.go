package user

import (
	"log"

	"github.com/useportcall/portcall/apps/api/portcall"
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type CreateUserRequest struct {
	ID    string `json:"id"` // Optional: custom public ID for the user
	Email string `json:"email" binding:"required,email"`
	Name  string `json:"name"`
}

func CreateUser(c *routerx.Context) {
	var body CreateUserRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	if body.Email == "" {
		c.BadRequest("Missing or invalid email for user")
		return
	}

	// Check if custom ID is provided and already exists
	if body.ID != "" {
		var existingUser models.User
		if err := c.DB().FindFirst(&existingUser, "app_id = ? AND public_id = ?", c.AppID(), body.ID); err == nil {
			// User with this ID already exists, return it
			log.Printf("User with ID %s already exists, returning existing user", body.ID)
			c.OK(new(apix.User).Set(&existingUser))
			return
		}
	}

	// Generate public ID or use custom one
	publicID := body.ID
	if publicID == "" {
		publicID = dbx.GenPublicID("user")
	}

	user := models.User{
		PublicID: publicID,
		AppID:    c.AppID(),
		Email:    body.Email,
		Name:     body.Name,
	}

	if err := c.DB().Create(&user); err != nil {
		c.ServerError("Failed to create user", err)
		return
	}

	// increment df max subscriptions by 1
	if err := c.Queue().Enqueue("df_increment", map[string]any{"user_id": c.PublicAppID(), "feature": portcall.Features.NumberOfUsers, "is_test": !c.IsLive()}, "billing_queue"); err != nil {
		log.Printf("Error enqueueing df_increment: %v", err)
		c.ServerError("error updating feature usage", err)
		return
	}

	c.OK(new(apix.User).Set(&user))
}
