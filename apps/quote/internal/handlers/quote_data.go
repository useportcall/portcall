package handlers

import (
	"github.com/useportcall/portcall/apps/quote/internal/i18n"
	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/routerx"
)

func buildQuoteItems(items []models.PlanItem, lang string, inst *i18n.I18n, planInterval string) ([]QuoteItemData, int64, string) {
	var total int64
	var basePrice string
	var out []QuoteItemData

	for _, item := range items {
		var tiers string
		var itemTotal int64

		switch item.PricingModel {
		case "fixed":
			itemTotal = int64(item.Quantity) * item.UnitAmount
			basePrice = toBasePrice(item.UnitAmount, planInterval)
		case "tiered", "block":
			tiers = inst.T(lang, "pricing.usage_based")
		default:
			tiers = inst.T(lang, "pricing.see_plan")
		}
		if item.PricingModel == "fixed" {
			total += itemTotal
		}

		out = append(out, QuoteItemData{
			Title:       item.PublicTitle,
			Description: item.PublicDescription,
			UnitLabel:   item.PublicUnitLabel,
			PricingType: item.PricingModel,
			Quantity:    item.Quantity,
			UnitAmount:  item.UnitAmount,
			TotalAmount: convertCentsToDollars(itemTotal),
			Tiers:       tiers,
		})
	}
	return out, total, basePrice
}

func loadFeatureNames(c *routerx.Context, planID uint) ([]FeatureData, error) {
	var planFeatures []models.PlanFeature
	if err := c.DB().List(&planFeatures, "plan_id = ?", planID); err != nil {
		return nil, err
	}
	if len(planFeatures) == 0 {
		return nil, nil
	}

	featureIDs := make([]uint, len(planFeatures))
	for i, pf := range planFeatures {
		featureIDs[i] = pf.FeatureID
	}
	var features []models.Feature
	if err := c.DB().List(&features, "id IN ?", featureIDs); err != nil {
		return nil, err
	}

	out := make([]FeatureData, len(features))
	for i, f := range features {
		out[i] = FeatureData{Name: convertSnakeToCamelCaps(f.PublicID)}
	}
	return out, nil
}
