package email

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ = fmt.Print

func TestNewEmailFromTemplate(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

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
	r.Nil(err)
	a.Equal("lekcije@lekcije.com", email.From.Address)
	a.Equal("oinume@gmail.com", email.Tos[0].Address)
	a.Equal("oinume", email.Tos[0].Name)
	a.Equal("テスト", email.Subject)
	a.Equal("text/html", email.BodyMIMEType)
	a.Equal("oinume さん\nこんにちは", email.BodyString())
}

func TestNewEmailsFromTemplate(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

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
	r.Nil(err)

	a.Equal(2, len(emails))
	a.Equal("lekcije@lekcije.com", emails[1].From.Address)
	a.Equal("akuwano@gmail.com", emails[1].Tos[0].Address)
	a.Equal("akuwano", emails[1].Tos[0].Name)
	a.Equal("テスト", emails[1].Subject)
	a.Equal("text/html", emails[1].BodyMIMEType)
	a.Equal("akuwano さん\nこんにちは", emails[1].BodyString())
}
