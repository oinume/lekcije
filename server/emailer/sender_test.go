package emailer

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type transport struct {
	called bool
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.called = true
	resp := &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: http.StatusOK,
		Status:     "OK",
		Body:       ioutil.NopCloser(strings.NewReader("OK")),
	}
	return resp, nil
}

func TestSendGridSender_Send(t *testing.T) {
	a := assert.New(t)

	s := `
From: lekcije@lekcije.com
To: gmail <oinume@gmail.com>, oinume@lampetty.net
Subject: テスト
Body: text/html
oinume さん
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
	email, err := NewEmailFromTemplate(template, data)
	a.Nil(err)

	tr := &transport{}
	err = NewSendGridSender(&http.Client{Transport: tr}).Send(email)
	a.Nil(err)
	a.True(tr.called)
}
