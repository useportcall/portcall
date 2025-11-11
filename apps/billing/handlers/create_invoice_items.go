package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

type CreateInvoiceItemsPayload struct {
	InvoiceID uint `json:"invoice_id"`
}

func CreateInvoiceItems(c server.IContext) error {
	var p CreateInvoiceItemsPayload
	if err := json.Unmarshal(c.Payload(), &p); err != nil {
		return err
	}

	var invoice models.Invoice
	if err := c.DB().FindForID(p.InvoiceID, &invoice); err != nil {
		return fmt.Errorf("failed to find invoice with ID %d: %w", p.InvoiceID, err)
	}

	var subscriptionItems []models.SubscriptionItem
	if err := c.DB().List(&subscriptionItems, "subscription_id = ?", invoice.SubscriptionID); err != nil {
		return err
	}

	for _, si := range subscriptionItems {
		invoiceItem := &models.InvoiceItem{
			PublicID:           dbx.GenPublicID("ii"),
			AppID:              si.AppID,
			InvoiceID:          p.InvoiceID,
			Quantity:           si.Quantity,
			SubscriptionItemID: si.ID,
			Title:              si.Title,
			Description:        si.Description,
			PricingModel:       si.PricingModel,
			Total:              calculateTotal(si.PricingModel, si.UnitAmount, si.Quantity, si.Usage),
			Amount:             getItemUnitAmount(si),
		}
		if err := c.DB().Create(invoiceItem); err != nil {
			return err
		}
	}

	payload := map[string]any{
		"invoice_id": p.InvoiceID,
	}
	if err := c.Queue().Enqueue("calculate_invoice_totals", payload, "billing_queue"); err != nil {
		return err
	}

	return nil
}

// TODO: look into more robust way to calculate large totals
func calculateTotal(pricingModel string, unitAmount int64, quantity int32, usage uint) int64 {
	switch pricingModel {
	case "fixed":
		return int64(unitAmount) * int64(quantity)
	case "unit":
		return int64(unitAmount) * int64(usage) * int64(quantity)
	case "tiered":
		// for _, tier := range price.TieredPricing.Tiers {
		// 	if tier.End == nil || *tier.End >= usage {
		// 		return int64(tier.Amount) * int64(quantity) * int64(usage)
		// 	}
		// }

		return 0 //TODO: or handle the case where no tier matches
	case "block":
		// TODO: implement block pricing logic
		return 0
	case "volume":
		// TODO: implement volume pricing logic
		return 0
	default:
		return 0 // TODO: or handle unknown pricing model
	}
}

func getItemUnitAmount(si models.SubscriptionItem) int64 {
	switch si.PricingModel {
	case "fixed":
		return si.UnitAmount
	case "unit":
		return si.UnitAmount
	case "tiered":
		if si.Tiers == nil {
			return 0
		}

		if len(*si.Tiers) > 0 {
			return int64((*si.Tiers)[0].Amount)
		}

		return 0 // or handle no tiers case
	case "block":
		// TODO: implement block pricing logic
		return 0
	case "volume":
		// TODO: implement volume pricing logic
		return 0
	default:
		return 0 // or handle unknown pricing model
	}
}
