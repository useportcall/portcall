package apix

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type AppConfig struct {
	DefaultConnection string    `json:"default_connection"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (a *AppConfig) Set(data *models.AppConfig) *AppConfig {
	a.DefaultConnection = data.DefaultConnection.PublicID
	a.CreatedAt = data.CreatedAt
	a.UpdatedAt = data.UpdatedAt
	return a
}
