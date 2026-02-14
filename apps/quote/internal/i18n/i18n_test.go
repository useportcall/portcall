package i18n

import (
	"net/http/httptest"
	"testing"
)

func newI18nForTest(t *testing.T) *I18n {
	t.Helper()
	i := &I18n{translations: map[string]map[string]interface{}{}}
	i.loadTranslations()
	return i
}

func TestGetLanguageFromQuery(t *testing.T) {
	i := newI18nForTest(t)
	req := httptest.NewRequest("GET", "/quotes/q_123?lang=ja-JP", nil)
	if got := i.GetLanguage(req); got != "ja" {
		t.Fatalf("expected ja, got %q", got)
	}
}

func TestGetLanguageFallbackForUnknownQuery(t *testing.T) {
	i := newI18nForTest(t)
	req := httptest.NewRequest("GET", "/quotes/q_123?lang=xx", nil)
	req.Header.Set("Accept-Language", "ja-JP,ja;q=0.9")
	if got := i.GetLanguage(req); got != "ja" {
		t.Fatalf("expected ja fallback from header, got %q", got)
	}
}

func TestGetLanguageFromHeader(t *testing.T) {
	i := newI18nForTest(t)
	req := httptest.NewRequest("GET", "/quotes/q_123", nil)
	req.Header.Set("Accept-Language", "ja-JP,ja;q=0.9,en;q=0.8")
	if got := i.GetLanguage(req); got != "ja" {
		t.Fatalf("expected ja, got %q", got)
	}
}

func TestTFallbackToEnglish(t *testing.T) {
	i := newI18nForTest(t)
	if got := i.T("xx", "quote_details"); got != "Quote details" {
		t.Fatalf("expected English fallback, got %q", got)
	}
}
