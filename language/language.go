package language

import (
	_ "embed"
	"encoding/json"
	"log"
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
}

func (l *Language) Get(section, key string) string {
	if sec, ok := l.Translations[section]; ok {
		if val, ok := sec[key]; ok {
			return val
		}
		return missingTranslationValue(section, key)
	}
	return missingTranslationValue(section, key)
}

func missingTranslationValue(section, key string) string {
	return "i18n-missing:'" + section + ":" + key + "'"
}

func LoadLanguage(data []byte) *Language {
	var lang Language
	if err := json.Unmarshal(data, &lang); err != nil {
		log.Fatal("failed to decode language data: %w", err)
	}
	return &lang
}
