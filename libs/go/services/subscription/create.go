package subscription

import (
	"log"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// Create creates a new subscription with all items in a single
// transaction. It also finds and sets the rollback plan for paid plans.
func (s *service) Create(input *CreateInput) (*CreateResult, error) {
	log.Printf("Creating subscription app=%d user=%d plan=%d",
		input.AppID, input.UserID, input.PlanID)

	plan, err := findPlan(s.db, input.PlanID)
	if err != nil {
		return nil, err
	}

	billingAddrID, err := findBillingAddr(s.db, input.UserID)
	if err != nil {
		return nil, err
	}

	nextReset, err := calcNextReset(plan.Interval, time.Now())
	if err != nil {
		return nil, err
	}

	sub := buildSubscription(plan, input.UserID, billingAddrID, nextReset)

	if !plan.IsFree {
		setRollbackPlan(s.db, plan.AppID, sub)
	}

	planItemIDs, err := listPlanItemIDs(s.db, plan.ID)
	if err != nil {
		return nil, err
	}

	itemCount, err := persistSubscriptionWithItems(s.db, sub, plan, planItemIDs)
	if err != nil {
		return nil, err
	}

	log.Printf("Subscription %d created", sub.ID)
	return &CreateResult{
		Subscription: sub,
		UserID:       input.UserID,
		PlanID:       plan.ID,
		ItemCount:    itemCount,
	}, nil
}

func findPlan(db dbx.IORM, id uint) (*models.Plan, error) {
	var plan models.Plan
	if err := db.FindForID(id, &plan); err != nil {
		return nil, err
	}
	return &plan, nil
}

func findBillingAddr(db dbx.IORM, userID uint) (*uint, error) {
	var user models.User
	if err := db.FindForID(userID, &user); err != nil {
		return nil, err
	}
	return user.BillingAddressID, nil
}

func calcNextReset(interval string, now time.Time) (time.Time, error) {
	next, err := NextReset(now, interval, now)
	if err != nil {
		return time.Time{}, err
	}
	return *next, nil
}

func buildSubscription(
	plan *models.Plan, userID uint,
	billingAddrID *uint, nextReset time.Time,
) *models.Subscription {
	return &models.Subscription{
		PublicID:             dbx.GenPublicID("sub"),
		AppID:                plan.AppID,
		Status:               "active",
		Currency:             plan.Currency,
		PlanID:               &plan.ID,
		UserID:               userID,
		BillingAddressID:     billingAddrID,
		BillingInterval:      plan.Interval,
		BillingIntervalCount: plan.IntervalCount,
		InvoiceDueByDays:     plan.InvoiceDueByDays,
		NextResetAt:          nextReset,
	}
}
