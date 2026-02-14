package ratelimitx

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// ParseRateFromEnv parses a rate config from an environment variable.
// Format: "100/1m" or "1000/1h"
func ParseRateFromEnv(envKey string, defaultLimit int, defaultWindow time.Duration) Config {
	rateStr := os.Getenv(envKey)
	if rateStr == "" {
		return Config{Limit: defaultLimit, Window: defaultWindow}
	}

	var limit int
	var windowStr string
	_, err := fmt.Sscanf(rateStr, "%d/%s", &limit, &windowStr)
	if err != nil {
		log.Printf("Invalid rate format in %s: %s, using defaults", envKey, rateStr)
		return Config{Limit: defaultLimit, Window: defaultWindow}
	}

	window, err := parseDuration(windowStr)
	if err != nil {
		log.Printf("Invalid duration in %s: %s, using defaults", envKey, windowStr)
		return Config{Limit: defaultLimit, Window: defaultWindow}
	}

	return Config{Limit: limit, Window: window}
}

func parseDuration(s string) (time.Duration, error) {
	if len(s) < 2 {
		return 0, fmt.Errorf("invalid duration: %s", s)
	}

	value, err := strconv.Atoi(s[:len(s)-1])
	if err != nil {
		return 0, err
	}

	unit := s[len(s)-1:]
	switch unit {
	case "s":
		return time.Duration(value) * time.Second, nil
	case "m":
		return time.Duration(value) * time.Minute, nil
	case "h":
		return time.Duration(value) * time.Hour, nil
	default:
		return 0, fmt.Errorf("unsupported duration unit: %s", unit)
	}
}
