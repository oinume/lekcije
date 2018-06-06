package registration_email

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/emailer"
	"github.com/oinume/lekcije/server/errors"
	"github.com/oinume/lekcije/server/model"
)

var httpClient = &http.Client{}
var sender = emailer.NewSendGridSender(httpClient)

func Send(user *model.User) error {
	t := emailer.NewTemplate("notifier", getEmailTemplate())
	data := struct {
		To                string
		TeacherNames      string
		TeacherIDs        []uint32
		Teachers          map[uint32]*model.Teacher
		LessonsPerTeacher map[uint32]*model.TeacherLessons
		WebURL            string
	}{
		To:                user.Email,
		WebURL:            config.WebURL(),
	}
	email, err := emailer.NewEmailFromTemplate(t, data)
	if err != nil {
		return errors.NewInternalError(
			errors.WithError(err),
			errors.WithMessagef("Failed to create emailer.Email from template: to=%v", user.Email),
		)
	}
	email.SetCustomArg("email_type", model.EmailTypeNewLessonNotifier)
	email.SetCustomArg("user_id", fmt.Sprint(user.ID))

	return sender.Send(email)
}

func getEmailTemplate() string {
	return strings.TrimSpace(`
From: lekcije <lekcije@lekcije.com>
To: {{ .To }}
Subject: lekcijeの登録が完了しました
Body: text/html
lekcijeにご登録いただきありがとうござます。

● DMM英会話の講師をフォローしてみましょう
<a href="{{ .WebURL }}/me">こちら</a>からDMM英会話のお気に入りの講師をフォローしてみましょう。フォローすると講師が空きレッスンを登録した時にメールで通知がくるようになります。

<a href="https://goo.gl/forms/CIGO3kpiQCGjtFD42">お問い合わせ</a>
	`)
}
