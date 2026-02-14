package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func listSecrets(ns string) error {
	raw, err := runCmdOut("kubectl", "get", "secrets", "-n", ns)
	if err != nil {
		return fmt.Errorf("kubectl get secrets: %w\n%s", err, raw)
	}
	if strings.TrimSpace(raw) == "" {
		info("No secrets found in namespace %s", ns)
		return nil
	}
	info("Secrets in namespace %s:", ns)
	fmt.Println()
	fmt.Println(raw)
	return nil
}

func getSecret(ns, name string) error {
	out, err := runCmdOut(
		"kubectl", "get", "secret", name, "-n", ns, "-o", "json",
	)
	if err != nil {
		return fmt.Errorf("secret %q not found in namespace %s:\n%s", name, ns, out)
	}
	var raw struct {
		Data map[string]string `json:"data"`
	}
	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		return fmt.Errorf("parse secret JSON: %w", err)
	}
	info("Secret: %s  (namespace: %s)", name, ns)
	if len(raw.Data) == 0 {
		warn("  (no keys)")
		return nil
	}
	keys := sortedKeys(raw.Data)
	fmt.Printf("  %-40s %s\n", "KEY", "VALUE (decoded)")
	fmt.Printf("  %-40s %s\n", strings.Repeat("─", 40), strings.Repeat("─", 40))
	for _, k := range keys {
		fmt.Printf("  %-40s %s\n", k, decodeB64(raw.Data[k]))
	}
	return nil
}

func setSecretKeys(ns, name string, pairs []string) error {
	kv, err := parsePairs(pairs)
	if err != nil {
		return err
	}
	if secretExists(ns, name) {
		return patchSecretKeys(ns, name, kv)
	}
	return createSecret(ns, name, kv)
}

func secretExists(ns, name string) bool {
	_, err := runCmdOut("kubectl", "get", "secret", name, "-n", ns, "-o", "name")
	return err == nil
}

func fetchSecretData(ns, name string) (map[string]string, error) {
	out, err := runCmdOut("kubectl", "get", "secret", name, "-n", ns, "-o", "json")
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
