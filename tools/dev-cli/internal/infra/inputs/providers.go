package inputs

import (
	"fmt"
	"strings"
)

var supportedProviders = map[string]bool{"digitalocean": true}

func ResolveProviders(providers []string, fallback string) ([]string, error) {
	selected := providers
	if len(selected) == 0 && strings.TrimSpace(fallback) != "" {
		selected = []string{fallback}
	}
	if len(selected) == 0 {
		selected = []string{"digitalocean"}
	}
	unique := []string{}
	seen := map[string]bool{}
	for _, provider := range selected {
		name := strings.ToLower(strings.TrimSpace(provider))
		if name == "" || seen[name] {
			continue
		}
		if !supportedProviders[name] {
			return nil, fmt.Errorf("infra init currently supports provider(s): digitalocean")
		}
		seen[name] = true
		unique = append(unique, name)
	}
	if len(unique) == 0 {
		return nil, fmt.Errorf("at least one provider must be selected")
	}
	return unique, nil
}
