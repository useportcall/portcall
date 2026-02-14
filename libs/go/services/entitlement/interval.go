package entitlement

import "time"

func nextResetAt(interval string, now time.Time) *time.Time {
	next := now
	switch interval {
	case "hour":
		next = now.Add(time.Hour)
	case "day":
		next = now.AddDate(0, 0, 1)
	case "week":
		next = now.AddDate(0, 0, 7)
	case "month":
		next = now.AddDate(0, 1, 0)
	case "year":
		next = now.AddDate(1, 0, 0)
	case "no_reset", "none":
		next = time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	return &next
}
