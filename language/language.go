package language

import (
	_ "embed"
	"encoding/json"
	"log"
	"strings"
)

//go:embed en.json
var enJSON []byte

//go:embed es.json
var esJSON []byte

var Languages = map[string]*Language{
	"EN": LoadLanguage(enJSON),
	"ES": LoadLanguage(esJSON),
}

type translation map[string]string

type Language struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Translations map[string]translation `json:"translations"`

	UsernameFirstWords  []string `json:"username_first_words"`
	UsernameSecondWords []string `json:"username_second_words"`
}

func (l *Language) Get(path string) string {
	parts := strings.SplitN(path, ".", 2)
	if len(parts) != 2 {
		return ""
	}
	section, key := parts[0], parts[1]

	if sec, ok := l.Translations[section]; ok {
		if val, ok := sec[key]; ok {
			return val
		}
		return missingTranslationValue(section, key)
	}
	return missingTranslationValue(section, key)
}

func missingTranslationValue(section, key string) string {
	return "i18n-missing:'" + section + "." + key + "'"
}

func LoadLanguage(data []byte) *Language {
	var lang Language
	if err := json.Unmarshal(data, &lang); err != nil {
		log.Fatal("failed to decode language data: %w", err)
	}
	return &lang
}
