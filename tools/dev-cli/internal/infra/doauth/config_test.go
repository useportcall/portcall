package doauth

import "testing"

func TestFindContextToken(t *testing.T) {
	cfg := `access-token: ""
auth-contexts:
  default: "true"
  main: dop_v1_main123
  portcall: dop_v1_portcall456
backup-policies:
  get:
    format: ""`
	if got := FindContextToken(cfg, "portcall"); got != "dop_v1_portcall456" {
		t.Fatalf("unexpected token: %q", got)
	}
	if got := FindContextToken(cfg, "missing"); got != "" {
		t.Fatalf("expected empty token for missing context, got %q", got)
	}
}
