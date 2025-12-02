// calculate_next_reset_test.go
package utils

import (
	"errors"
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
			// If both non-nil, compare messages
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Fatalf("expected error %q, got %q", tt.wantErr.Error(), err.Error())
			}

			// Check time result
			if !got.Equal(tt.wantTime) {
				t.Fatalf("expected time %v, got %v", tt.wantTime, got)
			}
		})
	}
}
