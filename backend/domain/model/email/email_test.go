package email

import (
	"strings"
	"testing"

	"github.com/oinume/lekcije/backend/internal/assertion"
)

func TestNewFromTemplate(t *testing.T) {
	s := `
From: lekcije@lekcije.com
To: {{ .Name }} <{{ .Email }}>
Subject: テスト
Body: text/html
{{ .Name }} さん
こんにちは
	`
	template := NewTemplate("TestNewEmailFromTemplate", strings.TrimSpace(s))
	data := struct {
		Name  string
		Email string
	}{
		"oinume",
		"oinume@gmail.com",
	}
	email, err := NewFromTemplate(template, data)
	if err != nil {
		t.Fatal(err)
	}

	assertion.AssertEqual(t, "lekcije@lekcije.com", email.From.Address, "")
	assertion.AssertEqual(t, "oinume@gmail.com", email.Tos[0].Address, "")
	assertion.AssertEqual(t, "oinume", email.Tos[0].Name, "")
	assertion.AssertEqual(t, "テスト", email.Subject, "")
	assertion.AssertEqual(t, "text/html", email.BodyMIMEType, "")
	assertion.AssertEqual(t, "oinume さん\nこんにちは", email.BodyString(), "")
}

func TestNewEmailsFromTemplate(t *testing.T) {
	s := `
From: lekcije@lekcije.com
To: {{ .Name }} <{{ .Email }}>
Subject: テスト
Body: text/html
{{ .Name }} さん
こんにちは
	`
	template := NewTemplate("TestNewEmailsFromTemplate", strings.TrimSpace(s))
	data := []struct {
		Name  string
		Email string
	}{
		{
			"oinume",
			"oinume@gmail.com",
		},
		{
			"akuwano",
			"akuwano@gmail.com",
		},
	}
	tmp := make([]interface{}, len(data))
	for i := range data {
		tmp[i] = data[i]
	}
	emails, err := NewEmailsFromTemplate(template, tmp)
	if err != nil {
		t.Fatal(err)
	}

	assertion.AssertEqual(t, 2, len(emails), "")
	assertion.AssertEqual(t, "lekcije@lekcije.com", emails[1].From.Address, "")
	assertion.AssertEqual(t, "akuwano@gmail.com", emails[1].Tos[0].Address, "")
	assertion.AssertEqual(t, "akuwano", emails[1].Tos[0].Name, "")
	assertion.AssertEqual(t, "テスト", emails[1].Subject, "")
	assertion.AssertEqual(t, "text/html", emails[1].BodyMIMEType, "")
	assertion.AssertEqual(t, "akuwano さん\nこんにちは", emails[1].BodyString(), "")
}
