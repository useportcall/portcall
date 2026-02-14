package on_subscription_reset

import (
	"github.com/useportcall/portcall/libs/go/qx/server"
	"github.com/useportcall/portcall/libs/go/services/entitlement"
)

func syncEntitlements(c server.IContext, userID uint, planID *uint) error {
	svc := entitlement.NewService(c.DB())
	if _, err := svc.ResetAll(&entitlement.ResetAllInput{UserID: userID}); err != nil {
		return err
	}
	if planID == nil {
		return nil
	}

	start, err := svc.StartUpsert(&entitlement.StartUpsertInput{
		UserID: userID,
		PlanID: *planID,
	})
	if err != nil {
		return err
	}
	for i := 0; i < len(start.PlanFeatureIDs); i++ {
		if _, err := svc.Upsert(&entitlement.UpsertInput{
			UserID: userID,
			Index:  i,
			Values: start.PlanFeatureIDs,
		}); err != nil {
			return err
		}
	}
	return nil
}
