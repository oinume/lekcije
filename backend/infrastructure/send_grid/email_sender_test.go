package send_grid

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"go.uber.org/zap/zapcore"

	model_email "github.com/oinume/lekcije/backend/domain/model/email"
	"github.com/oinume/lekcije/backend/internal/assertion"
	"github.com/oinume/lekcije/backend/logger"
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
		Body:       io.NopCloser(strings.NewReader("OK")),
	}
	return resp, nil
}

func TestSendGridSender_Send(t *testing.T) {
	s := `
From: lekcije@lekcije.com
To: gmail <oinume@gmail.com>, oinume@lampetty.net
Subject: テスト
Body: text/html
oinume さん
こんにちは
	`
	template := model_email.NewTemplate("TestNewEmailFromTemplate", strings.TrimSpace(s))
	data := struct {
		Name  string
		Email string
	}{
		"oinume",
		"oinume@gmail.com",
	}
	email, err := model_email.NewFromTemplate(template, data)
	if err != nil {
		t.Fatal(err)
	}

	email.SetCustomArg("userId", "1")
	email.SetCustomArg("teacherIds", "1,2,3")
	tr := &transport{}
	err = NewEmailSender(
		&http.Client{Transport: tr},
		logger.NewAppLogger(os.Stdout, zapcore.InfoLevel),
	).Send(context.Background(), email)
	if err != nil {
		t.Fatal(err)
	}
	assertion.AssertEqual(t, true, tr.called, "tr.called must be true")
}
