package harness

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

// JSONResponse wraps a parsed JSON response with status code.
type JSONResponse struct {
	Status   int
	Data     map[string]any   // contents of "data" key if it's an object
	DataList []map[string]any // contents of "data" key if it's a list
	Raw      map[string]any   // full response body
}

// DoJSON sends an HTTP request with optional JSON body and parses response.
func DoJSON(t *testing.T, method, url string, body any, headers map[string]string) *JSONResponse {
	t.Helper()

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request body: %v", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("HTTP %s %s: %v", method, url, err)
	}
	defer resp.Body.Close()

	raw := make(map[string]any)
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		t.Fatalf("decode response from %s: %v", url, err)
	}

	jr := &JSONResponse{Status: resp.StatusCode, Raw: raw}
	if data, ok := raw["data"]; ok {
		switch v := data.(type) {
		case map[string]any:
			jr.Data = v
		case []any:
			for _, item := range v {
				if m, ok := item.(map[string]any); ok {
					jr.DataList = append(jr.DataList, m)
				}
			}
		}
	}
	return jr
}

// MustOK asserts the response is 200 and returns the data map.
func (r *JSONResponse) MustOK(t *testing.T, context string) map[string]any {
	t.Helper()
	if r.Status != http.StatusOK {
		t.Fatalf("%s: expected 200, got %d — %v", context, r.Status, r.Raw)
	}
	if r.Data == nil {
		t.Fatalf("%s: response missing 'data' key: %v", context, r.Raw)
	}
	return r.Data
}

// MustOKList asserts the response is 200 and returns the data list.
func (r *JSONResponse) MustOKList(t *testing.T, context string) []map[string]any {
	t.Helper()
	if r.Status != http.StatusOK {
		t.Fatalf("%s: expected 200, got %d — %v", context, r.Status, r.Raw)
	}
	return r.DataList
}

// MustStatus asserts the response has the given status code.
func (r *JSONResponse) MustStatus(t *testing.T, status int, context string) {
	t.Helper()
	if r.Status != status {
		t.Fatalf("%s: expected %d, got %d — %v", context, status, r.Status, r.Raw)
	}
}

// GetString extracts a string from nested map data.
func GetString(m map[string]any, key string) string {
	v, _ := m[key].(string)
	return v
}

// GetSlice extracts a slice from nested map data.
func GetSlice(m map[string]any, key string) []any {
	v, _ := m[key].([]any)
	return v
}

// GetFloat extracts a float64 from nested map data.
func GetFloat(m map[string]any, key string) float64 {
	v, _ := m[key].(float64)
	return v
}

// MustString extracts a string or fails the test.
func MustString(t *testing.T, m map[string]any, key string) string {
	t.Helper()
	v, ok := m[key].(string)
	if !ok || v == "" {
		t.Fatalf("expected non-empty string for key %q, got %v", key, m[key])
	}
	return v
}

// DashPath builds a dashboard API path for the given app and subpath.
func DashPath(appID, path string) string {
	return fmt.Sprintf("/api/apps/%s/%s", appID, path)
}
