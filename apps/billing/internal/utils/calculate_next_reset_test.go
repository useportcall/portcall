// calculate_next_reset_test.go
package utils

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestCalculateNextReset(t *testing.T) {
	tests := []struct {
		name     string
		interval string
		now      time.Time
		wantErr  error
		wantTime time.Time
	}{
		{
			name:     "week - adds 7 days, preserves time of day",
			interval: "week",
			now:      time.Date(2024, 11, 15, 10, 30, 5, 123456789, time.UTC),
			wantErr:  nil,
			wantTime: time.Date(2024, 11, 22, 10, 30, 5, 123456789, time.UTC),
		},
		{
			name:     "month - simple next month same day",
			interval: "month",
			now:      time.Date(2024, 5, 10, 8, 0, 0, 0, time.UTC),
			wantErr:  nil,
			wantTime: time.Date(2024, 6, 10, 8, 0, 0, 0, time.UTC),
		},
		{
			name:     "month - Jan 31 → Feb 29 (leap year)",
			interval: "month",
			now:      time.Date(2024, 1, 31, 14, 0, 0, 0, time.UTC),
			wantErr:  nil,
			wantTime: time.Date(2024, 2, 29, 14, 0, 0, 0, time.UTC),
		},
		{
			name:     "month - Jan 30 → Feb 29 (leap year)",
			interval: "month",
			now:      time.Date(2024, 1, 30, 14, 0, 0, 0, time.UTC),
			wantErr:  nil,
			wantTime: time.Date(2024, 2, 29, 14, 0, 0, 0, time.UTC),
		},
		{
			name:     "month - March 31 → April 30",
			interval: "month",
			now:      time.Date(2024, 3, 31, 12, 0, 0, 0, time.UTC),
			wantErr:  nil,
			wantTime: time.Date(2024, 4, 30, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "month - December → January next year",
			interval: "month",
			now:      time.Date(2024, 12, 15, 9, 0, 0, 0, time.UTC),
			wantErr:  nil,
			wantTime: time.Date(2025, 1, 15, 9, 0, 0, 0, time.UTC),
		},
		{
			name:     "year - simple next year same day",
			interval: "year",
			now:      time.Date(2024, 6, 15, 9, 15, 0, 0, time.UTC),
			wantErr:  nil,
			wantTime: time.Date(2025, 6, 15, 9, 15, 0, 0, time.UTC),
		},
		{
			name:     "year - Feb 29 → Feb 28 (non-leap year)",
			interval: "year",
			now:      time.Date(2024, 2, 29, 7, 0, 0, 0, time.UTC),
			wantErr:  nil,
			wantTime: time.Date(2025, 2, 28, 7, 0, 0, 0, time.UTC),
		},
		{
			name:     "invalid interval",
			interval: "garbage",
			now:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr:  errors.New("invalid interval"),
			wantTime: time.Time{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err, got := CalculateNextReset(tt.interval, tt.now)

			// Check error presence
			if (err == nil) != (tt.wantErr == nil) {
				t.Fatalf("expected err=%v, got err=%v", tt.wantErr, err)
			}
			// If both non-nil, check that error contains expected message
			if err != nil && !strings.Contains(err.Error(), tt.wantErr.Error()) {
				t.Fatalf("expected error containing %q, got %q", tt.wantErr.Error(), err.Error())
			}

			// Check time result
			if !got.Equal(tt.wantTime) {
				t.Fatalf("expected time %v, got %v", tt.wantTime, got)
			}
		})
	}
}

func TestCalculateNextResetWithCount(t *testing.T) {
	tests := []struct {
		name     string
		interval string
		count    int
		now      time.Time
		wantTime time.Time
		wantErr  bool
	}{
		{
			name:     "week count=1 - adds 7 days",
			interval: "week",
			count:    1,
			now:      time.Date(2024, 11, 15, 10, 30, 5, 0, time.UTC),
			wantTime: time.Date(2024, 11, 22, 10, 30, 5, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "week count=2 - adds 14 days",
			interval: "week",
			count:    2,
			now:      time.Date(2024, 11, 15, 10, 30, 5, 0, time.UTC),
			wantTime: time.Date(2024, 11, 29, 10, 30, 5, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "month count=1 - adds 1 month",
			interval: "month",
			count:    1,
			now:      time.Date(2024, 5, 10, 8, 0, 0, 0, time.UTC),
			wantTime: time.Date(2024, 6, 10, 8, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "month count=3 - quarterly billing",
			interval: "month",
			count:    3,
			now:      time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC),
			wantTime: time.Date(2024, 4, 15, 8, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "month count=6 - semi-annual billing",
			interval: "month",
			count:    6,
			now:      time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC),
			wantTime: time.Date(2024, 7, 15, 8, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "month count=12 - annual via months",
			interval: "month",
			count:    12,
			now:      time.Date(2024, 1, 15, 8, 0, 0, 0, time.UTC),
			wantTime: time.Date(2025, 1, 15, 8, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "year count=1 - adds 1 year",
			interval: "year",
			count:    1,
			now:      time.Date(2024, 6, 15, 9, 15, 0, 0, time.UTC),
			wantTime: time.Date(2025, 6, 15, 9, 15, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "year count=2 - adds 2 years",
			interval: "year",
			count:    2,
			now:      time.Date(2024, 6, 15, 9, 15, 0, 0, time.UTC),
			wantTime: time.Date(2026, 6, 15, 9, 15, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "month count=0 defaults to 1",
			interval: "month",
			count:    0,
			now:      time.Date(2024, 5, 10, 8, 0, 0, 0, time.UTC),
			wantTime: time.Date(2024, 6, 10, 8, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "month count=-1 defaults to 1",
			interval: "month",
			count:    -1,
			now:      time.Date(2024, 5, 10, 8, 0, 0, 0, time.UTC),
			wantTime: time.Date(2024, 6, 10, 8, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "month with day clamping - Jan 31 + 3 months",
			interval: "month",
			count:    3,
			now:      time.Date(2024, 1, 31, 12, 0, 0, 0, time.UTC),
			wantTime: time.Date(2024, 4, 30, 12, 0, 0, 0, time.UTC), // April has 30 days
			wantErr:  false,
		},
		{
			name:     "invalid interval returns error",
			interval: "invalid",
			count:    1,
			now:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			wantTime: time.Time{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateNextResetWithCount(tt.interval, tt.count, tt.now)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !got.Equal(tt.wantTime) {
				t.Fatalf("expected time %v, got %v", tt.wantTime, got)
			}
		})
	}
}

func TestResolveEffectiveInterval(t *testing.T) {
	tests := []struct {
		name             string
		itemInterval     string
		fallbackInterval string
		want             string
	}{
		{
			name:             "inherit returns fallback",
			itemInterval:     "inherit",
			fallbackInterval: "month",
			want:             "month",
		},
		{
			name:             "empty returns fallback",
			itemInterval:     "",
			fallbackInterval: "year",
			want:             "year",
		},
		{
			name:             "specific interval returns itself",
			itemInterval:     "week",
			fallbackInterval: "month",
			want:             "week",
		},
		{
			name:             "month interval returns month",
			itemInterval:     "month",
			fallbackInterval: "year",
			want:             "month",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveEffectiveInterval(tt.itemInterval, tt.fallbackInterval)
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
