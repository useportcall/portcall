package utils

import (
	"errors"
	"time"
)

// CalculateNextReset calculates the next reset time based on interval.
// Deprecated: Use CalculateNextResetWithCount for interval count support.
func CalculateNextReset(interval string, now time.Time) (error, time.Time) {
	nextReset, err := CalculateNextResetWithCount(interval, 1, now)
	return err, nextReset
}

// CalculateNextResetWithCount calculates the next reset time based on interval and count.
// For example, interval="month" and count=3 gives a quarterly billing cycle.
func CalculateNextResetWithCount(interval string, count int, now time.Time) (time.Time, error) {
	if count <= 0 {
		count = 1
	}

	var nextReset time.Time
	switch interval {
	case "week":
		nextReset = now.AddDate(0, 0, 7*count)
	case "month":
		nextReset = addMonths(now, count)
	case "year":
		nextReset = addYears(now, count)
	default:
		return time.Time{}, errors.New("invalid interval: " + interval)
	}

	return nextReset, nil
}
