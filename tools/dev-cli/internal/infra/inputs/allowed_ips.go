package inputs

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ResolveAllowedIPs(explicit []string, sourcePath string) ([]string, error) {
	if len(explicit) > 0 {
		return UniqueIPs(explicit), nil
	}
	ips, err := ReadAdminAllowedIPs(sourcePath)
	if err != nil {
		return nil, err
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("no admin.ingress.allowedIPs found in %s", sourcePath)
	}
	return ips, nil
}

func ReadAdminAllowedIPs(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open values file: %w", err)
	}
	defer f.Close()
	inAdmin, inIngress, inList := false, false, false
	out := []string{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		trim := strings.TrimSpace(line)
		top := !strings.HasPrefix(line, " ") && strings.HasSuffix(trim, ":")
		if top {
			inAdmin, inIngress, inList = trim == "admin:", false, false
			continue
		}
		if inAdmin && trim == "ingress:" {
			inIngress, inList = true, false
			continue
		}
		if inAdmin && inIngress && strings.HasPrefix(trim, "allowedIPs:") {
			inList = true
			continue
		}
		if inList && strings.HasPrefix(trim, "- ") {
			v := strings.Trim(strings.TrimPrefix(trim, "- "), `"'`)
			if v != "" {
				out = append(out, v)
			}
			continue
		}
		if inList && trim != "" && !strings.HasPrefix(trim, "#") {
			inList = false
		}
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("scan values file: %w", err)
	}
	return UniqueIPs(out), nil
}

func UniqueIPs(in []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, value := range in {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" || seen[trimmed] {
			continue
		}
		seen[trimmed] = true
		out = append(out, trimmed)
	}
	return out
}
