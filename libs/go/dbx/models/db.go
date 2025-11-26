package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Email     string `gorm:"not null;unique"`
	FirstName string `gorm:"default:null"`
	LastName  string `gorm:"default:null"`
}

type App struct {
	gorm.Model
	PublicID     string  `gorm:"not null;unique"`
	Name         string  `gorm:"not null"`
	AccountID    uint    `gorm:"not null;index:idx_app_account"`
	Account      Account `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	IsLive       bool    `gorm:"not null;default:false"`
	Status       string  `gorm:"not null;default:'draft'"`
	PublicApiKey string  `gorm:"default:null;unique"`
}

type Address struct {
	gorm.Model
	PublicID   string `gorm:"not null;uniqueIndex:idx_public_app"`
	AppID      uint   `gorm:"not null;uniqueIndex:idx_public_app"`
	App        App    `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Line1      string `gorm:"not null"`
	Line2      string `gorm:"default:null"`
	City       string `gorm:"not null"`
	State      string `gorm:"default:null"`
	PostalCode string `gorm:"not null"`
	Country    string `gorm:"not null"`
}

type Company struct {
	gorm.Model
	AppID            uint    `gorm:"not null;unique"`
	App              App     `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name             string  `gorm:"default:null"`
	Alias            string  `gorm:"default:null"`
	FirstName        string  `gorm:"default:null"`
	LastName         string  `gorm:"default:null"`
	Email            string  `gorm:"default:null"`
	Phone            string  `gorm:"default:null"`
	VATNumber        string  `gorm:"default:null"`
	BusinessCategory string  `gorm:"default:null"`
	BillingAddressID uint    `gorm:"default:null"`
	BillingAddress   Address `gorm:"foreignKey:BillingAddressID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type Feature struct {
	gorm.Model
	PublicID  string `gorm:"not null;uniqueIndex:idx_feature_public_app"`
	AppID     uint   `gorm:"not null;uniqueIndex:idx_feature_public_app"`
	App       App    `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	IsMetered bool   `gorm:"not null;default:false"` // true if the feature is metered
}

type User struct {
	gorm.Model
	PublicID          string   `gorm:"not null;uniqueIndex:idx_public_app" json:"id"`
	AppID             uint     `gorm:"not null;uniqueIndex:idx_public_app" json:"-"`
	App               App      `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Name              string   `gorm:"not null"`
	Email             string   `gorm:"not null"`
	PaymentCustomerID string   `gorm:"default:null"` // Stripe customer ID for the user
	BillingAddressID  *uint    `gorm:"default:null"` // nullable, foreign key to the user's billing address
	BillingAddress    *Address `gorm:"foreignKey:BillingAddressID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type Entitlement struct {
	gorm.Model
	AppID           uint    `gorm:"not null;uniqueIndex:idx_app_user_feature"`
	App             App     `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID          uint    `gorm:"not null;uniqueIndex:idx_app_user_feature"`
	User            User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	FeaturePublicID string  `gorm:"not null;uniqueIndex:idx_app_user_feature"`
	Feature         Feature `gorm:"foreignKey:FeaturePublicID,AppID;references:PublicID,AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	//
	Interval  string         `gorm:"column:interval;not null;default:month"`
	Usage     int64          `gorm:"column:usage;default:0"`
	Quota     int64          `gorm:"column:quota;default:-1"`
	Rollover  int            `gorm:"not null;default:0"` // amount of quota that can be rolled over to the next period
	Topup     int64          `gorm:"column:topup;default:0"`
	Keys      datatypes.JSON `gorm:"column:keys;type:jsonb;not null;default:'[]'"`
	Tag       *string        `gorm:"column:tag"`
	IsMetered bool           `gorm:"not null;default:false"` // true if the entitlement is metered

	LastResetAt *time.Time `gorm:"column:last_reset_at;default:CURRENT_TIMESTAMP"`
	NextResetAt *time.Time `gorm:"column:next_reset_at;default:CURRENT_TIMESTAMP"`
	AnchorAt    *time.Time `gorm:"column:anchor_at"`
}

type Secret struct {
	gorm.Model
	PublicID   string     `gorm:"not null;unique"` // api key
	AppID      uint       `gorm:"not null"`
	App        App        `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	KeyHash    string     `gorm:"column:key_hash;not null;unique"`
	DisabledAt *time.Time `gorm:"column:disabled_at"`
	KeyType    string     `gorm:"column:key_type;not null"`
}

type Connection struct {
	gorm.Model
	PublicID               string  `gorm:"not null;uniqueIndex:idx_public_app" json:"id"`
	AppID                  uint    `gorm:"not null;uniqueIndex:idx_public_app" json:"-"`
	App                    App     `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Name                   string  `gorm:"not null;default:'New Connection'"`
	Source                 string  `gorm:"not null"`
	PublicKey              string  `gorm:"not null"`
	EncryptedKey           string  `gorm:"not null"`
	EncryptedWebhookSecret *string `gorm:"default:null"` // nullable, encrypted webhook secret for Stripe
}

type Quote struct {
	gorm.Model
	PublicID       string     `gorm:"not null;uniqueIndex:idx_public_app"`
	AppID          uint       `gorm:"not null;uniqueIndex:idx_public_app"`
	App            App        `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PlanID         uint       `gorm:"default:null"`
	Plan           Plan       `gorm:"default:null;foreignKey:PlanID"`
	UserID         *uint      `gorm:"not null"`
	User           *User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PublicTitle    string     `gorm:"not null"`
	PublicName     string     `gorm:"not null"`
	Status         string     `gorm:"not null;default:draft"` // created, sent, accepted, declined, expired
	DaysValid      int        `gorm:"not null;default:30"`    // number of days the quote is valid for
	IssuedAt       *time.Time `gorm:"default:null"`           // when the quote was issued
	ExpiresAt      *time.Time `gorm:"default:null"`           // when the quote expires
	URL            *string    `gorm:"default:null"`
	SignatureURL   *string    `gorm:"default:null"`
	TokenHash      *string    `gorm:"default:null"`
	DirectCheckout bool       `gorm:"not null;default:true"`
}

type PlanGroup struct {
	gorm.Model
	PublicID string `gorm:"not null;uniqueIndex:idx_public_app"`
	AppID    uint   `gorm:"not null;uniqueIndex:idx_public_app"`
	App      App    `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name     string `gorm:"not null" json:"name"`
}

type Plan struct {
	gorm.Model
	PublicID         string     `gorm:"not null;uniqueIndex:idx_public_app"`
	AppID            uint       `gorm:"not null;uniqueIndex:idx_public_app"`
	App              App        `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PlanGroupID      *uint      `gorm:"default:null"`
	PlanGroup        *PlanGroup `gorm:"foreignKey:PlanGroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name             string     `gorm:"not null"`
	Status           string     `gorm:"not null;default:draft"` // active, inactive, archived
	TrialPeriodDays  int        `gorm:"not null;default:0"`     // number of days for the trial period
	Interval         string     `gorm:"not null;default:month"` // billing interval, e.g., month, year
	IntervalCount    int        `gorm:"not null;default:1"`     // number of intervals for the billing cycle
	Currency         string     `gorm:"not null;default:USD"`   // currency code, e.g., USD, EUR
	InvoiceDueByDays int        `gorm:"not null;default:10"`
	IsFree           bool       `gorm:"not null;default:false"`
}

type PlanItem struct {
	gorm.Model
	PublicID          string        `gorm:"not null;uniqueIndex:idx_public_app"`
	AppID             uint          `gorm:"not null;uniqueIndex:idx_public_app"`
	App               App           `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PlanID            uint          `gorm:"not null"`
	Plan              Plan          `gorm:"foreignKey:PlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PricingModel      string        `gorm:"not null"`
	Quantity          int32         `gorm:"not null;default:1"` // quantity of the plan item
	UnitAmount        int64         `gorm:"not null"`           // unit amount in cents
	Maximum           *int          `gorm:"default:null"`
	Minimum           *int          `gorm:"default:null"`
	Tiers             *[]Tier       `gorm:"type:jsonb;serializer:json"` // tiers for tiered pricing
	PublicTitle       string        `gorm:"not null"`
	PublicDescription string        `gorm:"not null"`
	PublicUnitLabel   string        `gorm:"not null"`
	PlanFeatures      []PlanFeature `gorm:"constraint:OnDelete:CASCADE;"`
}

type PlanFeature struct {
	gorm.Model
	PublicID   string   `gorm:"not null;uniqueIndex:idx_public_app"`
	AppID      uint     `gorm:"not null;uniqueIndex:idx_public_app"`
	App        App      `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PlanID     uint     `gorm:"not null;"`
	Plan       Plan     `gorm:"foreignKey:PlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PlanItemID uint     `gorm:"default:null;uniqueIndex:idx_plan_item_feature"`
	PlanItem   PlanItem `gorm:"default:null;foreignKey:PlanItemID;constraint:OnDelete:CASCADE"`
	FeatureID  uint     `gorm:"default:null;uniqueIndex:idx_plan_item_feature"`
	Feature    Feature  `gorm:"default:null;foreignKey:FeatureID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Interval   string   `gorm:"not null;default:month"` // month, year, etc.
	Quota      int      `gorm:"not null;default:-1"`    // -1 means unlimited, 0 means no quota
	Rollover   int      `gorm:"not null;default:0"`     // amount of quota that can be rolled over to the next period
}

type Subscription struct {
	gorm.Model
	PublicID             string        `gorm:"not null;uniqueIndex:idx_public_app" json:"id"`
	AppID                uint          `gorm:"not null;uniqueIndex:idx_public_app" json:"-"`
	App                  App           `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	UserID               uint          `gorm:"not null"`
	User                 User          `gorm:"foreignKey:UserID"`
	Status               string        `gorm:"not null"`
	Currency             string        `gorm:"not null;default:USD"`
	InvoiceDueByDays     int           `gorm:"not null;default:10"`
	LastResetAt          time.Time     `gorm:"not null"`
	NextResetAt          time.Time     `gorm:"not null"`
	FinalResetAt         *time.Time    `gorm:"default:null"`
	InvoiceConfig        InvoiceConfig `gorm:"type:jsonb;not null"`
	BillingInterval      string        `gorm:"not null;default:month"` // billing interval, e.g., month, year
	BillingIntervalCount int           `gorm:"not null;default:1"`     // number of intervals for the billing cycle
	BillingAddressID     uint          `gorm:"not null"`
	BillingAddress       Address       `gorm:"foreignKey:BillingAddressID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	PlanID               *uint         `gorm:"default:null"`
	Plan                 *Plan         `gorm:"default:null;foreignKey:PlanID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	RollbackPlanID       *uint         `gorm:"default:null"`
	RollbackPlan         *Plan         `gorm:"default:null;foreignKey:RollbackPlanID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type SubscriptionItem struct {
	gorm.Model
	PublicID       string       `gorm:"not null;uniqueIndex:idx_public_app" json:"id"`
	AppID          uint         `gorm:"not null;uniqueIndex:idx_public_app" json:"-"`
	App            App          `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	SubscriptionID uint         `gorm:"not null"`
	Subscription   Subscription `gorm:"foreignKey:SubscriptionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	FeatureID      *uint        `gorm:"default:null"`
	Feature        *Feature     `gorm:"foreignKey:FeatureID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Usage          uint         `gorm:"not null;default:0"`
	PlanItemID     *uint        `gorm:"default:null"`
	PlanItem       *PlanItem    `gorm:"default:null;foreignKey:PlanItemID"`
	PricingModel   string       `gorm:"not null"` // fixed, tiered, volume, etc.
	UnitAmount     int64        `gorm:"not null"` // unit amount in cents
	Tiers          *[]Tier      `gorm:"type:jsonb;serializer:json"`
	Maximum        *int         `gorm:"default:null"`
	Minimum        *int         `gorm:"default:null"`
	Quantity       int32        `gorm:"not null"` // quantity of the subscription item
	Title          string       `gorm:"not null"` // title of the subscription item
	Description    string       `gorm:"not null"` // description of the subscription item
}

type CheckoutSession struct {
	gorm.Model
	PublicID             string    `gorm:"not null;unique"`
	AppID                uint      `gorm:"not null"`
	App                  App       `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PlanID               uint      `gorm:"default:null;index:idx_checkout_session_plan"`
	Plan                 Plan      `gorm:"foreignKey:PlanID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	UserID               uint      `gorm:"not null"`
	User                 User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ExpiresAt            time.Time `gorm:"not null"`
	ExternalSessionID    string    `gorm:"not null;unique"`
	ExternalClientSecret string    `gorm:"not null"`
	ExternalPublicKey    string    `gorm:"not null"`
	ExternalProvider     string    `gorm:"not null"` // e.g., stripe, local
	RedirectURL          *string   `gorm:"not null"` // URL to redirect after checkout
	CancelURL            *string   `gorm:"not null"` // URL to redirect if the user cancels the checkout
	BillingAddressID     *uint     `gorm:"default:null"`
	BillingAddress       *Address  `gorm:"foreignKey:BillingAddressID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CompanyAddressID     uint      `gorm:"default:null"`
	CompanyAddress       Address   `gorm:"foreignKey:CompanyAddressID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Status               string    `gorm:"not null;default:active"` // active, completed, canceled, pending
}

type Invoice struct {
	gorm.Model
	PublicID          string        `gorm:"not null;uniqueIndex:idx_public_app"`
	AppID             uint          `gorm:"not null;uniqueIndex:idx_public_app"`
	App               App           `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SubscriptionID    *uint         `gorm:"default:null"`
	Subscription      *Subscription `gorm:"foreignKey:SubscriptionID;constraint:OnDelete:CASCADE"`
	UserID            uint          `gorm:"not null"`
	User              User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Currency          string        `gorm:"not null"`
	Total             int64         `gorm:"not null"`
	SubTotal          int64         `gorm:"not null"`
	TaxAmount         int64         `gorm:"not null;default:0"`     // total tax amount in cents
	DiscountAmount    int64         `gorm:"not null;default:0"`     // total discount amount in cents
	DecimalPlaces     int           `gorm:"not null;default:2"`     // number of decimal places for the total
	Status            string        `gorm:"not null;default:draft"` // draft, paid, voided, refunded
	DueBy             time.Time     `gorm:"not null"`
	InvoiceNumber     string        `gorm:"not null"`
	PDFURL            string        `gorm:"default:null"` // nullable, URL to the PDF invoice
	EmailURL          string        `gorm:"default:null"` // nullable, URL to the email invoice
	CompanyName       string        `gorm:"not null"`     // name of the company issuing the invoice
	CompanyAddressID  uint          `gorm:"default:null"` // nullable, foreign key to the company address
	CompanyAddress    Address       `gorm:"default:null;foreignKey:CompanyAddressID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	BillingAddressID  uint          `gorm:"not null"`
	BillingAddress    Address       `gorm:"foreignKey:BillingAddressID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	ShippingAddressID *uint         `gorm:"default:null"` // nullable, foreign key to the shipping address
	ShippingAddress   *Address      `gorm:"foreignKey:ShippingAddressID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CustomerName      string        `gorm:"not null"` // name of the customer for the invoice
	CustomerEmail     string        `gorm:"not null"` // email of the customer for the
}

type InvoiceItem struct {
	gorm.Model
	PublicID           string           `gorm:"not null;uniqueIndex:idx_public_app"`
	AppID              uint             `gorm:"not null;uniqueIndex:idx_public_app"`
	App                App              `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	InvoiceID          uint             `gorm:"not null"`
	Invoice            Invoice          `gorm:"foreignKey:InvoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SubscriptionItemID uint             `gorm:"default:null"`
	SubscriptionItem   SubscriptionItem `gorm:"foreignKey:SubscriptionItemID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Total              int64            `gorm:"not null"`           // total amount in cents
	Amount             int64            `gorm:"not null"`           // amount in cents
	Quantity           int32            `gorm:"not null;default:1"` // quantity of the invoice item
	Title              string           `gorm:"not null"`           // title of the invoice item
	Description        string           `gorm:"not null"`           // description of the invoice item
	PricingModel       string           `gorm:"not null"`           // pricing model for the item, e.g., fixed, tiered, volume
}

type PaymentMethod struct {
	gorm.Model
	PublicID     string `gorm:"not null;uniqueIndex:idx_public_app"`
	AppID        uint   `gorm:"not null;uniqueIndex:idx_public_app"`
	App          App    `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID       uint   `gorm:"not null"`
	User         User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ExternalID   string `gorm:"not null"`
	ExternalType string `gorm:"not null"` // e.g., card, bank_account
}

type Email struct {
	gorm.Model
	to   string `gorm:"not null"`
	from string `gorm:"not null"`
	body string `gorm:"not null"`
}

type MeterEvent struct {
	gorm.Model
	AppID     uint      `gorm:"not null"`
	App       App       `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	FeatureID uint      `gorm:"not null"`
	Feature   Feature   `gorm:"foreignKey:FeatureID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Usage     int64     `gorm:"not null"`
	Timestamp time.Time `gorm:"not null"`
}

type AppConfig struct {
	gorm.Model
	AppID               uint       `gorm:"not null;unique"` // unique
	App                 App        `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	DefaultConnectionID uint       `gorm:"default:null"`
	DefaultConnection   Connection `gorm:"foreignKey:DefaultConnectionID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
