package entitlement

// StartUpsert initiates the entitlement upsert workflow by listing all
// plan features and returning the IDs for individual upsert tasks.
func (s *service) StartUpsert(input *StartUpsertInput) (*StartUpsertResult, error) {
	var ids []uint
	if err := s.db.ListIDs("plan_features", &ids, "plan_id = ?", input.PlanID); err != nil {
		return nil, err
	}
	if err := removeMissingEntitlements(s, input.UserID, ids); err != nil {
		return nil, err
	}

	return &StartUpsertResult{
		UserID:         input.UserID,
		PlanFeatureIDs: ids,
		HasMore:        len(ids) > 0,
	}, nil
}
