package dnscloudflare

import (
	"fmt"
	"net/url"
	"strings"
)

type Zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) ResolveZoneID(domain string, explicitZoneID string) (string, error) {
	if strings.TrimSpace(explicitZoneID) != "" {
		return strings.TrimSpace(explicitZoneID), nil
	}
	domain = strings.Trim(strings.ToLower(strings.TrimSpace(domain)), ".")
	if domain == "" {
		return "", fmt.Errorf("domain is required for cloudflare zone lookup")
	}
	for _, candidate := range zoneCandidates(domain) {
		id, err := c.zoneIDByName(candidate)
		if err != nil {
			return "", err
		}
		if id != "" {
			return id, nil
		}
	}
	return "", fmt.Errorf("no cloudflare zone found for domain %q", domain)
}

func zoneCandidates(domain string) []string {
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return []string{domain}
	}
	out := []string{}
	for i := 0; i <= len(parts)-2; i++ {
		out = append(out, strings.Join(parts[i:], "."))
	}
	return out
}

func (c *Client) zoneIDByName(name string) (string, error) {
	var zones []Zone
	path := "/zones?name=" + url.QueryEscape(name) + "&status=active&per_page=50"
	if err := c.request("GET", path, nil, &zones); err != nil {
		return "", fmt.Errorf("lookup zone %q: %w", name, err)
	}
	for _, zone := range zones {
		if strings.EqualFold(strings.TrimSpace(zone.Name), name) {
			return strings.TrimSpace(zone.ID), nil
		}
	}
	return "", nil
}
