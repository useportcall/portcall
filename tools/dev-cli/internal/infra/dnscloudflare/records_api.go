package dnscloudflare

import (
	"fmt"
	"net/url"
	"strings"
)

func (c *Client) listRecords(zoneID, host string) ([]DNSRecord, error) {
	var records []DNSRecord
	path := fmt.Sprintf("/zones/%s/dns_records?name=%s&per_page=100", zoneID, url.QueryEscape(host))
	if err := c.request("GET", path, nil, &records); err != nil {
		return nil, fmt.Errorf("list dns records for %q: %w", host, err)
	}
	return records, nil
}

func pickRecord(records []DNSRecord, recordType string) *DNSRecord {
	for i := range records {
		if strings.EqualFold(records[i].Type, recordType) {
			return &records[i]
		}
	}
	if len(records) > 0 {
		return &records[0]
	}
	return nil
}

func (c *Client) createRecord(zoneID, host, recordType, target string) error {
	body := map[string]any{"type": recordType, "name": host, "content": target, "ttl": 1, "proxied": false}
	if err := c.request("POST", "/zones/"+zoneID+"/dns_records", body, nil); err != nil {
		return fmt.Errorf("create dns record for %q: %w", host, err)
	}
	return nil
}

func (c *Client) updateRecord(zoneID, recordID, host, recordType, target string) error {
	body := map[string]any{"type": recordType, "name": host, "content": target, "ttl": 1, "proxied": false}
	path := fmt.Sprintf("/zones/%s/dns_records/%s", zoneID, recordID)
	if err := c.request("PUT", path, body, nil); err != nil {
		return fmt.Errorf("update dns record for %q: %w", host, err)
	}
	return nil
}
