package handlers

// QuoteItemData holds a single line-item for the quote template.
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

// FeatureData holds a feature name for the quote template.
type FeatureData struct {
	Name string
}

// QuoteData is the top-level view model passed to the quote HTML template.
type QuoteData struct {
	ID                       string
	AccessToken              string
	Status                   string
	StatusMessage            string
	CanAccept                bool
	Service                  string
	Amount                   string
	ValidUntil               string
	Items                    []QuoteItemData
	Features                 []FeatureData
	Total                    string
	RecipientName            string
	RecipientTitle           string
	RecipientCompanyName     string
	RecipientCompanyNameCaps string
	RecipientAddressLine1    string
	RecipientAddressLine2    string
	RecipientCity            string
	RecipientPostalCode      string
	RecipientCountry         string
	QuoteIssuedByName        string
	QuoteIssuedByEmail       string
	CompanyName              string
	CompanyAddressLine1      string
	CompanyAddressLine2      string
	CompanyCity              string
	CompanyPostalCode        string
	CompanyCountry           string
	CompanyIconURL           string
	PlanName                 string
	BasePrice                string
	PlanPrice                string
	BillingFrequency         string
	Discount                 string
	SubTotal                 string
	TaxAmount                string
	TotalAmount              string
	Toc                      string
	JurisdictionNote         string
	// i18n
	Labels         map[string]string
	CurrentLang    string
	AlternateLang  string
	AlternateLabel string
}
