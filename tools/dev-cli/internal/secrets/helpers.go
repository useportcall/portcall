package secrets

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type secretData struct {
	Data map[string]string `json:"data"`
}

func parsePairs(pairs []string) (map[string]string, error) {
	kv := map[string]string{}
	for _, p := range pairs {
		idx := strings.Index(p, "=")
		if idx < 1 {
			return nil, fmt.Errorf("invalid key=value pair: %q (expected KEY=VALUE)", p)
		}
		kv[p[:idx]] = p[idx+1:]
	}
	return kv, nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func decodeB64(s string) string {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "(decode error)"
	}
	return string(decoded)
}

func buildB64Patch(kv map[string]string) string {
	data := map[string]string{}
	for k, v := range kv {
		data[k] = base64.StdEncoding.EncodeToString([]byte(v))
	}
	b, _ := json.Marshal(map[string]any{"data": data})
	return string(b)
}

func buildNullPatch(keys []string) string {
	data := map[string]any{}
	for _, k := range keys {
		data[k] = nil
	}
	b, _ := json.Marshal(map[string]any{"data": data})
	return string(b)
}
