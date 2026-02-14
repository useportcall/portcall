package subscription

import (
	"errors"
	"log"

	modulebilling "github.com/useportcall/portcall/apps/dashboard/internal/modules/billing"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
	subsvc "github.com/useportcall/portcall/libs/go/services/subscription"
)

var errQuotaExceeded = errors.New("max subscriptions limit reached")

func checkQuota(c *routerx.Context) error {
	if c.GetBool("is_billing_exempt") {
		return nil
	}
	entitlement, err := modulebilling.CheckFeatureEntitlement(c, modulebilling.DFFeatures.MaxSubscriptions)
	if err != nil {
		c.ServerError("error checking subscription entitlement", err)
		return err
	}
	if !entitlement.Enabled {
		c.BadRequest("max subscriptions limit reached")
		return errQuotaExceeded
	}
	return nil
}

func hasPaymentMethod(c *routerx.Context, userID uint) (bool, error) {
	var pm models.PaymentMethod
	err := c.DB().FindFirst(&pm, "app_id = ? AND user_id = ?", c.AppID(), userID)
	if err == nil {
		return true, nil
	}
	if dbx.IsRecordNotFoundError(err) {
		return false, nil
	}
	return false, err
}

func runBillingSubscriptionFlow(c *routerx.Context, userID, planID uint) (bool, error) {
	result, err := subsvc.RunBillingFlow(&subsvc.BillingFlowInput{
		DB:     c.DB(),
		Crypto: c.Crypto(),
		Queue:  c.Queue(),
		AppID:  c.AppID(),
		UserID: userID,
		PlanID: planID,
	})
	if err != nil {
		return false, err
	}
	return result.CreatedNew, nil
}

func recordSubscriptionUsage(c *routerx.Context) error {
	if err := c.Queue().Enqueue("df_increment", map[string]any{
		"user_id": c.PublicAppID(),
		"feature": modulebilling.DFFeatures.MaxSubscriptions,
		"is_test": !c.IsLive(),
	}, "billing_queue"); err != nil {
		log.Printf("Error enqueueing df_increment: %v", err)
		return err
	}
	return nil
}
