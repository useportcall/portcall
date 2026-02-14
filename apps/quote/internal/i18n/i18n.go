package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localeFS embed.FS

type I18n struct {
	translations map[string]map[string]interface{}
	matcher      language.Matcher
	mu           sync.RWMutex
}

var instance *I18n
var once sync.Once

func GetInstance() *I18n {
	once.Do(func() {
		instance = &I18n{
			translations: make(map[string]map[string]interface{}),
		}
		instance.loadTranslations()
	})
	return instance
}

func (i *I18n) loadTranslations() {
	files, err := localeFS.ReadDir("locales")
	if err != nil {
		fmt.Printf("Error reading locales directory: %v\n", err)
		return
	}

	var tags []language.Tag
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		langCode := strings.TrimSuffix(file.Name(), ".json")
		content, err := localeFS.ReadFile("locales/" + file.Name())
		if err != nil {
			fmt.Printf("Error reading locale file %s: %v\n", file.Name(), err)
			continue
		}

		var data map[string]interface{}
		if err := json.Unmarshal(content, &data); err != nil {
			fmt.Printf("Error unmarshaling locale file %s: %v\n", file.Name(), err)
			continue
		}

		i.translations[langCode] = data
		tag, err := language.Parse(langCode)
		if err == nil {
			tags = append(tags, tag)
		}
	}

	// Set English as default fallback
	tags = append(tags, language.English)
	i.matcher = language.NewMatcher(tags)
}

func (i *I18n) GetLanguage(r *http.Request) string {
	// 1. Check query param.
	if lang := i.normalizeLanguage(r.URL.Query().Get("lang")); lang != "" {
		return lang
	}

	// 2. Check header.
	accept := r.Header.Get("Accept-Language")
	tag, _ := language.MatchStrings(i.matcher, accept)
	if lang := i.normalizeLanguage(tag.String()); lang != "" {
		return lang
	}
	return "en"
}

func (i *I18n) T(lang string, key string, args ...interface{}) string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	normalized := i.normalizeLanguage(lang)
	if normalized == "" {
		normalized = "en"
	}

	val, ok := i.getValue(normalized, key)
	if !ok {
		// Fallback to en
		val, ok = i.getValue("en", key)
		if !ok {
			return key
		}
	}

	strVal, ok := val.(string)
	if !ok {
		return key
	}

	if len(args) > 0 {
		return fmt.Sprintf(strVal, args...)
	}
	return strVal
}

func (i *I18n) getValue(lang string, key string) (interface{}, bool) {
	data, ok := i.translations[lang]
	if !ok {
		return nil, false
	}

	// Support dot notation: home.title
	parts := strings.Split(key, ".")
	var current interface{} = data
	for _, part := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil, false
		}
		current, ok = m[part]
		if !ok {
			return nil, false
		}
	}

	return current, true
}

func (i *I18n) normalizeLanguage(raw string) string {
	if raw == "" {
		return ""
	}
	tag, err := language.Parse(raw)
	if err != nil {
		return ""
	}
	base, _ := tag.Base()
	code := strings.ToLower(base.String())
	if _, ok := i.translations[code]; ok {
		return code
	}
	return ""
}
