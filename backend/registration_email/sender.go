package registration_email

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/morikuni/failure"
	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/domain/config"
	model_email "github.com/oinume/lekcije/backend/domain/model/email"
	"github.com/oinume/lekcije/backend/domain/repository"
	"github.com/oinume/lekcije/backend/infrastructure/send_grid"
	"github.com/oinume/lekcije/backend/model2"
)

type emailSender struct {
	sender repository.EmailSender
}

func NewEmailSender(httpClient *http.Client, appLogger *zap.Logger) *emailSender {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &emailSender{
		sender: send_grid.NewEmailSender(httpClient, appLogger),
	}
}

func (s *emailSender) Send(ctx context.Context, user *model2.User) error {
	t := model_email.NewTemplate("notifier", getEmailTemplate())
	data := struct {
		To     string
		WebURL string
	}{
		To:     user.Email,
		WebURL: config.WebURL(),
	}
	email, err := model_email.NewFromTemplate(t, data)
	if err != nil {
		return failure.Wrap(err, failure.Messagef("failed to create emailer.Email from template: to=%v", user.Email))
	}
	email.SetCustomArg("email_type", model2.EmailTypeRegistration)
	email.SetCustomArg("user_id", fmt.Sprint(user.ID))

	return s.sender.Send(ctx, email)
}

func getEmailTemplate() string {
	return strings.TrimSpace(`
From: lekcije <lekcije@lekcije.com>
To: {{ .To }}
Subject: lekcijeの登録が完了しました
Body: text/html
lekcijeにご登録いただきありがとうござます。

<a href="{{ .WebURL }}/me">こちら</a>からDMM英会話のお気に入りの講師をフォローしてみましょう。フォローすると講師が空きレッスンを登録した時にメールで通知が来るようになります。

ご質問などがございましたら、<a href="https://goo.gl/forms/CIGO3kpiQCGjtFD42">こちら</a>からお問い合わせ頂ければと思います。
lekcijeをよろしくお願いいたします。
	`)
}
