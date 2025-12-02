package utils

import (
	"errors"
	"time"
)

func CalculateNextReset(interval string, now time.Time) (error, time.Time) {
	var nextReset time.Time
	switch interval {
	case "week":
		nextReset = now.AddDate(0, 0, 7)
	case "month":
		// Calculate next month, but adjust day if needed
		year, month, day := now.Date()
		nextMonth := month + 1
		nextYear := year
		if nextMonth > 12 {
			nextMonth = 1
			nextYear++
		}
		// Find the last day of the next month
		firstOfNextMonth := time.Date(nextYear, nextMonth, 1, now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), now.Location())
		lastDay := firstOfNextMonth.AddDate(0, 1, -1).Day()
		if day > lastDay {
			day = lastDay
		}
		nextReset = time.Date(nextYear, nextMonth, day, now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), now.Location())
	case "year":
		// Calculate next year, adjust for leap years if needed
		year, month, day := now.Date()
		nextYear := year + 1
		// Find the last day of the same month in the next year
		firstOfNextMonth := time.Date(nextYear, month, 1, now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), now.Location())
		lastDay := firstOfNextMonth.AddDate(0, 1, -1).Day()
		if day > lastDay {
			day = lastDay
		}
		nextReset = time.Date(nextYear, month, day, now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), now.Location())
	default:
		return errors.New("invalid interval"), time.Time{}
	}

	return nil, nextReset
}
