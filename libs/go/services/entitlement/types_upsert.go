package entitlement

// StartUpsertInput is the input for starting an entitlement upsert workflow.
type StartUpsertInput struct {
	UserID uint `json:"user_id"`
	PlanID uint `json:"plan_id"`
}

// StartUpsertResult is the result of starting an entitlement upsert workflow.
type StartUpsertResult struct {
	UserID         uint
	PlanFeatureIDs []uint
	HasMore        bool
}

// UpsertInput is the input for upserting a single entitlement.
type UpsertInput struct {
	UserID uint   `json:"user_id"`
	Index  int    `json:"index"`
	Values []uint `json:"values"`
}

// UpsertResult is the result of upserting a single entitlement.
type UpsertResult struct {
	UserID         uint
	PlanFeatureIDs []uint
	NextIndex      int
	HasMore        bool
}
