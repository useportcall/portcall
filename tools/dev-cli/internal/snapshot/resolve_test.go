package snapshot

import "testing"

func TestLookupTargets_GroupAndTargetDedup(t *testing.T) {
	got, err := lookupTargets([]string{"invoice", "invoice-light"})
	if err != nil {
		t.Fatalf("lookupTargets error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 targets, got %d", len(got))
	}
}

func TestCollectTargets_InvalidMode(t *testing.T) {
	_, err := collectTargets(nil, "bad-mode")
	if err == nil {
		t.Fatal("expected error for invalid mode")
	}
}
