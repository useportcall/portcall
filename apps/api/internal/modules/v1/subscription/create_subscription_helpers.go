package subscription

import (
	"errors"
	"log"

	"github.com/useportcall/portcall/apps/api/portcall"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/dogfoodx"
	"github.com/useportcall/portcall/libs/go/routerx"
	subsvc "github.com/useportcall/portcall/libs/go/services/subscription"
)

var errMaxSubscriptionsReached = errors.New("max subscriptions limit reached")

func checkMaxSubscriptions(c *routerx.Context) error {
	if c.GetBool("is_billing_exempt") {
		return nil
	}
	entitlement, err := portcall.New().Check.MaxSubscriptions.For(c, c.PublicAppID())
	if err != nil {
		c.ServerError("error checking subscription entitlement", err)
		return err
	}
	if !entitlement.Enabled {
		c.BadRequest("max subscriptions limit reached")
		return errMaxSubscriptionsReached
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

func runBillingFlow(c *routerx.Context, userID, planID uint) (bool, error) {
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

func recordDogfoodUsage(c *routerx.Context, user *models.User) {
	var app models.App
	if err := c.DB().FindForID(user.AppID, &app); err != nil {
		log.Printf("Warning: failed loading app for usage record: %v", err)
		return
	}
	pc := portcall.New()
	if err := pc.Record.MaxSubscriptions.Increment(c, app.PublicID); err != nil {
		log.Printf("Warning: failed to increment max subscriptions: %v", err)
	}
	if err := dogfoodx.IncrementSubscriptionUsage(c.DB(), user.AppID); err != nil {
		log.Printf("Warning: failed to record subscription meter event: %v", err)
	}
	if err := c.Queue().Enqueue("df_increment", map[string]any{"user_id": app.PublicID, "feature": portcall.Features.MaxSubscriptions, "is_test": !app.IsLive}, "billing_queue"); err != nil {
		log.Printf("Warning: failed to enqueue df_increment: %v", err)
	}
}
