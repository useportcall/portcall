package handlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/useportcall/portcall/libs/go/dbx/models"
	"github.com/useportcall/portcall/libs/go/qx/server"
)

type EndSubscriptionResetPayload struct {
	SubscriptionID uint `json:"subscription_id"`
}

func EndSubscriptionReset(c server.IContext) error {
	var p EndSubscriptionResetPayload
	if err := json.Unmarshal(c.Payload(), &p); err != nil {
		return err
	}

	var subscription models.Subscription
	if err := c.DB().FindForID(p.SubscriptionID, &subscription); err != nil {
		return err
	}

	now := time.Now()

	nextReset, err := NextReset(subscription.CreatedAt, subscription.BillingInterval, now)
	if err != nil {
		return fmt.Errorf("failed to calculate next reset: %w", err)
	}

	subscription.LastResetAt = now
	subscription.NextResetAt = *nextReset

	if err := c.DB().Save(&subscription); err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	// TODO: send notification that subscription has been reset?

	return nil
}

// TODO: add tests
func NextReset(anchor time.Time, interval string, now time.Time) (*time.Time, error) {
	var recent time.Time
	switch interval {
	case "hour":
		// Keep the same minute and second, but find the most recent hour relative to `now`
		recent = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), anchor.Minute(), anchor.Second(), 0, now.Location())
		recent = recent.Add(time.Hour)
	case "day":
		// Keep the same hour, minute, and second, but find the most recent day relative to `now`
		recent = time.Date(now.Year(), now.Month(), now.Day(), anchor.Hour(), anchor.Minute(), anchor.Second(), 0, now.Location())
		recent = recent.AddDate(0, 0, 1)
	case "month":
		// TODO: test for leap year, 29/2, 31/1, 31/3, 28/2, etc.
		// Keep the same day, hour, minute, and second, but find the most recent month relative to `now`
		recent = time.Date(now.Year(), now.Month(), anchor.Day(), anchor.Hour(), anchor.Minute(), anchor.Second(), 0, now.Location())
		recent = AddMonthClamped(recent)
	case "year":
		// Keep the same month, day, hour, minute, and second, but find the most recent year relative to `now`
		recent = time.Date(now.Year(), anchor.Month(), anchor.Day(), anchor.Hour(), anchor.Minute(), anchor.Second(), 0, now.Location())
		recent = recent.AddDate(1, 0, 0)
	case "per_use":
		recent = now
	case "no_reset":
		// long time in the future
		recent = time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC)
	default:
		return nil, fmt.Errorf("unsupported interval: %s", interval)
	}

	return &recent, nil
}

// TODO: add tests
func AddMonthClamped(t time.Time) time.Time {
	year, month, day := t.Date()
	location := t.Location()

	// Add one month manually
	nextMonth := month + 1
	nextYear := year
	if nextMonth > 12 {
		nextMonth = 1
		nextYear++
	}

	// Find the last day of the next month
	firstOfNextNextMonth := time.Date(nextYear, nextMonth+1, 1, 0, 0, 0, 0, location)
	lastDayOfNextMonth := firstOfNextNextMonth.AddDate(0, 0, -1).Day()

	// Clamp the day
	if day > lastDayOfNextMonth {
		day = lastDayOfNextMonth
	}

	return time.Date(nextYear, nextMonth, day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), location)
}
