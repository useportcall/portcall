package subscription

import "github.com/useportcall/portcall/libs/go/dbx"

func findRollbackPlanID(db dbx.IORM, appID uint, current *uint) (*uint, error) {
	ids, err := listPlanIDs(db, "app_id = ? AND is_free = ?", appID, true)
	if err != nil {
		return nil, err
	}
	if id := firstDifferent(ids, current); id != nil {
		return id, nil
	}

	ids, err = listPlanIDs(db, "app_id = ?", appID)
	if err != nil {
		return nil, err
	}
	return cheapestFixedPlan(db, ids, current)
}

func listPlanIDs(db dbx.IORM, query string, args ...any) ([]uint, error) {
	ids := []uint{}
	conds := []any{query}
	conds = append(conds, args...)
	if err := db.ListIDs("plans", &ids, conds...); err != nil {
		return nil, err
	}
	return ids, nil
}

func firstDifferent(ids []uint, current *uint) *uint {
	for _, id := range ids {
		if current != nil && id == *current {
			continue
		}
		return &id
	}
	return nil
}
