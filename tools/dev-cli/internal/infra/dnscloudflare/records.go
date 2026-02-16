package dnscloudflare

import (
	"fmt"
	"net"
	"strings"
)

type DNSRecord struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type RecordOutcome struct {
	Host   string
	Status string
	Detail string
}

func (c *Client) EnsureRecords(zoneID string, hosts []string, target string, autoApply bool) ([]RecordOutcome, error) {
	recordType := "CNAME"
	if net.ParseIP(strings.TrimSpace(target)) != nil {
		recordType = "A"
	}
	out := []RecordOutcome{}
	for _, host := range hosts {
		item, err := c.ensureRecord(zoneID, host, recordType, target, autoApply)
		if err != nil {
			return out, err
		}
		out = append(out, item)
	}
	return out, nil
}

func (c *Client) ensureRecord(zoneID, host, recordType, target string, autoApply bool) (RecordOutcome, error) {
	records, err := c.listRecords(zoneID, host)
	if err != nil {
		return RecordOutcome{}, err
	}
	existing := pickRecord(records, recordType)
	if existing != nil && strings.TrimSpace(existing.Content) == strings.TrimSpace(target) {
		return RecordOutcome{Host: host, Status: "ok", Detail: "already configured"}, nil
	}
	if !autoApply {
		if existing == nil {
			return RecordOutcome{Host: host, Status: "missing", Detail: fmt.Sprintf("create %s -> %s", recordType, target)}, nil
		}
		return RecordOutcome{Host: host, Status: "mismatch", Detail: fmt.Sprintf("%s -> %s (expected %s)", existing.Type, existing.Content, target)}, nil
	}
	if existing == nil {
		if err := c.createRecord(zoneID, host, recordType, target); err != nil {
			return RecordOutcome{}, err
		}
		return RecordOutcome{Host: host, Status: "created", Detail: fmt.Sprintf("%s -> %s", recordType, target)}, nil
	}
	if err := c.updateRecord(zoneID, existing.ID, host, recordType, target); err != nil {
		return RecordOutcome{}, err
	}
	return RecordOutcome{Host: host, Status: "updated", Detail: fmt.Sprintf("%s -> %s", recordType, target)}, nil
}
