package utils

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func TestCalculateBillingTotal_Fixed(t *testing.T) {
	tests := []struct {
		name       string
		unitAmount int64
		quantity   int32
		expected   int64
	}{
		{"basic fixed", 1000, 1, 1000},
		{"fixed with quantity", 1000, 2, 2000},
		{"zero amount", 0, 1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateBillingTotal("fixed", tt.unitAmount, tt.quantity, 0, nil)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestCalculateBillingTotal_Unit(t *testing.T) {
	tests := []struct {
		name       string
		unitAmount int64
		quantity   int32
		usage      int64
		expected   int64
	}{
		{"basic unit", 10, 1, 100, 1000},         // 10 cents * 100 units
		{"unit with quantity", 10, 2, 100, 2000}, // 10 cents * 2 * 100 units
		{"zero usage", 10, 1, 0, 0},
		{"high volume", 1, 1, 100000, 100000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateBillingTotal("unit", tt.unitAmount, tt.quantity, tt.usage, nil)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestCalculateBillingTotal_Tiered(t *testing.T) {
	// Example tiers: 0-1000 free, 1001-5000 at $0.10, 5001+ at $0.05
	tiers := []models.Tier{
		{Start: 0, End: 1000, Amount: 0},     // First 1000 free
		{Start: 1000, End: 5000, Amount: 10}, // Next 4000 at 10 cents
		{Start: 5000, End: -1, Amount: 5},    // Beyond 5000 at 5 cents
	}

	tests := []struct {
		name     string
		usage    int64
		expected int64
	}{
		{"within free tier", 500, 0},
		{"exactly at free tier limit", 1000, 0},
		{"into second tier", 2000, 10000},       // (2000-1000) * 10 = 10000
		{"full second tier", 5000, 40000},       // (5000-1000) * 10 = 40000
		{"into third tier", 6000, 40000 + 5000}, // 40000 + (6000-5000)*5 = 45000
		{"high usage", 10000, 40000 + 25000},    // 40000 + (10000-5000)*5 = 65000
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateBillingTotal("tiered", 0, 1, tt.usage, &tiers)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestCalculateBillingTotal_Block(t *testing.T) {
	// Example block pricing: 0-1000 = $10, 1001-5000 = $40, 5001+ = $100
	tiers := []models.Tier{
		{Start: 0, End: 1000, Amount: 1000},    // $10 flat
		{Start: 1000, End: 5000, Amount: 4000}, // $40 flat
		{Start: 5000, End: -1, Amount: 10000},  // $100 flat
	}

	tests := []struct {
		name     string
		usage    int64
		expected int64
	}{
		{"first block", 500, 1000},
		{"at first block limit", 1000, 1000},
		{"second block", 2500, 4000},
		{"at second block limit", 5000, 4000},
		{"third block", 6000, 10000},
		{"high usage", 100000, 10000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateBillingTotal("block", 0, 1, tt.usage, &tiers)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestCalculateBillingTotal_Volume(t *testing.T) {
	// Volume pricing: 0-1000 at $0.10, 1001-5000 at $0.08, 5001+ at $0.05
	tiers := []models.Tier{
		{Start: 0, End: 1000, Amount: 10},   // 10 cents per unit for all
		{Start: 1000, End: 5000, Amount: 8}, // 8 cents per unit for all
		{Start: 5000, End: -1, Amount: 5},   // 5 cents per unit for all
	}

	tests := []struct {
		name     string
		usage    int64
		expected int64
	}{
		{"first tier", 500, 5000},    // 500 * 10 cents
		{"second tier", 2500, 20000}, // 2500 * 8 cents (volume applies to all)
		{"third tier", 10000, 50000}, // 10000 * 5 cents
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateBillingTotal("volume", 0, 1, tt.usage, &tiers)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestCalculateBillableUsage(t *testing.T) {
	tests := []struct {
		name      string
		usage     int64
		freeQuota int64
		expected  int64
	}{
		{"no free quota", 100, 0, 100},
		{"within free quota", 50, 100, 0},
		{"at free quota", 100, 100, 0},
		{"over free quota", 150, 100, 50},
		{"negative free quota treated as none", 100, -10, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateBillableUsage(tt.usage, tt.freeQuota)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
