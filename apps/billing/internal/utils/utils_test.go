package utils

import (
	"testing"

	"github.com/useportcall/portcall/libs/go/dbx/models"
)

func TestCalculateTotal(t *testing.T) {
	// Common tier sets used across tests.
	tieredTiers := []models.Tier{
		{Start: 0, End: 10, Amount: 100}, // usage 1–10
		{Start: 10, End: 20, Amount: 80}, // usage 11–20
		{Start: 20, End: -1, Amount: 50}, // usage 21+ (open-ended)
	}
	blockTiers := []models.Tier{
		{Start: 0, End: 100, Amount: 200}, // flat block price for usage 1–100
		{Start: 100, End: -1, Amount: 150},
	}

	type args struct {
		pricingModel string
		unitAmount   int64
		quantity     int32
		usage        uint
		tiers        *[]models.Tier
	}

	tests := []struct {
		name string
		args args
		want int64
	}{
		// fixed pricing
		{
			name: "fixed - basic multiplication",
			args: args{
				pricingModel: "fixed",
				unitAmount:   100,
				quantity:     2,
				usage:        0,   // ignored
				tiers:        nil, // ignored
			},
			want: 200,
		},
		{
			name: "fixed - zero quantity",
			args: args{
				pricingModel: "fixed",
				unitAmount:   100,
				quantity:     0,
				usage:        10,
				tiers:        nil,
			},
			want: 0,
		},

		// unit pricing
		{
			name: "unit - basic multiplication",
			args: args{
				pricingModel: "unit",
				unitAmount:   10,
				quantity:     3,
				usage:        5,
				tiers:        nil,
			},
			want: 10 * 3 * 5,
		},
		{
			name: "unit - zero usage",
			args: args{
				pricingModel: "unit",
				unitAmount:   10,
				quantity:     3,
				usage:        0,
				tiers:        nil,
			},
			want: 0,
		},

		// tiered pricing
		{
			name: "tiered - nil tiers returns 0",
			args: args{
				pricingModel: "tiered",
				unitAmount:   0,
				quantity:     1,
				usage:        10,
				tiers:        nil,
			},
			want: 0,
		},
		{
			name: "tiered - lower tier (usage in first tier)",
			args: args{
				pricingModel: "tiered",
				unitAmount:   0, // ignored for tiered
				quantity:     2,
				usage:        5, // 0 < 5 <= 10 => first tier
				tiers:        &tieredTiers,
			},
			want: 100 * 2 * 5,
		},
		{
			name: "tiered - upper closed bound (usage == end of first tier)",
			args: args{
				pricingModel: "tiered",
				unitAmount:   0,
				quantity:     1,
				usage:        10, // 0 < 10 <= 10 => first tier
				tiers:        &tieredTiers,
			},
			want: 100 * 1 * 10,
		},
		{
			name: "tiered - middle tier (usage between start and end)",
			args: args{
				pricingModel: "tiered",
				unitAmount:   0,
				quantity:     1,
				usage:        15, // 10 < 15 <= 20 => second tier
				tiers:        &tieredTiers,
			},
			want: 80 * 1 * 15,
		},
		{
			name: "tiered - open ended tier (end == -1)",
			args: args{
				pricingModel: "tiered",
				unitAmount:   0,
				quantity:     3,
				usage:        25, // falls in last tier (end == -1)
				tiers:        &tieredTiers,
			},
			want: 50 * 3 * 25,
		},
		{
			name: "tiered - no matching tier returns 0",
			args: args{
				pricingModel: "tiered",
				unitAmount:   0,
				quantity:     1,
				usage:        0, // 0 is not > any start
				tiers:        &tieredTiers,
			},
			want: 0,
		},

		// block pricing
		{
			name: "block - nil tiers returns 0",
			args: args{
				pricingModel: "block",
				unitAmount:   0,
				quantity:     1,
				usage:        10,
				tiers:        nil,
			},
			want: 0,
		},
		{
			name: "block - usage in first block",
			args: args{
				pricingModel: "block",
				unitAmount:   0,
				quantity:     2,
				usage:        50, // 0 < 50 <= 100 => first tier
				tiers:        &blockTiers,
			},
			want: 200 * 2,
		},
		{
			name: "block - usage in second block (open ended)",
			args: args{
				pricingModel: "block",
				unitAmount:   0,
				quantity:     1,
				usage:        150, // matches second tier (end == -1)
				tiers:        &blockTiers,
			},
			want: 150 * 1,
		},
		{
			name: "block - usage exactly at end of first block",
			args: args{
				pricingModel: "block",
				unitAmount:   0,
				quantity:     1,
				usage:        100, // 0 < 100 <= 100 => first tier
				tiers:        &blockTiers,
			},
			want: 200,
		},

		// default / unknown pricing model
		{
			name: "unknown pricing model returns 0",
			args: args{
				pricingModel: "unknown",
				unitAmount:   100,
				quantity:     10,
				usage:        5,
				tiers:        &tieredTiers,
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateTotal(
				tt.args.pricingModel,
				tt.args.unitAmount,
				tt.args.quantity,
				tt.args.usage,
				tt.args.tiers,
			)
			if got != tt.want {
				t.Fatalf("CalculateTotal(%q, %d, %d, %d, tiers) = %d, want %d",
					tt.args.pricingModel,
					tt.args.unitAmount,
					tt.args.quantity,
					tt.args.usage,
					got,
					tt.want,
				)
			}
		})
	}
}
