package subscription

func runCreateFlow(input *BillingFlowInput) error {
	svc := NewService(input.DB)
	result, err := svc.Create(&CreateInput{
		AppID:  input.AppID,
		UserID: input.UserID,
		PlanID: input.PlanID,
	})
	if err != nil {
		return err
	}

	invoiceID, err := createStandardInvoice(input.DB, result.Subscription.ID, result.ItemCount)
	if err != nil {
		return err
	}
	if err := syncEntitlements(input.DB, result.UserID, result.PlanID); err != nil {
		return err
	}
	if invoiceID == 0 {
		return nil
	}
	return payAndResolveInvoice(input, invoiceID)
}
