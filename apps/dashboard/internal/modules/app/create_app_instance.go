package app

import (
	"github.com/useportcall/portcall/apps/dashboard/internal/utils"
	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func createAppInstance(txn dbx.IORM, accountID uint, name string, isLive bool) (*models.App, error) {
	app := &models.App{
		PublicID:     utils.GenPublicID("app"),
		AccountID:    accountID,
		Name:         name,
		IsLive:       isLive,
		PublicApiKey: utils.GenPublicID("pk"),
	}
	if err := txn.Create(app); err != nil {
		return nil, err
	}

	address := &models.Address{AppID: app.ID, Line1: "123 Main St", City: "Anytown", State: "CA", PostalCode: "12345", Country: "USA"}
	if err := txn.Create(address); err != nil {
		return nil, err
	}

	company := &models.Company{AppID: app.ID, BillingAddressID: address.ID, Name: "Default Company", Email: "default@example.com", VATNumber: "123456789"}
	if err := txn.Create(company); err != nil {
		return nil, err
	}

	appConfig := &models.AppConfig{AppID: app.ID}
	if !isLive {
		connection := &models.Connection{AppID: app.ID, Source: "local", Name: "Mock Payment Provider (Test Only)"}
		if err := txn.Create(connection); err != nil {
			return nil, err
		}
		appConfig.DefaultConnectionID = connection.ID
	}
	if err := txn.Create(appConfig); err != nil {
		return nil, err
	}
	return app, nil
}
