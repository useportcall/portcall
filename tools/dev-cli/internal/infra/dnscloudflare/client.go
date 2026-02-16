package dnscloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const apiBaseURL = "https://api.cloudflare.com/client/v4"

type Client struct {
	token      string
	baseURL    string
	httpClient *http.Client
}

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewClient(token string) *Client {
	return NewClientWithHTTP(token, apiBaseURL, http.DefaultClient)
}

func NewClientWithHTTP(token string, baseURL string, httpClient *http.Client) *Client {
	return &Client{
		token:      strings.TrimSpace(token),
		baseURL:    strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		httpClient: httpClient,
	}
}

func (c *Client) request(method, path string, in any, out any) error {
	var body *bytes.Reader
	if in == nil {
		body = bytes.NewReader(nil)
	} else {
		buf, err := json.Marshal(in)
		if err != nil {
			return fmt.Errorf("encode request: %w", err)
		}
		body = bytes.NewReader(buf)
	}
	req, err := http.NewRequest(method, c.baseURL+path, body)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	var raw struct {
		Success bool            `json:"success"`
		Errors  []apiError      `json:"errors"`
		Result  json.RawMessage `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	if !raw.Success {
		return fmt.Errorf("cloudflare api error: %s", joinErrors(raw.Errors))
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(raw.Result, out); err != nil {
		return fmt.Errorf("decode result: %w", err)
	}
	return nil
}

func joinErrors(errors []apiError) string {
	parts := []string{}
	for _, e := range errors {
		if e.Message == "" {
			continue
		}
		parts = append(parts, e.Message)
	}
	if len(parts) == 0 {
		return "unknown error"
	}
	return strings.Join(parts, "; ")
}
