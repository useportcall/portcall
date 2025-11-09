package app

import (
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

type App struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	PublicAPIKey string    `json:"public_api_key"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (a *App) Set(app *models.App) *App {
	a.ID = app.PublicID
	a.Name = app.Name
	a.CreatedAt = app.CreatedAt
	a.UpdatedAt = app.UpdatedAt
	a.PublicAPIKey = app.PublicApiKey
	return a
}

type CreateAppRequest struct {
	Name string `json:"name"`
}
