package secrets

import (
	"encoding/json"
	"fmt"
	"strings"
)

func list(ns string) error {
	raw, err := runOut("kubectl", "get", "secrets", "-n", ns)
	if err != nil {
		return fmt.Errorf("kubectl get secrets: %w\n%s", err, raw)
	}
	if strings.TrimSpace(raw) == "" {
		fmt.Printf("No secrets found in namespace %s\n", ns)
		return nil
	}
	fmt.Printf("Secrets in namespace %s:\n\n%s\n", ns, raw)
	return nil
}

func get(ns, name string) error {
	out, err := runOut("kubectl", "get", "secret", name, "-n", ns, "-o", "json")
	if err != nil {
		return fmt.Errorf("secret %q not found in namespace %s:\n%s", name, ns, out)
	}
	var raw struct {
		Data map[string]string `json:"data"`
	}
	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		return fmt.Errorf("parse secret JSON: %w", err)
	}
	fmt.Printf("Secret: %s  (namespace: %s)\n", name, ns)
	if len(raw.Data) == 0 {
		fmt.Println("  (no keys)")
		return nil
	}
	keys := sortedKeys(raw.Data)
	fmt.Printf("  %-40s %s\n", "KEY", "VALUE (decoded)")
	fmt.Printf("  %-40s %s\n", strings.Repeat("-", 40), strings.Repeat("-", 40))
	for _, k := range keys {
		fmt.Printf("  %-40s %s\n", k, decodeB64(raw.Data[k]))
	}
	return nil
}

func secretExists(ns, name string) bool {
	_, err := runOut("kubectl", "get", "secret", name, "-n", ns, "-o", "name")
	return err == nil
}

func fetchData(ns, name string) (map[string]string, error) {
	out, err := runOut("kubectl", "get", "secret", name, "-n", ns, "-o", "json")
	if err != nil {
		return nil, fmt.Errorf("get secret %s: %w", name, err)
	}
	var s secretData
	if err := json.Unmarshal([]byte(out), &s); err != nil {
		return nil, fmt.Errorf("parse secret: %w", err)
	}
	if s.Data == nil {
		s.Data = map[string]string{}
	}
	return s.Data, nil
}
