package utils

import "time"

// addMonths adds the specified number of months to a time, clamping to valid days.
func addMonths(t time.Time, months int) time.Time {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	nsec := t.Nanosecond()
	loc := t.Location()

	// Calculate target month and year
	totalMonths := int(month) - 1 + months
	targetYear := year + totalMonths/12
	targetMonth := time.Month(totalMonths%12 + 1)
	if totalMonths%12 < 0 {
		targetMonth = time.Month(totalMonths%12 + 13)
		targetYear--
	}

	// Find the last day of the target month
	firstOfNextMonth := time.Date(targetYear, targetMonth+1, 1, 0, 0, 0, 0, loc)
	lastDayOfMonth := firstOfNextMonth.AddDate(0, 0, -1).Day()

	// Clamp the day
	if day > lastDayOfMonth {
		day = lastDayOfMonth
	}

	return time.Date(targetYear, targetMonth, day, hour, min, sec, nsec, loc)
}

// addYears adds the specified number of years to a time, handling leap years.
func addYears(t time.Time, years int) time.Time {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	nsec := t.Nanosecond()
	loc := t.Location()

	targetYear := year + years

	// Find the last day of the same month in the target year
	firstOfNextMonth := time.Date(targetYear, month+1, 1, 0, 0, 0, 0, loc)
	lastDayOfMonth := firstOfNextMonth.AddDate(0, 0, -1).Day()

	// Clamp the day (e.g., Feb 29 -> Feb 28 in non-leap years)
	if day > lastDayOfMonth {
		day = lastDayOfMonth
	}

	return time.Date(targetYear, month, day, hour, min, sec, nsec, loc)
}

// ResolveEffectiveInterval resolves the effective interval for an item.
// If the item's interval is "inherit", it uses the fallback interval.
func ResolveEffectiveInterval(itemInterval, fallbackInterval string) string {
	if itemInterval == "" || itemInterval == "inherit" {
		return fallbackInterval
	}
	return itemInterval
}

// FindEarliestTime returns the earliest non-nil time from a list of times.
// Returns nil if no valid times are provided.
func FindEarliestTime(times ...*time.Time) *time.Time {
	var earliest *time.Time
	for _, t := range times {
		if t == nil {
			continue
		}
		if earliest == nil || t.Before(*earliest) {
			earliest = t
		}
	}
	return earliest
}
