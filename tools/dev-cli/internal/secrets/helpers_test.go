package secrets

import "testing"

func TestParsePairs(t *testing.T) {
	kv, err := parsePairs([]string{"A=1", "B=two"})
	if err != nil {
		t.Fatalf("parsePairs error: %v", err)
	}
	if kv["A"] != "1" || kv["B"] != "two" {
		t.Fatal("parsePairs produced wrong values")
	}
}

func TestBuildNullPatch(t *testing.T) {
	got := buildNullPatch([]string{"A", "B"})
	if got == "" {
		t.Fatal("expected patch json")
	}
}
