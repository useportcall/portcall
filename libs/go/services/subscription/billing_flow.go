package subscription

import (
	"fmt"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// RunBillingFlow applies create/update subscription billing orchestration.
func RunBillingFlow(input *BillingFlowInput) (*BillingFlowResult, error) {
	if input == nil || input.DB == nil {
		return nil, fmt.Errorf("billing flow requires db")
	}
	if input.Crypto == nil {
		return nil, fmt.Errorf("billing flow requires crypto")
	}

	var existing models.Subscription
	err := input.DB.FindFirst(
		&existing,
		"app_id = ? AND user_id = ? AND status = ?",
		input.AppID,
		input.UserID,
		"active",
	)
	if err == nil {
		if existing.PlanID != nil && *existing.PlanID == input.PlanID {
			return &BillingFlowResult{CreatedNew: false}, nil
		}
		if err := runUpdateFlow(input, existing.ID); err != nil {
			return nil, err
		}
		return &BillingFlowResult{CreatedNew: false}, nil
	}
	if !dbx.IsRecordNotFoundError(err) {
		return nil, err
	}
	if err := runCreateFlow(input); err != nil {
		return nil, err
	}
	return &BillingFlowResult{CreatedNew: true}, nil
}
