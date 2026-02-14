package main

import (
	"fmt"
	"sort"
	"strings"
)

func createSecret(ns, name string, kv map[string]string) error {
	args := []string{"create", "secret", "generic", name, "-n", ns}
	for k, v := range kv {
		args = append(args, fmt.Sprintf("--from-literal=%s=%s", k, v))
	}
	info("Creating secret %s with %d key(s)", name, len(kv))
	if err := runCmd("kubectl", args...); err != nil {
		return fmt.Errorf("create secret: %w", err)
	}
	ok("Created secret %s", name)
	return nil
}

func patchSecretKeys(ns, name string, kv map[string]string) error {
	existing, err := fetchSecretData(ns, name)
	if err != nil {
		return err
	}

	var created, updated []string
	for k := range kv {
		if _, exists := existing[k]; exists {
			updated = append(updated, k)
		} else {
			created = append(created, k)
		}
	}
	sort.Strings(created)
	sort.Strings(updated)

	if len(created) > 0 {
		info("  New keys:     %s", strings.Join(created, ", "))
	}
	if len(updated) > 0 {
		warn("  Overwriting:  %s", strings.Join(updated, ", "))
	}

	patch := buildB64Patch(kv)
	if err := mergeSecretPatch(ns, name, patch); err != nil {
		return err
	}
	ok("Updated secret %s (%d key(s))", name, len(kv))
	return nil
}

func deleteSecret(ns, name string) error {
	if !secretExists(ns, name) {
		return fmt.Errorf("secret %q does not exist in namespace %s", name, ns)
	}
	if !askYesNo(fmt.Sprintf("Delete entire secret %q? [y/N]: ", name), false) {
		return nil
	}
	if err := runCmd("kubectl", "delete", "secret", name, "-n", ns); err != nil {
		return fmt.Errorf("delete secret: %w", err)
	}
	ok("Deleted secret %s", name)
	return nil
}

func deleteSecretKeys(ns, name string, keys []string) error {
	if !secretExists(ns, name) {
		return fmt.Errorf("secret %q does not exist in namespace %s", name, ns)
	}
	existing, err := fetchSecretData(ns, name)
	if err != nil {
		return err
	}
	for _, k := range keys {
		if _, found := existing[k]; !found {
			return fmt.Errorf("key %q not found in secret %s", k, name)
		}
	}
	info("Removing keys from %s: %s", name, strings.Join(keys, ", "))
	if !askYesNo("Proceed? [y/N]: ", false) {
		return nil
	}
	patch := buildNullPatch(keys)
	if err := mergeSecretPatch(ns, name, patch); err != nil {
		return err
	}
	ok("Removed %d key(s) from secret %s", len(keys), name)
	return nil
}
