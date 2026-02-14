package models

import (
	"time"

	"gorm.io/gorm"
)

type PaymentLink struct {
	gorm.Model
	PublicID              string    `gorm:"not null;uniqueIndex:idx_public_app"`
	AppID                 uint      `gorm:"not null;uniqueIndex:idx_public_app"`
	App                   App       `gorm:"foreignKey:AppID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PlanID                uint      `gorm:"not null;index:idx_payment_link_plan"`
	Plan                  Plan      `gorm:"foreignKey:PlanID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID                uint      `gorm:"not null;index:idx_payment_link_user"`
	User                  User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ExpiresAt             time.Time `gorm:"not null"`
	Status                string    `gorm:"not null;default:active"` // active, archived
	RedirectURL           *string   `gorm:"default:null"`
	CancelURL             *string   `gorm:"default:null"`
	RequireBillingAddress bool      `gorm:"not null;default:false"`
}
