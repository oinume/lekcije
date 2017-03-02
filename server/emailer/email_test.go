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
テストテスト
	`
	template := NewTemplate("parse", strings.TrimSpace(s))
	err := template.Parse()
	a.Nil(err)
	err = template.Execute(nil)
	a.Nil(err)

	a.Equal("Kazuhiro Oinuma", template.from.Name)
	a.Equal("oinume@gmail.com", template.from.Address)
	a.Equal("lekcije@lekcije.com", template.tos[0].Address)
	a.Equal("テスト", template.subject)
	a.Equal("text/html", template.mimeType)
}
