package handlers

import (
	"encoding/json"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

type CreatePaymentMethodPayload struct {
	ExternalSessionID       string `json:"external_session_id"`
	ExternalPaymentMethodID string `json:"external_payment_method_id"`
}

func CreatePaymentMethod(c server.IContext) error {
	var p CreatePaymentMethodPayload
	if err := json.Unmarshal(c.Payload(), &p); err != nil {
		return err
	}

	// ensure checkout session is locked
	checkoutSession := models.CheckoutSession{Status: "pending"}
	if err := c.DB().Update(&checkoutSession, "external_session_id = ? AND status = ?", p.ExternalSessionID, "active"); err != nil {
		return err
	}

	paymentMethod := models.PaymentMethod{}
	paymentMethod.PublicID = dbx.GenPublicID("pm")
	paymentMethod.AppID = checkoutSession.AppID
	paymentMethod.UserID = checkoutSession.UserID
	paymentMethod.ExternalID = p.ExternalPaymentMethodID
	paymentMethod.ExternalType = "card"
	if err := c.DB().FindFirst(&paymentMethod, "external_id = ?", p.ExternalPaymentMethodID); err != nil {
		if !dbx.IsRecordNotFoundError(err) {
			return err
		}

		if err := c.DB().Create(&paymentMethod); err != nil {
			return err
		}
	}

	payload := map[string]any{
		"app_id":  checkoutSession.AppID,
		"user_id": checkoutSession.UserID,
		"plan_id": checkoutSession.PlanID,
	}
	if err := c.Queue().Enqueue("create_subscription", payload, "billing_queue"); err != nil {
		return err
	}

	return nil
}
