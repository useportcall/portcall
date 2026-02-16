package authstore

import "testing"

func TestSaveAndLoadCloudflareToken(t *testing.T) {
	root := t.TempDir()
	if err := SaveCloudflareToken(root, "cf-token-123"); err != nil {
		t.Fatalf("save token: %v", err)
	}
	state, err := Load(root)
	if err != nil {
		t.Fatalf("load token: %v", err)
	}
	if state.CloudflareAPIToken != "cf-token-123" {
		t.Fatalf("unexpected token: %q", state.CloudflareAPIToken)
	}
}
