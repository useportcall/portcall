package subscription

import (
	"testing"
	"time"
)

func TestNextReset_Week(t *testing.T) {
	anchor := time.Date(2025, 1, 6, 9, 15, 0, 0, time.UTC) // Monday
	now := time.Date(2026, 2, 8, 10, 0, 0, 0, time.UTC)    // Sunday
	next, err := NextReset(anchor, "week", now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2026, 2, 9, 9, 15, 0, 0, time.UTC)
	if !next.Equal(want) {
		t.Fatalf("want %s, got %s", want, next)
	}
}

func TestNextReset_UnsupportedInterval(t *testing.T) {
	_, err := NextReset(time.Now(), "quarter", time.Now())
	if err == nil {
		t.Fatal("expected error for unsupported interval")
	}
}
