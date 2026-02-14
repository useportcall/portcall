//go:build integration

package testutil

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// SeedAccount creates a minimal Account row.
func SeedAccount(t *testing.T, db dbx.IORM) models.Account {
	t.Helper()
	a := models.Account{
		Email: "test-" + dbx.GenPublicID("") + "@example.com",
	}
	if err := db.Create(&a); err != nil {
		t.Fatalf("seed account: %v", err)
	}
	return a
}

// SeedApp creates an App linked to the given Account.
func SeedApp(t *testing.T, db dbx.IORM, accountID uint) models.App {
	t.Helper()
	app := models.App{
		PublicID:  dbx.GenPublicID("app"),
		Name:      "test-app",
		AccountID: accountID,
	}
	if err := db.Create(&app); err != nil {
		t.Fatalf("seed app: %v", err)
	}
	return app
}

// SeedUser creates a User linked to the given App.
func SeedUser(t *testing.T, db dbx.IORM, appID uint) models.User {
	t.Helper()
	u := models.User{
		PublicID: dbx.GenPublicID("usr"),
		AppID:    appID,
		Name:     "Test User",
		Email:    "user@example.com",
	}
	if err := db.Create(&u); err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return u
}
