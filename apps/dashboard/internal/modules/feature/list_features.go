package feature

import (
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ListFeatures(c *routerx.Context) {
	isMetered := c.QueryBool("is_metered", false)

	features := []models.Feature{}
	if err := c.DB().List(&features, "app_id = ? AND is_metered = ?", c.AppID(), isMetered); err != nil {
		c.ServerError("Failed to list features", err)
		return
	}

	response := make([]Feature, len(features))
	for i, feature := range features {
		response[i].Set(&feature)
	}

	c.OK(&response)
}
