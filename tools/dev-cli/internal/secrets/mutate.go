package secrets

import (
	"fmt"
	"sort"
	"strings"
)

func setKeys(ns, name string, pairs []string) error {
	kv, err := parsePairs(pairs)
	if err != nil {
		return err
	}
	if secretExists(ns, name) {
		return patchKeys(ns, name, kv)
	}
	return create(ns, name, kv)
}

func create(ns, name string, kv map[string]string) error {
	args := []string{"create", "secret", "generic", name, "-n", ns}
	for k, v := range kv {
		args = append(args, fmt.Sprintf("--from-literal=%s=%s", k, v))
	}
	fmt.Printf("Creating secret %s with %d key(s)\n", name, len(kv))
	return run("kubectl", args...)
}

func patchKeys(ns, name string, kv map[string]string) error {
	existing, err := fetchData(ns, name)
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
		fmt.Printf("New keys: %s\n", strings.Join(created, ", "))
	}
	if len(updated) > 0 {
		fmt.Printf("Overwriting: %s\n", strings.Join(updated, ", "))
	}
	return run("kubectl", "patch", "secret", name, "-n", ns, "--type", "merge", "-p", buildB64Patch(kv))
}

func deleteSecret(ns, name string) error {
	if !secretExists(ns, name) {
		return fmt.Errorf("secret %q does not exist in namespace %s", name, ns)
	}
	if !askYesNo(fmt.Sprintf("Delete entire secret %q? [y/N]: ", name), false) {
		return nil
	}
	return run("kubectl", "delete", "secret", name, "-n", ns)
}

func deleteKeys(ns, name string, keys []string) error {
	if !secretExists(ns, name) {
		return fmt.Errorf("secret %q does not exist in namespace %s", name, ns)
	}
	existing, err := fetchData(ns, name)
	if err != nil {
		return err
	}
	for _, k := range keys {
		if _, found := existing[k]; !found {
			return fmt.Errorf("key %q not found in secret %s", k, name)
		}
	}
	fmt.Printf("Removing keys from %s: %s\n", name, strings.Join(keys, ", "))
	if !askYesNo("Proceed? [y/N]: ", false) {
		return nil
	}
	return run("kubectl", "patch", "secret", name, "-n", ns, "--type", "merge", "-p", buildNullPatch(keys))
}
