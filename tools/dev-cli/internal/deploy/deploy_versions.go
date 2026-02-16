package deploy

import (
	"fmt"
	"strconv"
	"strings"
)

func deployName(app string) string {
	switch app {
	case "billing":
		return "billing-worker"
	case "email":
		return "email-worker"
	case "file":
		return "file-api"
	default:
		return app
	}
}

func helmValueName(app string) string {
	switch app {
	case "billing":
		return "billingWorker"
	case "email":
		return "emailWorker"
	case "file":
		return "fileApi"
	default:
		return app
	}
}

func bumpVersion(current, bump string) (string, error) {
	v := strings.TrimPrefix(strings.TrimSpace(current), "v")
	parts := strings.Split(v, ".")
	for len(parts) < 3 {
		parts = append(parts, "0")
	}
	maj, e1 := strconv.Atoi(parts[0])
	min, e2 := strconv.Atoi(parts[1])
	pat, e3 := strconv.Atoi(parts[2])
	if e1 != nil || e2 != nil || e3 != nil {
		return "", fmt.Errorf("invalid version: %s", current)
	}
	switch bump {
	case "major":
		maj, min, pat = maj+1, 0, 0
	case "minor":
		min, pat = min+1, 0
	case "patch":
		pat++
	case "skip":
	default:
		return "", fmt.Errorf("invalid bump type: %s", bump)
	}
	return fmt.Sprintf("v%d.%d.%d", maj, min, pat), nil
}
