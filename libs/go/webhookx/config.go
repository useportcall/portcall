package webhookx

import (
	"net"
	"os"
	"strconv"
	"strings"
)

func loadPositiveInt(name string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func parseCIDRList(raw string) []*net.IPNet {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]*net.IPNet, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item == "" {
			continue
		}
		if !strings.Contains(item, "/") {
			item += "/32"
		}
		_, cidr, err := net.ParseCIDR(item)
		if err == nil {
			out = append(out, cidr)
		}
	}
	return out
}
