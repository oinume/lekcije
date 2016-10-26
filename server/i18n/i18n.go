package i18n

import (
	"fmt"

	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/nicksnyder/go-i18n/i18n/translation"
	"github.com/nicksnyder/go-i18n/i18n/language"
	"gopkg.in/yaml.v2"
)

var T i18n.TranslateFunc

func Load(defaultLocale string) error {
	switch defaultLocale {
	case "en-US":
		LoadJSON(defaultLocale, En_US1)
		LoadJSON(defaultLocale, En_US2)
	default:
		return fmt.Errorf("Unsupported Locale: %s", defaultLocale)
	}

	// Obtain the default translation function for use.
	var err error
	T, err = NewTranslation(defaultLocale)
	if err != nil {
		return err
	}

	return nil

}

// NewTranslation obtains a translation function object for the
// specified locales.
func NewTranslation(userLocale string) (t i18n.TranslateFunc, err error) {
	t, err = i18n.Tfunc(userLocale)
	if err != nil {
		return t, err
	}
	return t, err
}

// LoadJSON takes a json document of translations and manually
// loads them into the system.
func loadYAML(userLocale string, translationDocument string) {
	tranDocuments := []map[string]interface{}{}
	err := yaml.Unmarshal([]byte(translationDocument), &tranDocuments)
	if err != nil {
		return err
	}

	for _, tranDocument := range tranDocuments {
		tran, err := translation.NewTranslation(tranDocument)
		if err != nil {
			return err
		}
		languages := language.MustParse(userLocale)
		i18n.AddTranslation(languages[0], tran)
	}

	return nil
}

/*
package main

import (
	"encoding/json"

	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/nicksnyder/go-i18n/i18n/language"
	"github.com/nicksnyder/go-i18n/i18n/translation"
	"fmt"
)

var (
// T is the translate function for the specified user
// locale and default locale specified during the load.
	T i18n.TranslateFunc
)

var En_US1 = `[
	{
		"id": "message1",
		"translation": "message1-1."
	},
	{
		"id": "message2",
		"translation": "message2."
	}
]`

var En_US2 = `[
	{
		"id": "message1",
		"translation": "message1-2."
	},
	{
		"id": "message3",
		"translation": "message3."
	}
]`

// Init initializes the local environment.
func Init(defaultLocale string) error {
	switch defaultLocale {
	case "en-US":
		LoadJSON(defaultLocale, En_US1)
		LoadJSON(defaultLocale, En_US2)
	default:
		return fmt.Errorf("Unsupported Locale: %s", defaultLocale)
	}

	// Obtain the default translation function for use.
	var err error
	T, err = NewTranslation(defaultLocale)
	if err != nil {
		return err
	}

	return nil
}

// NewTranslation obtains a translation function object for the
// specified locales.
func NewTranslation(userLocale string) (t i18n.TranslateFunc, err error) {
	t, err = i18n.Tfunc(userLocale)
	if err != nil {
		return t, err
	}

	return t, err
}

// LoadJSON takes a json document of translations and manually
// loads them into the system.
func LoadJSON(userLocale string, translationDocument string) error {
	tranDocuments := []map[string]interface{}{}
	err := json.Unmarshal([]byte(translationDocument), &tranDocuments)
	if err != nil {
		return err
	}

	for _, tranDocument := range tranDocuments {
		tran, err := translation.NewTranslation(tranDocument)
		if err != nil {
			return err
		}
		languages := language.MustParse(userLocale)
		i18n.AddTranslation(languages[0], tran)
	}

	return nil
}

func main() {
	Init("en-US")

	fmt.Println(T("message1"))
	fmt.Println(T("message2"))
	fmt.Println(T("message3"))
	fmt.Println(T("message4"))
}

 */
