package apix

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type Account struct {
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *Account) Set(account *models.Account) *Account {
	a.Email = account.Email
	a.FirstName = account.FirstName
	a.LastName = account.LastName
	a.CreatedAt = account.CreatedAt
	a.UpdatedAt = account.UpdatedAt
	return a
}
