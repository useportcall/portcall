package entitlement

// ResetAllInput is the input for resetting all entitlements for a user.
type ResetAllInput struct {
	UserID uint `json:"user_id"`
}

// ResetAllResult is the result of resetting all entitlements.
type ResetAllResult struct {
	ResetCount int
}
