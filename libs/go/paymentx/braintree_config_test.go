package paymentx

import (
	"os"
	"testing"
)

func TestParseBT_JSON(t *testing.T) {
	raw := `{"merchant_id":"m1","private_key":"pk","environment":"sandbox"}`
	c, err := ParseBraintreeCredentials(raw)
	assertNoErr(t, err)
	assertEqual(t, c.MerchantID, "m1")
	assertEqual(t, c.Environment, "sandbox")
}

func TestParseBT_Colon(t *testing.T) {
	c, err := ParseBraintreeCredentials("mid:pkey:production:acct1")
	assertNoErr(t, err)
	assertEqual(t, c.MerchantID, "mid")
	assertEqual(t, c.PrivateKey, "pkey")
	assertEqual(t, c.MerchantAccount, "acct1")
}

func TestParseBT_Pipe(t *testing.T) {
	c, err := ParseBraintreeCredentials("mid|pkey|sandbox")
	assertNoErr(t, err)
	assertEqual(t, c.MerchantID, "mid")
	assertEqual(t, c.Environment, "sandbox")
}

func TestParseBT_Empty(t *testing.T) {
	_, err := ParseBraintreeCredentials("")
	if err == nil {
		t.Fatal("expected error for empty input")
	}
}

func TestParseBT_Invalid(t *testing.T) {
	_, err := ParseBraintreeCredentials("onlyone")
	if err == nil {
		t.Fatal("expected error for single-part input")
	}
}

func TestParseBT_DefaultEnv(t *testing.T) {
	os.Unsetenv("BRAINTREE_ENVIRONMENT")
	c, err := ParseBraintreeCredentials("mid:pkey")
	assertNoErr(t, err)
	assertEqual(t, c.Environment, "production")
}

func TestParseBT_EnvOverride(t *testing.T) {
	t.Setenv("BRAINTREE_ENVIRONMENT", "sandbox")
	c, err := ParseBraintreeCredentials("mid:pkey")
	assertNoErr(t, err)
	assertEqual(t, c.Environment, "sandbox")
}

func TestParseBT_BadEnv(t *testing.T) {
	_, err := ParseBraintreeCredentials("mid:pkey:badenv")
	if err == nil {
		t.Fatal("expected error for unsupported environment")
	}
}

func TestParseBT_JSONMerchantAcct(t *testing.T) {
	raw := `{"merchant_id":"m","private_key":"k","merchant_account_id":"acct"}`
	c, err := ParseBraintreeCredentials(raw)
	assertNoErr(t, err)
	assertEqual(t, c.MerchantAccount, "acct")
}

func assertNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func assertEqual(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
