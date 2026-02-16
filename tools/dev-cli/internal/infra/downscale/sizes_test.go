package downscale

import "testing"

func TestIsDowngradeSize(t *testing.T) {
	if !IsDowngradeSize("db-s-1vcpu-2gb", "db-s-1vcpu-1gb") {
		t.Fatalf("expected downgrade detection")
	}
	if IsDowngradeSize("db-s-1vcpu-2gb", "db-s-2vcpu-4gb") {
		t.Fatalf("did not expect upgrade to be flagged as downgrade")
	}
	if IsDowngradeSize("unknown", "db-s-1vcpu-1gb") {
		t.Fatalf("unknown size should not trigger downgrade")
	}
}
