package controller

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = fmt.Print

func TestParseHTMLTemplates(t *testing.T) {
	a := assert.New(t)

	os.Chdir("../../")
	template := ParseHTMLTemplates(TemplatePath("indexLogout.html"))
	type Data struct {
		commonTemplateData
	}
	data := &Data{}
	var b bytes.Buffer
	err := template.Execute(&b, data)
	a.Nil(err)
	a.Contains(b.String(), "starter-template")
	a.Contains(b.String(), `<div id="func"></div>`)
	//fmt.Printf("--- output ---\n%v\n", b.String())
}

//func TestTemplate(t *testing.T) {
//	a := assert.New(t)
//
//	parent := template.New("parent")
//	parent.Funcs(map[string]interface{}{
//		"T": T,
//	})
//	_, err := parent.Parse(`parent: {{ T }}`)
//	a.Nil(err)
//
//	child := template.New("child")
//	child.Funcs(map[string]interface{}{
//		"T": T,
//	})
//	_, err = child.Parse(`child: {{ T }}`)
//	a.Nil(err)
//
//	var b bytes.Buffer
//	err = parent.Execute(&b, nil)
//	a.Nil(err)
//	fmt.Printf("--- output ---\n%v\n", b.String())
//}
