package emailer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = fmt.Print

func TestTemplate_Execute(t *testing.T) {
	a := assert.New(t)
	s := `
From: Kazuhiro Oinuma <oinume@gmail.com>
To: lekcije@lekcije.com
Subject: テスト
Body: text/html
oinume さん
こんにちは
	`
	template := NewTemplate("TestTemplate_Execute", strings.TrimSpace(s))
	err := template.Parse()
	a.Nil(err)
	err = template.Execute(nil)
	a.Nil(err)

	email := template.emails[0]
	a.Equal("Kazuhiro Oinuma", email.From.Name)
	a.Equal("oinume@gmail.com", email.From.Address)
	a.Equal("lekcije@lekcije.com", email.Tos[0].Address)
	a.Equal("テスト", email.Subject)
	a.Equal("text/html", email.BodyMIMEType)
	a.Equal("oinume さん\nこんにちは", email.BodyString())
}

//func TestNewEmailFromTemplate(t *testing.T) {
//	a := assert.New(t)
//	template := NewTemplate("TestNewEmailFromTemplate", strings.TrimSpace(s))
//	NewEmailFromTemplate(template, )
//}
