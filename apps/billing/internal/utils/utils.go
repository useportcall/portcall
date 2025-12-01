package utils

func IsPricingMetered(pricingModel string) bool {
	switch pricingModel {
	case "fixed":
		return false
	default:
		return true
	}
}
