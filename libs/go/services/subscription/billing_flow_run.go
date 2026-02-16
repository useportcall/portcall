package subscription

// RunCreateFlow runs the create-subscription billing flow.
func RunCreateFlow(input *BillingFlowInput) error {
	return runCreateFlow(input)
}

// RunUpdateFlow runs the update-subscription billing flow.
func RunUpdateFlow(input *BillingFlowInput, subscriptionID uint) error {
	return runUpdateFlow(input, subscriptionID)
}
