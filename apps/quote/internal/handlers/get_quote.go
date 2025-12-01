package handlers

import (
	"fmt"
	"net/http"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

type QuoteItemData struct {
	Title       string
	Description string
	UnitLabel   string
	PricingType string // fixed, tiered, block, etc.
	Quantity    int32
	UnitAmount  int64
	TotalAmount string
	Tiers       string // for display, if applicable
}

type QuoteData struct {
	ID         string
	Service    string
	Amount     string
	ValidUntil string
	Items      []QuoteItemData
	Total      string
}

func GetQuote(c *routerx.Context) {
	var quote models.Quote
	if err := c.DB().FindFirst(&quote, "public_id = ?", c.Param("id")); err != nil {
		c.NotFound("quote not found")
		return
	}

	var plan models.Plan
	if err := c.DB().FindForID(quote.PlanID, &plan); err != nil {
		c.NotFound("quote not found")
		return
	}

	var items []models.PlanItem
	if err := c.DB().List(&items, "plan_id = ?", plan.ID); err != nil {
		c.NotFound("quote not found")
		return
	}

	// Calculate total and build item details
	var total int64 = 0
	var itemDatas []QuoteItemData
	for _, item := range items {
		pricingType := item.PricingModel
		unitLabel := item.PublicUnitLabel
		title := item.PublicTitle
		desc := item.PublicDescription
		quantity := item.Quantity
		unitAmount := item.UnitAmount
		var tiers string
		var itemTotal int64

		switch pricingType {
		case "fixed":
			itemTotal = int64(quantity) * unitAmount
		case "tiered", "block":
			// For demo, just show unit amount and note it's usage-based
			itemTotal = 0
			tiers = "Usage-based pricing applies"
		default:
			itemTotal = 0
			tiers = "See plan for details"
		}
		if pricingType == "fixed" {
			total += itemTotal
		}

		itemDatas = append(itemDatas, QuoteItemData{
			Title:       title,
			Description: desc,
			UnitLabel:   unitLabel,
			PricingType: pricingType,
			Quantity:    quantity,
			UnitAmount:  unitAmount,
			TotalAmount: fmt.Sprintf("$%.2f", float64(itemTotal)/100.0),
			Tiers:       tiers,
		})
	}

	// Example: get data from query params or use mock data
	data := QuoteData{
		ID:         quote.PublicID,
		Service:    c.DefaultQuery("service", plan.Name),
		Amount:     fmt.Sprintf("$%.2f", float64(total)/100.0),
		ValidUntil: c.DefaultQuery("valid_until", "2025-09-30"),
		Items:      itemDatas,
		Total:      fmt.Sprintf("$%.2f", float64(total)/100.0),
	}
	c.HTML(http.StatusOK, "quote.html", data)
}
