package envx

import (
	"os"
	"strings"
)

type Mode string

const (
	Production  Mode = "production"
	Development Mode = "development"
	Test        Mode = "test"
)

func CurrentMode() Mode {
	raw := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	switch raw {
	case "prod", "production":
		return Production
	case "test":
		return Test
	default:
		return Development
	}
}

func IsProd() bool { return CurrentMode() == Production }
func IsDev() bool  { return CurrentMode() == Development }
func IsTest() bool { return CurrentMode() == Test }
