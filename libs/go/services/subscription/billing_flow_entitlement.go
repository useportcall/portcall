package subscription

import (
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/services/entitlement"
)

func syncEntitlements(db dbx.IORM, userID, planID uint) error {
	svc := entitlement.NewService(db)
	start, err := svc.StartUpsert(&entitlement.StartUpsertInput{
		UserID: userID,
		PlanID: planID,
	})
	if err != nil || !start.HasMore {
		return err
	}

	for i := 0; i < len(start.PlanFeatureIDs); i++ {
		_, err := svc.Upsert(&entitlement.UpsertInput{
			UserID: start.UserID,
			Index:  i,
			Values: start.PlanFeatureIDs,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
