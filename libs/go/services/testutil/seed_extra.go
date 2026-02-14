//go:build integration

package testutil

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx"
	"github.com/useportcall/portcall/libs/go/dbx/models"
)

// SeedPlan creates a Plan linked to the given App with the supplied status.
func SeedPlan(
	t *testing.T, db dbx.IORM, appID uint, status string,
) models.Plan {
	t.Helper()
	p := models.Plan{
		PublicID: dbx.GenPublicID("plan"),
		AppID:    appID,
		Name:     "Test Plan",
		Status:   status,
	}
	if err := db.Create(&p); err != nil {
		t.Fatalf("seed plan: %v", err)
	}
	return p
}

// SeedAddress creates an Address for the given App.
func SeedAddress(t *testing.T, db dbx.IORM, appID uint) models.Address {
	t.Helper()
	a := models.Address{
		PublicID:   dbx.GenPublicID("addr"),
		AppID:      appID,
		Line1:      "123 Main St",
		City:       "Testville",
		PostalCode: "12345",
		Country:    "US",
	}
	if err := db.Create(&a); err != nil {
		t.Fatalf("seed address: %v", err)
	}
	return a
}

// SeedCompany creates a Company for the given App.
func SeedCompany(
	t *testing.T, db dbx.IORM, appID, billingAddrID uint,
) models.Company {
	t.Helper()
	c := models.Company{
		AppID:            appID,
		Name:             "Test Co",
		BillingAddressID: billingAddrID,
	}
	if err := db.Create(&c); err != nil {
		t.Fatalf("seed company: %v", err)
	}
	return c
}
