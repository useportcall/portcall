package subscription

import (
	"fmt"
	"time"
)

// NextReset calculates the next reset time based on the interval and anchor.
func NextReset(anchor time.Time, interval string, now time.Time) (*time.Time, error) {
	var recent time.Time
	switch interval {
	case "hour":
		recent = time.Date(now.Year(), now.Month(), now.Day(),
			now.Hour(), anchor.Minute(), anchor.Second(), 0, now.Location())
		recent = recent.Add(time.Hour)
	case "day":
		recent = time.Date(now.Year(), now.Month(), now.Day(),
			anchor.Hour(), anchor.Minute(), anchor.Second(), 0, now.Location())
		recent = recent.AddDate(0, 0, 1)
	case "week":
		recent = time.Date(now.Year(), now.Month(), now.Day(),
			anchor.Hour(), anchor.Minute(), anchor.Second(), 0, now.Location())
		days := int(anchor.Weekday() - recent.Weekday())
		recent = recent.AddDate(0, 0, days)
		if !recent.After(now) {
			recent = recent.AddDate(0, 0, 7)
		}
	case "month":
		recent = time.Date(now.Year(), now.Month(), anchor.Day(),
			anchor.Hour(), anchor.Minute(), anchor.Second(), 0, now.Location())
		recent = addMonthClamped(recent)
	case "year":
		recent = time.Date(now.Year(), anchor.Month(), anchor.Day(),
			anchor.Hour(), anchor.Minute(), anchor.Second(), 0, now.Location())
		recent = recent.AddDate(1, 0, 0)
	case "per_use":
		recent = now
	case "no_reset", "none":
		recent = time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)
	default:
		return nil, fmt.Errorf("unsupported interval: %s", interval)
	}
	return &recent, nil
}

// addMonthClamped adds one month, clamping to the last day of the month.
func addMonthClamped(t time.Time) time.Time {
	year, month, day := t.Date()
	loc := t.Location()

	nextMonth := month + 1
	nextYear := year
	if nextMonth > 12 {
		nextMonth = 1
		nextYear++
	}

	firstOfNextNext := time.Date(nextYear, nextMonth+1, 1, 0, 0, 0, 0, loc)
	lastDay := firstOfNextNext.AddDate(0, 0, -1).Day()

	if day > lastDay {
		day = lastDay
	}

	return time.Date(nextYear, nextMonth, day,
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
}
