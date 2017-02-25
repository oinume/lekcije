package emailer

import (
	"testing"
	"strings"

	"github.com/stretchr/testify/assert"
)

func TestTemplate_Parse(t *testing.T) {
	a := assert.New(t)
	s := `
AAAAA
	`
	template := NewTemplate("parse", strings.TrimSpace(s))
	err := template.Parse()
	a.Nil(err)
}
