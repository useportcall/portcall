package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/useportcall/portcall/apps/quote/internal/i18n"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func buildI18nLabels(inst *i18n.I18n, lang string) map[string]string {
	return map[string]string{
		"ValidUntil":       inst.T(lang, "valid_until"),
		"PreparedFor":      inst.T(lang, "prepared_for"),
		"PreparedBy":       inst.T(lang, "prepared_by"),
		"QuoteDetails":     inst.T(lang, "quote_details"),
		"PlanName":         inst.T(lang, "plan_name"),
		"BasePrice":        inst.T(lang, "base_price"),
		"FeaturesIncluded": inst.T(lang, "features_included"),
		"QuoteSummary":     inst.T(lang, "quote_summary"),
		"Product":          inst.T(lang, "product"),
		"Price":            inst.T(lang, "price"),
		"Quantity":         inst.T(lang, "quantity"),
		"Frequency":        inst.T(lang, "frequency"),
		"Discount":         inst.T(lang, "discount"),
		"Total":            inst.T(lang, "total"),
		"Subtotal":         inst.T(lang, "subtotal"),
		"SignHere":         inst.T(lang, "sign_here"),
		"ClearSignature":   inst.T(lang, "clear_signature"),
		"LegalConfirm":     inst.T(lang, "legal_confirm"),
		"AcceptSubmit":     inst.T(lang, "accept_submit"),
		"DeclineQuote":     inst.T(lang, "decline_quote"),
		"PaymentTerms":     inst.T(lang, "payment_terms"),
		"JurisdictionNote": inst.T(lang, "jurisdiction_note"),
	}
}

func assembleQuoteData(
	quote *models.Quote, plan *models.Plan, company *models.Company,
	companyAddr, recipientAddr *models.Address,
	items []QuoteItemData, features []FeatureData,
	total int64, basePrice, accessToken string,
	inst *i18n.I18n, lang string,
) QuoteData {
	validUntil := "n/a"
	if quote.ExpiresAt != nil {
		validUntil = quote.ExpiresAt.Format("2006-01-02")
	}

	discount := (total * int64(plan.DiscountPct)) / 100
	preparedByEmail := quote.PreparedByEmail
	if preparedByEmail == "" {
		preparedByEmail = company.Email
	}

	statusKey := "status." + quote.Status
	statusMessage := inst.T(lang, statusKey)
	if statusMessage == statusKey {
		statusMessage = quoteStatusBanner(quote.Status)
	}

	canAccept := quoteCanBeAccepted(quote, time.Now().UTC())
	if quote.Status == "sent" && !canAccept {
		statusMessage = inst.T(lang, "status.expired")
	}

	altLang, altLabel := "ja", "日本語"
	if lang == "ja" {
		altLang, altLabel = "en", "English"
	}

	return QuoteData{
		ID: quote.PublicID, AccessToken: accessToken,
		Status: quote.Status, StatusMessage: statusMessage, CanAccept: canAccept,
		Service: plan.Name,
		Amount:  fmt.Sprintf("$%.2f", float64(total)/100.0),
		Items:   items, Total: fmt.Sprintf("$%.2f", float64(total)/100.0),
		ValidUntil:               validUntil,
		RecipientName:            quote.PublicName,
		RecipientTitle:           quote.PublicTitle,
		RecipientCompanyName:     quote.CompanyName,
		RecipientCompanyNameCaps: strings.ToUpper(quote.CompanyName),
		RecipientAddressLine1:    recipientAddr.Line1,
		RecipientAddressLine2:    recipientAddr.Line2,
		RecipientCity:            recipientAddr.City,
		RecipientPostalCode:      recipientAddr.PostalCode,
		RecipientCountry:         recipientAddr.Country,
		QuoteIssuedByName:        fmt.Sprintf("%s %s", company.FirstName, company.LastName),
		QuoteIssuedByEmail:       preparedByEmail,
		CompanyName:              company.Name,
		CompanyAddressLine1:      companyAddr.Line1,
		CompanyAddressLine2:      companyAddr.Line2,
		CompanyCity:              companyAddr.City,
		CompanyPostalCode:        companyAddr.PostalCode,
		CompanyCountry:           companyAddr.Country,
		CompanyIconURL:           company.IconLogoURL,
		PlanName:                 plan.Name, BasePrice: basePrice,
		Features:         features,
		PlanPrice:        convertCentsToDollars(total),
		BillingFrequency: capitalizeFirstLetter(plan.Interval),
		Discount:         fmt.Sprintf("%d", plan.DiscountPct),
		SubTotal:         convertCentsToDollars(total - discount),
		TaxAmount:        convertCentsToDollars(0),
		TotalAmount:      convertCentsToDollars(total - discount),
		Toc:              quote.Toc,
		JurisdictionNote: inst.T(lang, "jurisdiction_note"),
		Labels:           buildI18nLabels(inst, lang),
		CurrentLang:      lang, AlternateLang: altLang, AlternateLabel: altLabel,
	}
}
