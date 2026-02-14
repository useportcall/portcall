package invoice

import (
	"fmt"
	"log"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func listSubscriptionItemIDs(db dbx.IORM, subscriptionID uint) ([]uint, error) {
	var ids []uint
	if err := db.ListIDs("subscription_items", &ids, "subscription_id = ?", subscriptionID); err != nil {
		return nil, fmt.Errorf("listing items for subscription %d: %w", subscriptionID, err)
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("no subscription items for subscription %d", subscriptionID)
	}
	return ids, nil
}

func buildItems(
	db dbx.IORM, inv *models.Invoice, subscriptionID uint, subItemIDs []uint,
) ([]models.InvoiceItem, error) {
	items := make([]models.InvoiceItem, 0, len(subItemIDs))
	for _, siID := range subItemIDs {
		var si models.SubscriptionItem
		if err := db.FindForID(siID, &si); err != nil {
			return nil, fmt.Errorf("subscription item %d not found: %w", siID, err)
		}
		items = append(items, buildItem(db, inv, subscriptionID, &si))
	}
	return items, nil
}

func buildItem(
	db dbx.IORM, inv *models.Invoice, subscriptionID uint, si *models.SubscriptionItem,
) models.InvoiceItem {
	total := calculateItemTotal(si.PricingModel, si.UnitAmount, si.Quantity, int64(si.Usage), si.Tiers)
	if meteredTotal, ok := meterTotalForItem(db, subscriptionID, si); ok {
		total = meteredTotal
	}
	return models.InvoiceItem{
		PublicID:           dbx.GenPublicID("ii"),
		AppID:              si.AppID,
		InvoiceID:          inv.ID,
		Quantity:           si.Quantity,
		SubscriptionItemID: si.ID,
		Title:              si.Title,
		Description:        si.Description,
		PricingModel:       si.PricingModel,
		Total:              total,
		Amount:             getItemUnitAmount(*si),
	}
}

func applyTotals(invoice *models.Invoice, items []models.InvoiceItem) {
	var subtotal int64
	for _, item := range items {
		subtotal += item.Total
	}
	invoice.SubTotal = subtotal
	invoice.TaxAmount = 0
	invoice.DiscountAmount = calculateDiscount(subtotal, invoice.DiscountPct)
	invoice.Total = subtotal - invoice.DiscountAmount
	invoice.Status = "issued"
}

func persistInvoice(db dbx.IORM, invoice *models.Invoice, items []models.InvoiceItem) error {
	return db.Txn(func(tx dbx.IORM) error {
		if err := tx.Create(invoice); err != nil {
			return err
		}
		for i := range items {
			items[i].InvoiceID = invoice.ID
			if err := tx.Create(&items[i]); err != nil {
				return err
			}
		}
		if err := tx.Save(invoice); err != nil {
			return err
		}
		log.Printf("Created invoice %d with %d items, total=%d", invoice.ID, len(items), invoice.Total)
		return nil
	})
}
