package subscription

import (
	"log"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// Find looks up an active subscription and returns whether to
// create a new one or update the existing one.
func (s *service) Find(input *FindInput) (*FindResult, error) {
	var sub models.Subscription
	err := s.db.FindFirst(&sub,
		"app_id = ? AND user_id = ? AND status = ?",
		input.AppID, input.UserID, "active")

	if err != nil {
		if !dbx.IsRecordNotFoundError(err) {
			return nil, err
		}
		log.Printf("No subscription for app=%d user=%d, creating",
			input.AppID, input.UserID)
		return &FindResult{
			Action: "create",
			AppID:  input.AppID,
			UserID: input.UserID,
			PlanID: input.PlanID,
		}, nil
	}

	return &FindResult{
		Action:         "update",
		SubscriptionID: sub.ID,
		PlanID:         input.PlanID,
		AppID:          input.AppID,
	}, nil
}
