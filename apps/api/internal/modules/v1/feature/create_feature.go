package feature

import (
	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type CreateFeatureRequest struct {
	FeatureID string `json:"feature_id" binding:"required"`
	IsMetered bool   `json:"is_metered"`
}

func CreateFeature(c *routerx.Context) {
	var body CreateFeatureRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.BadRequest("Invalid request body")
		return
	}

	if body.FeatureID == "" {
		c.BadRequest("feature_id is required")
		return
	}

	feature := models.Feature{
		PublicID:  body.FeatureID,
		IsMetered: body.IsMetered,
		AppID:     c.AppID(),
	}

	if err := c.DB().Create(&feature); err != nil {
		c.ServerError("Failed to create feature", err)
		return
	}

	c.OK(new(apix.Feature).Set(&feature))
}
