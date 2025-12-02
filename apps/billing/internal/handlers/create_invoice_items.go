package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/useportcall/portcall/apps/billing/internal/utils"
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

	subscriptionID := invoice.SubscriptionID
	if subscriptionID == nil {
		return fmt.Errorf("invoice with ID %d has no associated subscription", p.InvoiceID)
	}

	subscriptionItems := []models.SubscriptionItem{}
	if err := c.DB().List(&subscriptionItems, "subscription_id = ?", *subscriptionID); err != nil {
		return err
	}

	log.Printf("Found %d subscription items for Subscription ID: %d\n", len(subscriptionItems), *subscriptionID)
	if len(subscriptionItems) == 0 {
		return fmt.Errorf("no subscription items found for subscription ID %d", *subscriptionID)
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
			Total:              utils.CalculateTotal(si.PricingModel, si.UnitAmount, si.Quantity, si.Usage, si.Tiers),
			Amount:             getItemUnitAmount(si),
		}
		if err := c.DB().Create(invoiceItem); err != nil {
			return err
		}

		log.Printf("Created Invoice Item: ID=%d, Total=%d\n", invoiceItem.ID, invoiceItem.Total)
	}

	payload := map[string]any{
		"invoice_id": p.InvoiceID,
	}
	if err := c.Queue().Enqueue("calculate_invoice_totals", payload, "billing_queue"); err != nil {
		return err
	}

	return nil
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
