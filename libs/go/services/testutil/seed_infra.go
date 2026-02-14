//go:build integration

package testutil

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// SeedConnection creates a "local" Connection for the given App.
func SeedConnection(t *testing.T, db dbx.IORM, appID uint) models.Connection {
	t.Helper()
	c := models.Connection{
		PublicID:     dbx.GenPublicID("conn"),
		AppID:        appID,
		Name:         "Local",
		Source:       "local",
		PublicKey:    "pk_test",
		EncryptedKey: "unused",
	}
	if err := db.Create(&c); err != nil {
		t.Fatalf("seed connection: %v", err)
	}
	return c
}

// SeedAppConfig creates an AppConfig pointing to the given Connection.
func SeedAppConfig(
	t *testing.T, db dbx.IORM, appID, connID uint,
) models.AppConfig {
	t.Helper()
	cfg := models.AppConfig{AppID: appID, DefaultConnectionID: connID}
	if err := db.Create(&cfg); err != nil {
		t.Fatalf("seed app config: %v", err)
	}
	return cfg
}

// RequireNoErr fails the test immediately on non-nil error.
func RequireNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
