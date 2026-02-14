package payment

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func recoverSubscriptionStatus(db dbx.IORM, subscriptionID *uint) error {
	if subscriptionID == nil {
		return nil
	}

	var sub models.Subscription
	if err := db.FindForID(*subscriptionID, &sub); err != nil {
		return err
	}
	if sub.Status != "past_due" {
		return nil
	}

	sub.Status = "active"
	return db.Save(&sub)
}
