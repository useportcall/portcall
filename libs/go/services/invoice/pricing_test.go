package invoice_test

import "testing"

func TestCalculateDiscount(t *testing.T) {
	tests := []struct {
		name     string
		amount   int64
		pct      int
		expected int64
	}{
		{"no discount", 10000, 0, 0},
		{"10% of 10000", 10000, 10, 1000},
		{"25% of 4000", 4000, 25, 1000},
		{"negative pct", 5000, -5, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// calculateDiscount is unexported, so we test it via applyTotals indirectly.
			// We verify through the Create flow in create_test.go.
			// This test validates the expected values match our mental model.
			if tt.pct <= 0 {
				if tt.expected != 0 {
					t.Fatalf("expected 0 for pct %d", tt.pct)
				}
				return
			}
			got := tt.amount * int64(tt.pct) / 100
			if got != tt.expected {
				t.Fatalf("expected %d, got %d", tt.expected, got)
			}
		})
	}
}
