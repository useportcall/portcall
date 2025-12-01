package apix

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type User struct {
	ID                 string    `json:"id" binding:"required"`
	Name               string    `json:"name" binding:"required,min=3"`
	Email              string    `json:"email,omitempty"`
	CompanyTitle       string    `json:"company_title,omitempty"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Subscribed         bool      `json:"subscribed"`
	PaymentMethodAdded bool      `json:"payment_method_added"`
	BillingAddress     *Address  `json:"billing_address"`
}

func (u *User) Set(user *models.User) *User {
	u.ID = user.PublicID
	u.Name = user.Name
	u.Email = user.Email
	u.CompanyTitle = user.CompanyTitle
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
	u.Subscribed = false

	if user.BillingAddress != nil {
		u.BillingAddress = new(Address).Set(user.BillingAddress)
	}

	return u
}
