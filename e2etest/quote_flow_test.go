package e2etest

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestE2E_QuoteIssueAcceptAndSignatureRetrieval(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}
	h := NewHarness(t)

	quote := h.DashPost(dashPath(h.AppPublicID, "quotes"), nil).
		MustOK(t, "create quote")
	quoteID := mustString(t, quote, "id")
	plan := quote["plan"].(map[string]any)
	planID := mustString(t, plan, "id")

	userID := mustString(t, h.DashPost(dashPath(h.AppPublicID, "users"), map[string]any{
		"email": "quote-go-e2e@test.dev",
	}).MustOK(t, "create user"), "id")

	items := h.DashGet(dashPath(h.AppPublicID, "plan-items?plan_id="+planID)).
		MustOKList(t, "list plan items")
	itemID := mustString(t, items[0], "id")
	h.DashPost(dashPath(h.AppPublicID, "plan-items/"+itemID), map[string]any{
		"unit_amount": 2999,
	}).MustOK(t, "update quote amount")
	h.DashPost(dashPath(h.AppPublicID, "quotes/"+quoteID), map[string]any{
		"user_id":                 userID,
		"recipient_email":         "recipient-go-e2e@test.dev",
		"recipient_name":          "Recipient",
		"recipient_title":         "CEO",
		"company_name":            "Recipient Co",
		"direct_checkout_enabled": false,
	}).MustOK(t, "update quote")
	h.DashPost(dashPath(h.AppPublicID, "quotes/"+quoteID+"/send"), nil).
		MustOK(t, "send quote")

	quoteData := h.DashGet(dashPath(h.AppPublicID, "quotes/"+quoteID)).
		MustOK(t, "get quote")
	quoteURL := mustString(t, quoteData, "url")

	form := url.Values{}
	form.Set("signatureData", "data:image/png;base64,"+validSignatureBase64(t))
	req, err := http.NewRequest(http.MethodPost, quoteURL, strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("create submit request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("submit quote: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("expected redirect on quote submit, got %d", resp.StatusCode)
	}

	accepted := h.DashGet(dashPath(h.AppPublicID, "quotes/"+quoteID)).
		MustOK(t, "get accepted quote")
	if getString(accepted, "status") != "accepted" {
		t.Fatalf("expected accepted status, got %v", accepted["status"])
	}

	signatureReq, err := http.NewRequest(
		http.MethodGet,
		h.Dashboard.URL+dashPath(h.AppPublicID, "quotes/"+quoteID+"/signature"),
		nil,
	)
	if err != nil {
		t.Fatalf("create signature request: %v", err)
	}
	for k, v := range h.DashHeaders() {
		signatureReq.Header.Set(k, v)
	}
	signatureResp, err := http.DefaultClient.Do(signatureReq)
	if err != nil {
		t.Fatalf("download signature: %v", err)
	}
	defer signatureResp.Body.Close()
	if signatureResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for signature download, got %d", signatureResp.StatusCode)
	}
}

func validSignatureBase64(t *testing.T) string {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 10, G: 10, B: 10, A: 255})
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode signature image: %v", err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}
