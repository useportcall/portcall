package subscription

import (
	"strings"

	"github.com/useportcall/portcall/libs/go/apix"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func ListSubscriptions(c *routerx.Context) {
	appID := c.AppID()

	where := []string{"app_id = ?"}
	args := []any{appID}

	if userID := c.Query("user_id"); userID != "" {
		var user models.User
		if err := c.DB().GetForPublicID(c.AppID(), userID, &user); err != nil {
			c.ServerError("Internal server error")
			return
		}

		where = append(where, "user_id = ?")
		args = append(args, user.ID)
	}

	if status := c.Query("status"); status != "" {
		where = append(where, "status = ?")
		args = append(args, status)
	}

	query := strings.Join(where, " AND ")

	conds := []any{query}
	conds = append(conds, args...)

	subscriptions := []models.Subscription{}
	if err := c.DB().List(&subscriptions, conds...); err != nil {
		if !dbx.IsRecordNotFoundError(err) {
			c.ServerError("Internal server error")
			return
		}
	}

	response := make([]apix.Subscription, len(subscriptions))
	for i, subscription := range subscriptions {
		response[i] = *new(apix.Subscription).Set(&subscription)
	}

	c.OK(response)
}
