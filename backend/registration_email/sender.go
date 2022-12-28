package registration_email

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/oinume/lekcije/backend/domain/config"
	"github.com/oinume/lekcije/backend/emailer"
	"github.com/oinume/lekcije/backend/errors"
	"github.com/oinume/lekcije/backend/model2"
)

type emailSender struct {
	sender emailer.Sender
}

func NewEmailSender(httpClient *http.Client, appLogger *zap.Logger) *emailSender {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &emailSender{
		sender: emailer.NewSendGridSender(httpClient, appLogger),
	}
}

func (s *emailSender) Send(ctx context.Context, user *model2.User) error {
	t := emailer.NewTemplate("notifier", getEmailTemplate())
	data := struct {
		To     string
		WebURL string
	}{
		To:     user.Email,
		WebURL: config.WebURL(),
	}
	email, err := emailer.NewEmailFromTemplate(t, data)
	if err != nil {
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessagef("Failed to create emailer.Email from template: to=%v", user.Email),
		)
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
