package subscription

import "log"

// FindResets finds all subscription IDs that need to be reset.
func (s *service) FindResets() (*FindResetsResult, error) {
	var ids []uint
	if err := s.db.ListIDs(
		"subscriptions",
		&ids,
		"(status = ? OR status = ?) AND next_reset_at <= CURRENT_TIMESTAMP",
		"active",
		"canceled",
	); err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		log.Println("No subscriptions to reset found")
	}

	return &FindResetsResult{SubscriptionIDs: ids}, nil
}
