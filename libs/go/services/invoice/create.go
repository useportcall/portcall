package invoice

import (
	"fmt"
	"log"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// Create creates an invoice with all items and totals in a single transaction.
// It looks up the subscription, checks idempotency, builds invoice + items,
// calculates totals, and persists everything atomically.
func (s *service) Create(input *CreateInput) (*CreateResult, error) {
	log.Printf("Processing CreateInvoice for subscription ID %d", input.SubscriptionID)

	sub, err := findSubscription(s.db, input.SubscriptionID)
	if err != nil {
		return nil, err
	}

	skipped, err := isIdempotent(s.db, sub)
	if err != nil {
		return nil, err
	}
	if skipped {
		return &CreateResult{Skipped: true}, nil
	}

	company, count, err := lookupCompanyAndCount(s.db, sub.AppID)
	if err != nil {
		return nil, err
	}

	invoice, err := buildInvoice(sub, company, count)
	if err != nil {
		return nil, err
	}

	subItemIDs, err := listSubscriptionItemIDs(s.db, sub.ID)
	if err != nil {
		return nil, err
	}

	items, err := buildItems(s.db, invoice, sub.ID, subItemIDs)
	if err != nil {
		return nil, err
	}

	applyTotals(invoice, items)

	if err := persistInvoice(s.db, invoice, items); err != nil {
		return nil, err
	}

	return &CreateResult{Invoice: invoice}, nil
}

func findSubscription(db dbx.IORM, id uint) (*models.Subscription, error) {
	var sub models.Subscription
	if err := db.FindForID(id, &sub); err != nil {
		return nil, fmt.Errorf("subscription %d not found: %w", id, err)
	}
	return &sub, nil
}

func isIdempotent(db dbx.IORM, sub *models.Subscription) (bool, error) {
	var existing models.Invoice
	windowStart := invoiceWindowStart(sub)
	err := db.FindFirst(&existing,
		"subscription_id = ? AND status IN (?, ?, ?) AND created_at >= ?",
		sub.ID, "pending", "issued", "paid", windowStart)
	if err == nil {
		log.Printf("Invoice exists for subscription %d in cycle window, skipping", sub.ID)
		return true, nil
	}
	if !dbx.IsRecordNotFoundError(err) {
		return false, err
	}
	return false, nil
}

func lookupCompanyAndCount(db dbx.IORM, appID uint) (*models.Company, int64, error) {
	var company models.Company
	if err := db.FindFirst(&company, "app_id = ?", appID); err != nil {
		return nil, 0, fmt.Errorf("company for app %d: %w", appID, err)
	}
	var count int64
	if err := db.Count(&count, &models.Invoice{}, "app_id = ?", appID); err != nil {
		return nil, 0, err
	}
	return &company, count, nil
}
