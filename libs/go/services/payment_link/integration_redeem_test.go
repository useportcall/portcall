//go:build integration

package payment_link_test

import (
	"net/url"
	"strings"
	"testing"

	pl "github.com/useportcall/portcall/libs/go/services/payment_link"
)

func TestRedeem_CreatesCheckoutSession(t *testing.T) {
	appID, user, plan := newEnv(t)
	t.Setenv("CHECKOUT_URL", "https://checkout.test")
	svc := pl.NewService(testDB, noopCrypto{})

	created, err := svc.Create(&pl.CreateInput{AppID: appID, PlanID: plan.PublicID, UserID: user.PublicID})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	parsed, err := url.Parse(created.URL)
	if err != nil {
		t.Fatalf("url.Parse() error = %v", err)
	}
	redeemed, err := svc.Redeem(&pl.RedeemInput{ID: parsed.Query().Get("pl"), Token: parsed.Query().Get("pt")})
	if err != nil {
		t.Fatalf("Redeem() error = %v", err)
	}
	if redeemed.Session == nil || !strings.HasPrefix(redeemed.Session.PublicID, "cs_") {
		t.Fatalf("expected checkout session to be created, got %+v", redeemed.Session)
	}
	if !strings.Contains(redeemed.CheckoutURL, "id=") {
		t.Fatalf("expected checkout URL with session id, got %s", redeemed.CheckoutURL)
	}
}

func TestRedeem_UsesFallbackConnectionWhenDefaultMissing(t *testing.T) {
	appID, user, plan := newEnv(t)
	t.Setenv("CHECKOUT_URL", "https://checkout.test")
	svc := pl.NewService(testDB, noopCrypto{})

	if err := testDB.Exec(
		"UPDATE app_configs SET default_connection_id = NULL WHERE app_id = ?",
		appID,
	); err != nil {
		t.Fatalf("app config update failed: %v", err)
	}

	created, err := svc.Create(&pl.CreateInput{AppID: appID, PlanID: plan.PublicID, UserID: user.PublicID})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	parsed, err := url.Parse(created.URL)
	if err != nil {
		t.Fatalf("url.Parse() error = %v", err)
	}
	redeemed, err := svc.Redeem(&pl.RedeemInput{ID: parsed.Query().Get("pl"), Token: parsed.Query().Get("pt")})
	if err != nil {
		t.Fatalf("Redeem() error = %v", err)
	}
	if redeemed.Session == nil || !strings.HasPrefix(redeemed.Session.PublicID, "cs_") {
		t.Fatalf("expected checkout session to be created, got %+v", redeemed.Session)
	}
}

func TestRedeem_WithoutToken_CreatesCheckoutSession(t *testing.T) {
	appID, user, plan := newEnv(t)
	t.Setenv("CHECKOUT_URL", "https://checkout.test")
	svc := pl.NewService(testDB, noopCrypto{})

	created, err := svc.Create(&pl.CreateInput{AppID: appID, PlanID: plan.PublicID, UserID: user.PublicID})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	parsed, err := url.Parse(created.URL)
	if err != nil {
		t.Fatalf("url.Parse() error = %v", err)
	}
	redeemed, err := svc.Redeem(&pl.RedeemInput{ID: parsed.Query().Get("pl")})
	if err != nil {
		t.Fatalf("Redeem() error = %v", err)
	}
	if redeemed.Session == nil || !strings.HasPrefix(redeemed.Session.PublicID, "cs_") {
		t.Fatalf("expected checkout session to be created, got %+v", redeemed.Session)
	}
}
